package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ReggieJTech/SuperAgent/internal/plugin"
	"github.com/ReggieJTech/SuperAgent/internal/queue"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// WebhookPlugin implements the webhook receiver plugin
type WebhookPlugin struct {
	*plugin.BasePlugin
	config     Config
	server     *http.Server
	router     *mux.Router
	queue      queue.Queue
	stopChan   chan struct{}

	// Sources
	sources map[string]*webhookSource

	// Metrics
	requestsReceived atomic.Int64
	requestsAccepted atomic.Int64
	requestsRejected atomic.Int64
	eventsProcessed  atomic.Int64
}

// webhookSource represents a configured webhook source
type webhookSource struct {
	config       SourceConfig
	auth         *Authenticator
	transformer  *Transformer
	rateLimiter  *rateLimiter
}

// rateLimiter implements a simple token bucket rate limiter
type rateLimiter struct {
	rate       int
	burst      int
	tokens     float64
	lastUpdate time.Time
}

// NewWebhookPlugin creates a new webhook plugin
func NewWebhookPlugin() plugin.Plugin {
	base := plugin.NewBasePlugin("webhook", "1.0.0", "HTTP/HTTPS webhook receiver")
	return &WebhookPlugin{
		BasePlugin: base,
		stopChan:   make(chan struct{}),
		sources:    make(map[string]*webhookSource),
	}
}

// Init initializes the webhook plugin
func (p *WebhookPlugin) Init(ctx context.Context, cfg plugin.PluginConfig) error {
	// Call base init
	if err := p.BasePlugin.Init(ctx, cfg); err != nil {
		return err
	}

	// Get queue
	if cfg.Queue != nil {
		if q, ok := cfg.Queue.(queue.Queue); ok {
			p.queue = q
		}
	}

	// Load configuration
	if cfg.ConfigFile != "" {
		if err := p.loadConfig(cfg.ConfigFile); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
	}

	// Set defaults
	p.config.SetDefaults()

	// Initialize sources
	for _, sourceCfg := range p.config.Sources {
		if !sourceCfg.Enabled {
			continue
		}

		source := &webhookSource{
			config:      sourceCfg,
			auth:        NewAuthenticator(sourceCfg.Auth),
			transformer: NewTransformer(sourceCfg.Transform),
		}

		// Initialize rate limiter if enabled
		if p.config.Global.RateLimit.Enabled {
			source.rateLimiter = &rateLimiter{
				rate:       p.config.Global.RateLimit.RequestsPerSecond,
				burst:      p.config.Global.RateLimit.Burst,
				tokens:     float64(p.config.Global.RateLimit.Burst),
				lastUpdate: time.Now(),
			}
		}

		p.sources[sourceCfg.Path] = source
		log.Info().
			Str("name", sourceCfg.Name).
			Str("path", sourceCfg.Path).
			Str("method", sourceCfg.Method).
			Msg("Webhook source registered")
	}

	// Setup HTTP router
	p.setupRouter()

	log.Info().
		Str("listen", p.config.ListenAddress).
		Int("sources", len(p.sources)).
		Msg("Webhook plugin initialized")

	return nil
}

// loadConfig loads configuration from file
func (p *WebhookPlugin) loadConfig(configFile string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	// Expand environment variables
	configStr := os.ExpandEnv(string(data))

	if err := yaml.Unmarshal([]byte(configStr), &p.config); err != nil {
		return err
	}

	return nil
}

// setupRouter sets up the HTTP router
func (p *WebhookPlugin) setupRouter() {
	p.router = mux.NewRouter()

	// Register webhook endpoints
	for path, source := range p.sources {
		source := source // capture loop variable
		p.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			p.handleWebhook(w, r, source)
		}).Methods(source.config.Method)
	}

	// Health check endpoint
	p.router.HandleFunc("/health", p.handleHealth).Methods("GET")
}

// Start starts the webhook plugin
func (p *WebhookPlugin) Start(ctx context.Context) error {
	log.Info().Str("plugin", p.Name()).Msg("Starting webhook plugin")

	// Create HTTP server
	p.server = &http.Server{
		Addr:         p.config.ListenAddress,
		Handler:      p.router,
		ReadTimeout:  p.config.Global.Timeout,
		WriteTimeout: p.config.Global.Timeout,
		IdleTimeout:  2 * time.Minute,
	}

	// Start server in goroutine
	go func() {
		log.Info().Str("address", p.config.ListenAddress).Msg("Webhook server starting")

		var err error
		if p.config.TLS.Enabled {
			err = p.server.ListenAndServeTLS(p.config.TLS.CertFile, p.config.TLS.KeyFile)
		} else {
			err = p.server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Webhook server error")
		}
	}()

	// Mark as started
	p.MarkStarted()

	return nil
}

// Stop stops the webhook plugin
func (p *WebhookPlugin) Stop(ctx context.Context) error {
	log.Info().Str("plugin", p.Name()).Msg("Stopping webhook plugin")

	close(p.stopChan)

	if p.server != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := p.server.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Error shutting down webhook server")
			return err
		}
	}

	// Mark as stopped
	p.MarkStopped()

	return nil
}

// handleWebhook handles incoming webhook requests
func (p *WebhookPlugin) handleWebhook(w http.ResponseWriter, r *http.Request, source *webhookSource) {
	p.requestsReceived.Add(1)

	if p.config.Logging.LogRequests {
		log.Debug().
			Str("source", source.config.Name).
			Str("path", r.URL.Path).
			Str("method", r.Method).
			Str("remote_addr", r.RemoteAddr).
			Msg("Webhook request received")
	}

	// Check rate limit
	if source.rateLimiter != nil && !source.rateLimiter.allow() {
		p.requestsRejected.Add(1)
		p.sendError(w, p.config.Response.Error.RateLimit)
		return
	}

	// Check IP whitelist
	if len(source.config.AllowedIPs) > 0 {
		clientIP := extractIP(r.RemoteAddr)
		if !isIPAllowed(clientIP, source.config.AllowedIPs) {
			p.requestsRejected.Add(1)
			p.sendError(w, p.config.Response.Error.Authentication)
			return
		}
	}

	// Read body
	body, err := io.ReadAll(io.LimitReader(r.Body, p.config.Global.MaxBodySize))
	if err != nil {
		p.requestsRejected.Add(1)
		p.sendError(w, p.config.Response.Error.InvalidPayload)
		return
	}
	defer r.Body.Close()

	// Authenticate
	if err := source.auth.Authenticate(r, body); err != nil {
		p.requestsRejected.Add(1)
		log.Debug().Err(err).Str("source", source.config.Name).Msg("Authentication failed")
		p.sendError(w, p.config.Response.Error.Authentication)
		return
	}

	// Transform payload
	events, err := source.transformer.Transform(body, source.config.Name)
	if err != nil {
		p.requestsRejected.Add(1)
		log.Error().Err(err).Str("source", source.config.Name).Msg("Transformation failed")
		p.sendError(w, p.config.Response.Error.InvalidPayload)
		return
	}

	// Queue events
	for _, event := range events {
		if err := p.queue.Enqueue(context.Background(), event); err != nil {
			log.Error().Err(err).Msg("Failed to queue event")
		} else {
			p.eventsProcessed.Add(1)
		}
	}

	p.requestsAccepted.Add(1)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(p.config.Response.Success.StatusCode)
	w.Write([]byte(p.config.Response.Success.Body))
}

// handleHealth handles health check requests
func (p *WebhookPlugin) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := p.Health()
	w.Header().Set("Content-Type", "application/json")
	if health.Status == "healthy" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(w).Encode(health)
}

// sendError sends an error response
func (p *WebhookPlugin) sendError(w http.ResponseWriter, errResp ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errResp.StatusCode)
	w.Write([]byte(errResp.Body))
}

// Stats returns plugin statistics
func (p *WebhookPlugin) Stats() map[string]interface{} {
	return map[string]interface{}{
		"requests_received": p.requestsReceived.Load(),
		"requests_accepted": p.requestsAccepted.Load(),
		"requests_rejected": p.requestsRejected.Load(),
		"events_processed":  p.eventsProcessed.Load(),
		"sources":           len(p.sources),
	}
}

// allow checks if the rate limiter allows a request
func (rl *rateLimiter) allow() bool {
	now := time.Now()
	elapsed := now.Sub(rl.lastUpdate).Seconds()

	// Add tokens based on elapsed time
	rl.tokens += elapsed * float64(rl.rate)
	if rl.tokens > float64(rl.burst) {
		rl.tokens = float64(rl.burst)
	}

	rl.lastUpdate = now

	// Check if we have tokens available
	if rl.tokens >= 1.0 {
		rl.tokens -= 1.0
		return true
	}

	return false
}

// extractIP extracts IP address from remote address
func extractIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}

// isIPAllowed checks if an IP is in the allowed list
func isIPAllowed(ip string, allowedIPs []string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, allowed := range allowedIPs {
		// Check if it's a CIDR range
		if strings.Contains(allowed, "/") {
			_, ipNet, err := net.ParseCIDR(allowed)
			if err != nil {
				continue
			}
			if ipNet.Contains(parsedIP) {
				return true
			}
		} else {
			// Direct IP match
			if allowed == ip {
				return true
			}
		}
	}

	return false
}
