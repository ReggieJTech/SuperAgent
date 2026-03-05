package snmp

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bigpanda/super-agent/internal/plugin"
	"github.com/bigpanda/super-agent/internal/queue"
	"github.com/gosnmp/gosnmp"
	"github.com/rs/zerolog/log"
)

// SNMPPlugin implements the SNMP trap receiver plugin
type SNMPPlugin struct {
	*plugin.BasePlugin
	config       Config
	configMgr    *EventConfigManager
	filter       *Filter
	trapListener *gosnmp.TrapListener
	stopChan     chan struct{}

	// Rate limiting
	sourceRates  map[string]*rateLimiter
	globalRate   *rateLimiter

	// Metrics
	trapsFiltered atomic.Int64
	trapsUnknown  atomic.Int64
}

// rateLimiter implements a simple token bucket rate limiter
type rateLimiter struct {
	rate       int
	burst      int
	tokens     float64
	lastUpdate time.Time
}

// NewSNMPPlugin creates a new SNMP plugin
func NewSNMPPlugin() plugin.Plugin {
	base := plugin.NewBasePlugin("snmp", "1.0.0", "SNMP trap receiver")
	return &SNMPPlugin{
		BasePlugin:  base,
		stopChan:    make(chan struct{}),
		sourceRates: make(map[string]*rateLimiter),
	}
}

// Init initializes the SNMP plugin
func (p *SNMPPlugin) Init(ctx context.Context, config plugin.PluginConfig) error {
	// Call base init
	if err := p.BasePlugin.Init(ctx, config); err != nil {
		return err
	}

	// Parse SNMP-specific config
	if config.Config != nil {
		cfg, ok := config.Config.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid config type")
		}

		// Convert map to Config struct (simplified, would use proper unmarshaling)
		p.parseConfig(cfg)
	}

	// Set defaults
	p.config.SetDefaults()

	// Initialize event config manager
	var err error
	p.configMgr, err = NewEventConfigManager(p.config.EventConfigsDir)
	if err != nil {
		return fmt.Errorf("failed to initialize event config manager: %w", err)
	}

	// Initialize filter
	p.filter, err = NewFilter(p.config.Filtering)
	if err != nil {
		return fmt.Errorf("failed to initialize filter: %w", err)
	}

	// Initialize rate limiters
	if p.config.RateLimiting.Enabled {
		p.globalRate = &rateLimiter{
			rate:       p.config.RateLimiting.Global,
			burst:      p.config.RateLimiting.Burst,
			tokens:     float64(p.config.RateLimiting.Burst),
			lastUpdate: time.Now(),
		}
	}

	log.Info().
		Str("listen", p.config.ListenAddress).
		Int("event_configs", p.configMgr.Count()).
		Msg("SNMP plugin initialized")

	return nil
}

// Start starts the SNMP plugin
func (p *SNMPPlugin) Start(ctx context.Context) error {
	if p.IsStarted() {
		return nil
	}

	log.Info().Str("plugin", p.Name()).Msg("Starting SNMP plugin")

	// Parse listen address
	parts := strings.Split(p.config.ListenAddress, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid listen address: %s", p.config.ListenAddress)
	}
	port := parts[1]

	// Create trap listener
	p.trapListener = gosnmp.NewTrapListener()
	p.trapListener.OnNewTrap = p.handleTrap
	p.trapListener.Params = gosnmp.Default

	// Set SNMP version
	switch p.config.SNMPVersion {
	case "1", "v1":
		p.trapListener.Params.Version = gosnmp.Version1
	case "2c", "v2c":
		p.trapListener.Params.Version = gosnmp.Version2c
	case "3", "v3":
		p.trapListener.Params.Version = gosnmp.Version3
		p.configureV3()
	default:
		p.trapListener.Params.Version = gosnmp.Version2c
	}

	// Set community string
	p.trapListener.Params.Community = p.config.Community

	// Start listener
	go func() {
		defer p.RecoverFromPanic()

		addr := "0.0.0.0:" + port
		log.Info().Str("address", addr).Msg("SNMP trap listener starting")

		if err := p.trapListener.Listen(addr); err != nil {
			log.Error().Err(err).Msg("SNMP trap listener error")
			p.SetUnhealthy("Listener failed", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	// Start config reload loop if enabled
	if p.config.AutoReload {
		go p.configReloadLoop(ctx)
	}

	p.MarkStarted()
	p.SetHealthy()

	return nil
}

// Stop stops the SNMP plugin
func (p *SNMPPlugin) Stop(ctx context.Context) error {
	if !p.IsStarted() {
		return nil
	}

	log.Info().Str("plugin", p.Name()).Msg("Stopping SNMP plugin")

	// Stop trap listener
	if p.trapListener != nil {
		p.trapListener.Close()
	}

	close(p.stopChan)
	p.MarkStopped()

	return nil
}

// Stats returns plugin statistics
func (p *SNMPPlugin) Stats() map[string]interface{} {
	stats := p.BasePlugin.Stats()

	stats["traps_filtered"] = p.trapsFiltered.Load()
	stats["traps_unknown"] = p.trapsUnknown.Load()
	stats["event_configs"] = p.configMgr.Count()
	stats["filter"] = p.filter.Stats()

	return stats
}

// handleTrap handles an incoming SNMP trap
func (p *SNMPPlugin) handleTrap(packet *gosnmp.SnmpPacket, addr *net.UDPAddr) {
	defer p.RecoverFromPanic()

	sourceIP := addr.IP.String()

	// Increment received counter
	p.IncrementReceived()

	if p.config.Logging.LogReceivedTraps {
		log.Debug().
			Str("source", sourceIP).
			Int("varbinds", len(packet.Variables)).
			Msg("Trap received")
	}

	// Parse trap
	trap, err := ParseTrap(packet, sourceIP)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse trap")
		p.IncrementErrors()
		return
	}

	// Check rate limiting
	if p.config.RateLimiting.Enabled {
		if !p.checkRateLimit(sourceIP) {
			log.Warn().Str("source", sourceIP).Msg("Rate limit exceeded, dropping trap")
			p.IncrementDropped()
			return
		}
	}

	// Check filter
	if p.filter.ShouldDrop(trap.TrapOID, sourceIP) {
		if p.config.Logging.LogFilteredTraps {
			log.Debug().
				Str("oid", trap.TrapOID).
				Str("source", sourceIP).
				Msg("Trap filtered")
		}
		p.trapsFiltered.Add(1)
		p.IncrementDropped()
		return
	}

	// Find event config
	eventConfig, found := p.findEventConfig(trap)
	if !found {
		p.handleUnknownTrap(trap)
		return
	}

	// Transform trap to event
	event, err := TransformTrap(trap, eventConfig)
	if err != nil {
		log.Error().Err(err).Str("trap", trap.TrapName).Msg("Failed to transform trap")
		p.IncrementErrors()
		return
	}

	// Send to queue
	if err := p.sendToQueue(event); err != nil {
		log.Error().Err(err).Msg("Failed to send event to queue")
		p.IncrementDropped()
		p.IncrementErrors()
		return
	}

	p.IncrementSent()
}

// findEventConfig finds the event configuration for a trap
func (p *SNMPPlugin) findEventConfig(trap *TrapData) (*EventConfig, bool) {
	// Try by trap name first
	if trap.TrapName != "" {
		if config, ok := p.configMgr.Get(trap.TrapName); ok {
			return config, true
		}
	}

	// Try by trap OID
	if trap.TrapOID != "" {
		if config, ok := p.configMgr.GetByOID(trap.TrapOID); ok {
			return config, true
		}
	}

	return nil, false
}

// handleUnknownTrap handles an unknown trap
func (p *SNMPPlugin) handleUnknownTrap(trap *TrapData) {
	p.trapsUnknown.Add(1)

	switch p.config.UnknownTraps.Action {
	case "drop":
		p.IncrementDropped()

	case "forward":
		// Create generic event
		event := queue.NewEvent()
		event.Status = "critical"
		if p.config.UnknownTraps.ForwardAsCritical {
			event.Status = "critical"
		}
		event.PrimaryKey = trap.SourceIP
		event.SecondaryKey = trap.TrapName
		event.Description = fmt.Sprintf("Unknown SNMP trap: %s", trap.TrapOID)
		event.SourceModule = "snmp"
		event.Tags["snmp_trap_oid"] = trap.TrapOID
		event.Tags["snmp_trap_name"] = trap.TrapName
		event.Tags["snmp_source_ip"] = trap.SourceIP

		if err := p.sendToQueue(event); err != nil {
			log.Error().Err(err).Msg("Failed to send unknown trap event")
			p.IncrementDropped()
		} else {
			p.IncrementSent()
		}

	case "log":
		fallthrough
	default:
		if p.config.Logging.LogUnknownTraps {
			if p.config.UnknownTraps.LogDetails {
				log.Warn().
					Str("trap_oid", trap.TrapOID).
					Str("trap_name", trap.TrapName).
					Str("source", trap.SourceIP).
					Int("varbinds", len(trap.VarBinds)).
					Interface("varbinds", trap.VarBinds).
					Msg("Unknown trap received")
			} else {
				log.Warn().
					Str("trap_oid", trap.TrapOID).
					Str("source", trap.SourceIP).
					Msg("Unknown trap received")
			}
		}
		p.IncrementDropped()
	}
}

// sendToQueue sends an event to the queue
func (p *SNMPPlugin) sendToQueue(event *queue.Event) error {
	q, ok := p.Queue().(queue.Queue)
	if !ok {
		return fmt.Errorf("invalid queue type")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return q.Enqueue(ctx, event)
}

// checkRateLimit checks if a trap from the source should be rate limited
func (p *SNMPPlugin) checkRateLimit(sourceIP string) bool {
	now := time.Now()

	// Check global rate limit
	if p.globalRate != nil {
		if !p.globalRate.allow(now) {
			return false
		}
	}

	// Check per-source rate limit
	if p.config.RateLimiting.PerSource > 0 {
		limiter, exists := p.sourceRates[sourceIP]
		if !exists {
			limiter = &rateLimiter{
				rate:       p.config.RateLimiting.PerSource,
				burst:      p.config.RateLimiting.Burst,
				tokens:     float64(p.config.RateLimiting.Burst),
				lastUpdate: now,
			}
			p.sourceRates[sourceIP] = limiter
		}

		if !limiter.allow(now) {
			return false
		}
	}

	return true
}

// allow checks if an action is allowed under rate limit
func (rl *rateLimiter) allow(now time.Time) bool {
	elapsed := now.Sub(rl.lastUpdate).Seconds()

	// Add tokens based on elapsed time
	rl.tokens += elapsed * float64(rl.rate)
	if rl.tokens > float64(rl.burst) {
		rl.tokens = float64(rl.burst)
	}

	rl.lastUpdate = now

	// Check if we have a token available
	if rl.tokens >= 1.0 {
		rl.tokens -= 1.0
		return true
	}

	return false
}

// configReloadLoop periodically reloads event configurations
func (p *SNMPPlugin) configReloadLoop(ctx context.Context) {
	ticker := time.NewTicker(p.config.ReloadInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopChan:
			return
		case <-ticker.C:
			if err := p.configMgr.Reload(); err != nil {
				log.Error().Err(err).Msg("Failed to reload event configs")
			}
		}
	}
}

// configureV3 configures SNMPv3 security
func (p *SNMPPlugin) configureV3() {
	// Configure SNMPv3 security parameters
	p.trapListener.Params.SecurityModel = gosnmp.UserSecurityModel

	msgFlags := gosnmp.NoAuthNoPriv
	switch p.config.V3.SecurityLevel {
	case "authNoPriv":
		msgFlags = gosnmp.AuthNoPriv
	case "authPriv":
		msgFlags = gosnmp.AuthPriv
	}
	p.trapListener.Params.MsgFlags = msgFlags

	// Set security name
	p.trapListener.Params.SecurityParameters = &gosnmp.UsmSecurityParameters{
		UserName:                 p.config.V3.SecurityName,
		AuthenticationProtocol:   p.getAuthProtocol(),
		AuthenticationPassphrase: p.config.V3.AuthPassword,
		PrivacyProtocol:          p.getPrivProtocol(),
		PrivacyPassphrase:        p.config.V3.PrivPassword,
	}
}

// getAuthProtocol returns the authentication protocol
func (p *SNMPPlugin) getAuthProtocol() gosnmp.SnmpV3AuthProtocol {
	switch strings.ToUpper(p.config.V3.AuthProtocol) {
	case "MD5":
		return gosnmp.MD5
	case "SHA":
		return gosnmp.SHA
	case "SHA224":
		return gosnmp.SHA224
	case "SHA256":
		return gosnmp.SHA256
	case "SHA384":
		return gosnmp.SHA384
	case "SHA512":
		return gosnmp.SHA512
	default:
		return gosnmp.NoAuth
	}
}

// getPrivProtocol returns the privacy protocol
func (p *SNMPPlugin) getPrivProtocol() gosnmp.SnmpV3PrivProtocol {
	switch strings.ToUpper(p.config.V3.PrivProtocol) {
	case "DES":
		return gosnmp.DES
	case "AES":
		return gosnmp.AES
	case "AES192":
		return gosnmp.AES192
	case "AES256":
		return gosnmp.AES256
	case "AES192C":
		return gosnmp.AES192C
	case "AES256C":
		return gosnmp.AES256C
	default:
		return gosnmp.NoPriv
	}
}

// parseConfig parses configuration from map
func (p *SNMPPlugin) parseConfig(cfg map[string]interface{}) {
	// Simplified config parsing - in production would use proper unmarshaling
	if val, ok := cfg["listen_address"].(string); ok {
		p.config.ListenAddress = val
	}
	if val, ok := cfg["snmp_version"].(string); ok {
		p.config.SNMPVersion = val
	}
	if val, ok := cfg["community"].(string); ok {
		p.config.Community = val
	}
	if val, ok := cfg["event_configs_dir"].(string); ok {
		p.config.EventConfigsDir = val
	}
	// ... continue for other fields
}
