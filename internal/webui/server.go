package webui

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Server represents the web UI server
type Server struct {
	config     Config
	agent      AgentInterface
	router     *mux.Router
	httpServer *http.Server
	wsHub      *WebSocketHub
	mu         sync.RWMutex
	started    bool
}

// Config contains web UI server configuration
type Config struct {
	ListenAddress string
	TLS           TLSConfig
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	IdleTimeout   time.Duration
}

// TLSConfig contains TLS configuration
type TLSConfig struct {
	Enabled      bool
	CertFile     string
	KeyFile      string
	AutoGenerate bool
}

// AgentInterface defines the interface to the agent
type AgentInterface interface {
	Health() map[string]interface{}
	Stats() map[string]interface{}
	Config() interface{}
	IsStarted() bool
	Uptime() time.Duration
}

// New creates a new web UI server
func New(config Config, agent AgentInterface) *Server {
	s := &Server{
		config: config,
		agent:  agent,
		router: mux.NewRouter(),
		wsHub:  NewWebSocketHub(),
	}

	// Set defaults
	if s.config.ReadTimeout == 0 {
		s.config.ReadTimeout = 15 * time.Second
	}
	if s.config.WriteTimeout == 0 {
		s.config.WriteTimeout = 15 * time.Second
	}
	if s.config.IdleTimeout == 0 {
		s.config.IdleTimeout = 60 * time.Second
	}

	// Start WebSocket hub
	go s.wsHub.Run()

	// Setup routes
	s.setupRoutes()

	return s
}

// Start starts the web UI server
func (s *Server) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return fmt.Errorf("server already started")
	}

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         s.config.ListenAddress,
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		var err error
		if s.config.TLS.Enabled {
			log.Info().
				Str("address", s.config.ListenAddress).
				Msg("Starting HTTPS server")

			// Check if auto-generate is enabled
			if s.config.TLS.AutoGenerate {
				if err := s.generateSelfSignedCert(); err != nil {
					log.Error().Err(err).Msg("Failed to generate self-signed certificate")
				}
			}

			// Configure TLS
			tlsConfig := &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				},
			}
			s.httpServer.TLSConfig = tlsConfig

			err = s.httpServer.ListenAndServeTLS(s.config.TLS.CertFile, s.config.TLS.KeyFile)
		} else {
			log.Info().
				Str("address", s.config.ListenAddress).
				Msg("Starting HTTP server")
			err = s.httpServer.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Web UI server error")
		}
	}()

	s.started = true
	log.Info().
		Str("address", s.config.ListenAddress).
		Bool("tls", s.config.TLS.Enabled).
		Msg("Web UI server started")

	return nil
}

// Stop stops the web UI server
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		return nil
	}

	log.Info().Msg("Stopping web UI server")

	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Error shutting down web UI server")
			return err
		}
	}

	s.started = false
	log.Info().Msg("Web UI server stopped")
	return nil
}

// setupRoutes sets up the HTTP routes
func (s *Server) setupRoutes() {
	// API routes
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Health endpoints
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")
	s.router.HandleFunc("/health/live", s.handleHealthLive).Methods("GET")
	s.router.HandleFunc("/health/ready", s.handleHealthReady).Methods("GET")

	// Stats endpoint
	api.HandleFunc("/stats", s.handleStats).Methods("GET")

	// Agent endpoints
	api.HandleFunc("/agent/info", s.handleAgentInfo).Methods("GET")
	api.HandleFunc("/agent/config", s.handleAgentConfig).Methods("GET")

	// Queue endpoints
	api.HandleFunc("/queue/stats", s.handleQueueStats).Methods("GET")
	api.HandleFunc("/queue/size", s.handleQueueSize).Methods("GET")

	// Plugin endpoints
	api.HandleFunc("/plugins", s.handlePluginsList).Methods("GET")
	api.HandleFunc("/plugins/{name}", s.handlePluginInfo).Methods("GET")
	api.HandleFunc("/plugins/{name}/stats", s.handlePluginStats).Methods("GET")

	// Event endpoints
	api.HandleFunc("/events/recent", s.handleRecentEvents).Methods("GET")
	api.HandleFunc("/events/dlq", s.handleDLQEvents).Methods("GET")

	// SNMP endpoints
	api.HandleFunc("/snmp/configs", s.handleSNMPConfigs).Methods("GET")
	api.HandleFunc("/snmp/unknown", s.handleSNMPUnknown).Methods("GET")

	// WebSocket endpoint for real-time events
	api.HandleFunc("/events/stream", s.handleWebSocket)

	// Authentication endpoints
	api.HandleFunc("/auth/login", s.handleLogin).Methods("POST")
	api.HandleFunc("/auth/logout", s.handleLogout).Methods("POST")
	api.HandleFunc("/auth/refresh", s.handleRefreshToken).Methods("POST")

	// Configuration Management endpoints
	// BigPanda configuration
	api.HandleFunc("/config/bigpanda", s.handleGetBigPandaConfig).Methods("GET")
	api.HandleFunc("/config/bigpanda", s.handleUpdateBigPandaConfig).Methods("PUT")

	// SNMP configuration
	api.HandleFunc("/config/snmp", s.handleGetSNMPConfig).Methods("GET")
	api.HandleFunc("/config/snmp", s.handleUpdateSNMPConfig).Methods("PUT")

	// Webhook configuration
	api.HandleFunc("/config/webhook", s.handleGetWebhookConfig).Methods("GET")
	api.HandleFunc("/config/webhook", s.handleUpdateWebhookConfig).Methods("PUT")

	// SNMP MIB and Event Config Management
	api.HandleFunc("/snmp/mibs/upload", s.handleUploadMIB).Methods("POST")
	api.HandleFunc("/snmp/events/generate", s.handleGenerateEventConfig).Methods("POST")
	api.HandleFunc("/snmp/events", s.handleListEventConfigs).Methods("GET")
	api.HandleFunc("/snmp/events/{name}", s.handleGetEventConfig).Methods("GET")
	api.HandleFunc("/snmp/events/{name}", s.handleUpdateEventConfig).Methods("PUT")
	api.HandleFunc("/snmp/events/{name}", s.handleDeleteEventConfig).Methods("DELETE")

	// Webhook Endpoint Management
	api.HandleFunc("/webhook/endpoints", s.handleCreateWebhookEndpoint).Methods("POST")
	api.HandleFunc("/webhook/endpoints/{name}", s.handleUpdateWebhookEndpoint).Methods("PUT")
	api.HandleFunc("/webhook/endpoints/{name}", s.handleDeleteWebhookEndpoint).Methods("DELETE")

	// Static files (for embedded React app)
	// s.router.PathPrefix("/").Handler(http.FileServer(http.FS(staticFiles)))

	// Apply middleware
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.corsMiddleware)
}

// generateSelfSignedCert generates a self-signed certificate
func (s *Server) generateSelfSignedCert() error {
	// Check if certificate already exists
	// This is a placeholder - actual implementation would use crypto/x509
	log.Info().Msg("Using existing certificate or generating self-signed certificate")
	return nil
}

// IsStarted returns whether the server is started
func (s *Server) IsStarted() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.started
}
