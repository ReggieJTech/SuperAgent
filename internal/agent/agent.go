package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ReggieJTech/SuperAgent/internal/forwarder"
	"github.com/ReggieJTech/SuperAgent/internal/modules/snmp"
	"github.com/ReggieJTech/SuperAgent/internal/plugin"
	"github.com/ReggieJTech/SuperAgent/internal/queue"
	"github.com/ReggieJTech/SuperAgent/internal/webui"
	"github.com/rs/zerolog/log"
)

// Agent represents the main BigPanda Super Agent
type Agent struct {
	config *Config
	mu     sync.RWMutex

	// Components
	queue        queue.Queue
	forwarder    *forwarder.Forwarder
	pluginLoader *plugin.Loader
	webui        *webui.Server
	// monitor   *monitoring.Monitor

	// State
	started   bool
	startTime time.Time
	stopChan  chan struct{}
	wg        sync.WaitGroup
}

// New creates a new Agent instance
func New(cfg *Config) (*Agent, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	agent := &Agent{
		config:   cfg,
		stopChan: make(chan struct{}),
	}

	return agent, nil
}

// Start starts the agent and all its components
func (a *Agent) Start(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.started {
		return fmt.Errorf("agent already started")
	}

	log.Info().Msg("Starting BigPanda Super Agent")

	// Initialize components (placeholders for now)
	if err := a.initializeComponents(ctx); err != nil {
		return fmt.Errorf("failed to initialize components: %w", err)
	}

	// Start components
	if err := a.startComponents(ctx); err != nil {
		return fmt.Errorf("failed to start components: %w", err)
	}

	a.started = true
	a.startTime = time.Now()

	log.Info().Msg("Agent started successfully")
	return nil
}

// Stop stops the agent and all its components gracefully
func (a *Agent) Stop(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.started {
		return fmt.Errorf("agent not started")
	}

	log.Info().Msg("Stopping BigPanda Super Agent")

	// Signal all goroutines to stop
	close(a.stopChan)

	// Wait for all components to stop with timeout
	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info().Msg("All components stopped gracefully")
	case <-ctx.Done():
		log.Warn().Msg("Shutdown timeout exceeded, forcing stop")
		return ctx.Err()
	}

	// Stop components in reverse order
	if err := a.stopComponents(ctx); err != nil {
		log.Error().Err(err).Msg("Error stopping components")
		return err
	}

	a.started = false
	log.Info().Msg("Agent stopped successfully")
	return nil
}

// Reload reloads the agent configuration
func (a *Agent) Reload(configPath string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	log.Info().Str("config_path", configPath).Msg("Reloading configuration")

	// Load new configuration
	newConfig, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load new configuration: %w", err)
	}

	// Validate new configuration
	if err := newConfig.Validate(); err != nil {
		return fmt.Errorf("invalid new configuration: %w", err)
	}

	// Apply new configuration
	// TODO: Implement hot-reload logic for each component
	// For now, we just update the config reference
	a.config = newConfig

	log.Info().Msg("Configuration reloaded successfully")
	return nil
}

// Health returns the health status of the agent
func (a *Agent) Health() map[string]interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	health := map[string]interface{}{
		"status":  "healthy",
		"started": a.started,
		"uptime":  time.Since(a.startTime).String(),
	}

	if !a.started {
		health["status"] = "stopped"
		return health
	}

	// Check component health
	if a.queue != nil {
		health["queue"] = a.queue.Health()
	}
	if a.forwarder != nil {
		health["forwarder"] = a.forwarder.Health()
	}
	if a.pluginLoader != nil {
		health["plugins"] = a.pluginLoader.HealthCheck()
	}

	return health
}

// Stats returns agent statistics
func (a *Agent) Stats() map[string]interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	stats := map[string]interface{}{
		"uptime":     time.Since(a.startTime).String(),
		"start_time": a.startTime.Format(time.RFC3339),
	}

	if !a.started {
		return stats
	}

	// Collect stats from components
	if a.queue != nil {
		stats["queue"] = a.queue.Stats()
	}
	if a.forwarder != nil {
		stats["forwarder"] = a.forwarder.Stats()
	}
	if a.pluginLoader != nil {
		stats["plugins"] = a.pluginLoader.StatsReport()
	}

	return stats
}

// initializeComponents initializes all agent components
func (a *Agent) initializeComponents(ctx context.Context) error {
	log.Debug().Msg("Initializing components")

	// Initialize queue
	var err error
	if a.config.Queue.Persistence {
		a.queue, err = queue.NewBadgerQueue(queue.Config{
			Path:        a.config.Queue.Path,
			MaxSize:     a.config.Queue.MaxSize,
			Persistence: a.config.Queue.Persistence,
			SyncWrites:  a.config.Queue.SyncWrites,
			DLQ: queue.DLQConfig{
				Enabled:   a.config.Queue.DLQ.Enabled,
				MaxSize:   a.config.Queue.DLQ.MaxSize,
				Retention: a.config.Queue.DLQ.Retention,
			},
		})
	} else {
		a.queue = queue.NewMemoryQueue(queue.Config{
			Path:        a.config.Queue.Path,
			MaxSize:     a.config.Queue.MaxSize,
			Persistence: false,
		})
	}
	if err != nil {
		return fmt.Errorf("failed to initialize queue: %w", err)
	}
	log.Info().Msg("Queue initialized")

	// Initialize forwarder
	a.forwarder = forwarder.New(forwarder.Config{
		APIURL:   a.config.BigPanda.APIURL,
		Token:    a.config.BigPanda.Token,
		AppKey:   a.config.BigPanda.AppKey,
		Batching: forwarder.BatchConfig{
			Enabled:  a.config.BigPanda.Batching.Enabled,
			MaxSize:  a.config.BigPanda.Batching.MaxSize,
			MaxWait:  a.config.BigPanda.Batching.MaxWait,
			MaxBytes: a.config.BigPanda.Batching.MaxBytes,
		},
		Retry: forwarder.RetryConfig{
			MaxAttempts:       a.config.BigPanda.Retry.MaxAttempts,
			InitialBackoff:    a.config.BigPanda.Retry.InitialBackoff,
			MaxBackoff:        a.config.BigPanda.Retry.MaxBackoff,
			BackoffMultiplier: a.config.BigPanda.Retry.BackoffMultiplier,
		},
		RateLimit: forwarder.RateLimitConfig{
			EventsPerSecond: a.config.BigPanda.RateLimit.EventsPerSecond,
			Burst:           a.config.BigPanda.RateLimit.Burst,
		},
		Timeout: forwarder.TimeoutConfig{
			Connect: a.config.BigPanda.Timeout.Connect,
			Request: a.config.BigPanda.Timeout.Request,
			Idle:    a.config.BigPanda.Timeout.Idle,
		},
	}, a.queue)
	log.Info().Msg("Forwarder initialized")

	// Initialize plugin system
	a.pluginLoader = plugin.NewLoader()

	// Register built-in plugins
	a.pluginLoader.Registry().Register("mock", plugin.NewMockPlugin)
	a.pluginLoader.Registry().Register("snmp", snmp.NewSNMPPlugin)
	// TODO: Register webhook plugin

	// Load plugins from configuration
	pluginConfigs := make([]plugin.PluginConfig, 0, len(a.config.Modules))
	for _, modCfg := range a.config.Modules {
		pluginConfigs = append(pluginConfigs, plugin.PluginConfig{
			Name:       modCfg.Name,
			Enabled:    modCfg.Enabled,
			ConfigFile: modCfg.ConfigFile,
			Queue:      a.queue,
		})
	}

	if err := a.pluginLoader.LoadPlugins(ctx, pluginConfigs); err != nil {
		log.Warn().Err(err).Msg("Some plugins failed to load")
	}
	log.Info().Int("count", len(a.pluginLoader.Registry().List())).Msg("Plugins initialized")

	// Initialize web UI
	a.webui = webui.New(webui.Config{
		ListenAddress: a.config.Server.ListenAddress,
		TLS: webui.TLSConfig{
			Enabled:      a.config.Server.TLS.Enabled,
			CertFile:     a.config.Server.TLS.CertFile,
			KeyFile:      a.config.Server.TLS.KeyFile,
			AutoGenerate: a.config.Server.TLS.AutoGenerate,
		},
	}, a)
	log.Info().Msg("Web UI initialized")

	// TODO: Initialize monitoring
	// if a.config.Monitoring.Enabled {
	//     a.monitor = monitoring.New(a.config.Monitoring, a)
	// }

	log.Debug().Msg("Components initialized")
	return nil
}

// startComponents starts all agent components
func (a *Agent) startComponents(ctx context.Context) error {
	log.Debug().Msg("Starting components")

	// Queue is always ready (no-op)

	// Start forwarder
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.forwarder.Run(ctx, a.stopChan)
	}()
	log.Info().Msg("Forwarder started")

	// Start plugins
	if err := a.pluginLoader.StartPlugins(ctx); err != nil {
		log.Error().Err(err).Msg("Error starting some plugins")
	}

	// Start plugin health monitor
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.pluginLoader.MonitorPlugins(ctx, 30*time.Second)
	}()
	log.Info().Msg("Plugin monitor started")

	// Start web UI
	if err := a.webui.Start(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to start web UI")
		return err
	}
	log.Info().Msg("Web UI started")

	// TODO: Start monitoring
	// if a.config.Monitoring.Enabled {
	//     a.wg.Add(1)
	//     go func() {
	//         defer a.wg.Done()
	//         a.monitor.Run(ctx, a.stopChan)
	//     }()
	// }

	log.Debug().Msg("Components started")
	return nil
}

// stopComponents stops all agent components
func (a *Agent) stopComponents(ctx context.Context) error {
	log.Debug().Msg("Stopping components")

	// Stop in reverse order
	// TODO: Stop monitoring first
	// if a.monitor != nil {
	//     a.monitor.Stop(ctx)
	// }

	// Stop web UI
	if a.webui != nil {
		if err := a.webui.Stop(ctx); err != nil {
			log.Error().Err(err).Msg("Error stopping web UI")
		}
	}

	// Stop plugins
	if a.pluginLoader != nil {
		if err := a.pluginLoader.StopPlugins(ctx); err != nil {
			log.Error().Err(err).Msg("Error stopping plugins")
		}
	}

	// Stop forwarder
	if a.forwarder != nil {
		a.forwarder.Stop(ctx)
	}

	// Close queue last (after all producers are stopped)
	if a.queue != nil {
		if err := a.queue.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing queue")
		}
	}

	log.Debug().Msg("Components stopped")
	return nil
}

// Config returns a copy of the current configuration
func (a *Agent) Config() interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.config
}

// IsStarted returns whether the agent is started
func (a *Agent) IsStarted() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.started
}

// Uptime returns how long the agent has been running
func (a *Agent) Uptime() time.Duration {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if !a.started {
		return 0
	}
	return time.Since(a.startTime)
}
