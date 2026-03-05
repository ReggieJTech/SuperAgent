package plugin

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

// BasePlugin provides common functionality for plugins
type BasePlugin struct {
	name        string
	version     string
	description string
	config      PluginConfig
	queue       interface{}

	started   atomic.Bool
	startTime time.Time
	mu        sync.RWMutex

	// Metrics
	eventsReceived atomic.Int64
	eventsSent     atomic.Int64
	eventsDropped  atomic.Int64
	errors         atomic.Int64

	// Health
	health HealthStatus
}

// NewBasePlugin creates a new base plugin
func NewBasePlugin(name, version, description string) *BasePlugin {
	return &BasePlugin{
		name:        name,
		version:     version,
		description: description,
		health:      NewHealthyStatus(),
	}
}

// Name returns the plugin name
func (p *BasePlugin) Name() string {
	return p.name
}

// Version returns the plugin version
func (p *BasePlugin) Version() string {
	return p.version
}

// Description returns the plugin description
func (p *BasePlugin) Description() string {
	return p.description
}

// Init initializes the plugin
func (p *BasePlugin) Init(ctx context.Context, config PluginConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.config = config
	p.queue = config.Queue

	log.Info().
		Str("plugin", p.name).
		Str("version", p.version).
		Msg("Plugin initialized")

	return nil
}

// Health returns the plugin health status
func (p *BasePlugin) Health() HealthStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Update timestamp
	health := p.health
	health.Timestamp = time.Now()

	// Add basic details
	if health.Details == nil {
		health.Details = make(map[string]interface{})
	}
	health.Details["started"] = p.started.Load()
	if p.started.Load() {
		health.Details["uptime"] = time.Since(p.startTime).String()
	}

	return health
}

// Stats returns plugin statistics
func (p *BasePlugin) Stats() map[string]interface{} {
	stats := map[string]interface{}{
		"events_received": p.eventsReceived.Load(),
		"events_sent":     p.eventsSent.Load(),
		"events_dropped":  p.eventsDropped.Load(),
		"errors":          p.errors.Load(),
		"started":         p.started.Load(),
	}

	if p.started.Load() {
		stats["uptime"] = time.Since(p.startTime).String()
		stats["start_time"] = p.startTime.Format(time.RFC3339)
	}

	return stats
}

// SetHealth updates the plugin health status
func (p *BasePlugin) SetHealth(status HealthStatus) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.health = status
}

// SetHealthy marks the plugin as healthy
func (p *BasePlugin) SetHealthy() {
	p.SetHealth(NewHealthyStatus())
}

// SetDegraded marks the plugin as degraded
func (p *BasePlugin) SetDegraded(message string, details map[string]interface{}) {
	p.SetHealth(NewDegradedStatus(message, details))
}

// SetUnhealthy marks the plugin as unhealthy
func (p *BasePlugin) SetUnhealthy(message string, details map[string]interface{}) {
	p.SetHealth(NewUnhealthyStatus(message, details))
}

// IsStarted returns whether the plugin is started
func (p *BasePlugin) IsStarted() bool {
	return p.started.Load()
}

// MarkStarted marks the plugin as started
func (p *BasePlugin) MarkStarted() {
	p.started.Store(true)
	p.startTime = time.Now()
}

// MarkStopped marks the plugin as stopped
func (p *BasePlugin) MarkStopped() {
	p.started.Store(false)
}

// Config returns the plugin configuration
func (p *BasePlugin) Config() PluginConfig {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.config
}

// Queue returns the event queue
func (p *BasePlugin) Queue() interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.queue
}

// IncrementReceived increments the received counter
func (p *BasePlugin) IncrementReceived() {
	p.eventsReceived.Add(1)
}

// IncrementSent increments the sent counter
func (p *BasePlugin) IncrementSent() {
	p.eventsSent.Add(1)
}

// IncrementDropped increments the dropped counter
func (p *BasePlugin) IncrementDropped() {
	p.eventsDropped.Add(1)
}

// IncrementErrors increments the error counter
func (p *BasePlugin) IncrementErrors() {
	p.errors.Add(1)
}

// RecoverFromPanic recovers from a panic in a plugin goroutine
func (p *BasePlugin) RecoverFromPanic() {
	if r := recover(); r != nil {
		log.Error().
			Str("plugin", p.name).
			Interface("panic", r).
			Msg("Plugin panic recovered")

		p.IncrementErrors()
		p.SetUnhealthy("Plugin panicked", map[string]interface{}{
			"panic": r,
		})
	}
}
