package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

// Loader handles plugin loading and lifecycle management
type Loader struct {
	registry *Registry
}

// NewLoader creates a new plugin loader
func NewLoader() *Loader {
	return &Loader{
		registry: NewRegistry(),
	}
}

// Registry returns the plugin registry
func (l *Loader) Registry() *Registry {
	return l.registry
}

// LoadPlugins loads multiple plugins from configuration
func (l *Loader) LoadPlugins(ctx context.Context, configs []PluginConfig) error {
	log.Info().Int("count", len(configs)).Msg("Loading plugins")

	var errors []error
	loaded := 0

	for _, config := range configs {
		if !config.Enabled {
			log.Info().Str("plugin", config.Name).Msg("Plugin disabled, skipping")
			continue
		}

		if err := l.registry.Load(ctx, config); err != nil {
			log.Error().
				Err(err).
				Str("plugin", config.Name).
				Msg("Failed to load plugin")
			errors = append(errors, err)
			continue
		}

		loaded++
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to load %d plugin(s), loaded %d successfully", len(errors), loaded)
	}

	log.Info().Int("count", loaded).Msg("All plugins loaded successfully")
	return nil
}

// StartPlugins starts all loaded plugins
func (l *Loader) StartPlugins(ctx context.Context) error {
	log.Info().Msg("Starting plugins")

	if err := l.registry.StartAll(ctx); err != nil {
		return err
	}

	log.Info().Msg("All plugins started successfully")
	return nil
}

// StopPlugins stops all loaded plugins
func (l *Loader) StopPlugins(ctx context.Context) error {
	log.Info().Msg("Stopping plugins")

	if err := l.registry.StopAll(ctx); err != nil {
		return err
	}

	log.Info().Msg("All plugins stopped successfully")
	return nil
}

// StartPluginWithRecovery starts a plugin with automatic recovery
func (l *Loader) StartPluginWithRecovery(ctx context.Context, name string) {
	plugin, err := l.registry.Get(name)
	if err != nil {
		log.Error().Err(err).Str("plugin", name).Msg("Plugin not found")
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Str("plugin", name).
					Interface("panic", r).
					Msg("Plugin crashed, attempting restart")

				// Wait before restart
				time.Sleep(5 * time.Second)

				// Try to restart
				if err := l.registry.Start(ctx, name); err != nil {
					log.Error().
						Err(err).
						Str("plugin", name).
						Msg("Failed to restart plugin")
				} else {
					log.Info().Str("plugin", name).Msg("Plugin restarted successfully")
				}
			}
		}()

		if err := plugin.Start(ctx); err != nil {
			log.Error().
				Err(err).
				Str("plugin", name).
				Msg("Plugin start failed")
		}
	}()
}

// MonitorPlugins monitors plugin health and restarts unhealthy plugins
func (l *Loader) MonitorPlugins(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			health := l.registry.Health()

			for name, status := range health {
				if status.Status == "unhealthy" {
					log.Warn().
						Str("plugin", name).
						Str("message", status.Message).
						Msg("Plugin unhealthy, attempting restart")

					// Try to restart
					if err := l.registry.Stop(ctx, name); err != nil {
						log.Error().Err(err).Str("plugin", name).Msg("Failed to stop unhealthy plugin")
						continue
					}

					time.Sleep(2 * time.Second)

					if err := l.registry.Start(ctx, name); err != nil {
						log.Error().Err(err).Str("plugin", name).Msg("Failed to restart unhealthy plugin")
					} else {
						log.Info().Str("plugin", name).Msg("Unhealthy plugin restarted")
					}
				}
			}
		}
	}
}

// HealthCheck performs a health check on all plugins
func (l *Loader) HealthCheck() map[string]interface{} {
	health := l.registry.Health()

	result := map[string]interface{}{
		"status":  "healthy",
		"plugins": make(map[string]interface{}),
	}

	unhealthyCount := 0
	degradedCount := 0

	for name, status := range health {
		result["plugins"].(map[string]interface{})[name] = map[string]interface{}{
			"status":  status.Status,
			"message": status.Message,
			"details": status.Details,
		}

		switch status.Status {
		case "unhealthy":
			unhealthyCount++
		case "degraded":
			degradedCount++
		}
	}

	// Overall status
	if unhealthyCount > 0 {
		result["status"] = "unhealthy"
		result["unhealthy_count"] = unhealthyCount
	} else if degradedCount > 0 {
		result["status"] = "degraded"
		result["degraded_count"] = degradedCount
	}

	result["total_plugins"] = len(health)

	return result
}

// StatsReport generates a statistics report for all plugins
func (l *Loader) StatsReport() map[string]interface{} {
	stats := l.registry.Stats()

	report := map[string]interface{}{
		"plugins": stats,
	}

	// Calculate totals
	var totalReceived, totalSent, totalDropped, totalErrors int64
	for _, pluginStats := range stats {
		if received, ok := pluginStats["events_received"].(int64); ok {
			totalReceived += received
		}
		if sent, ok := pluginStats["events_sent"].(int64); ok {
			totalSent += sent
		}
		if dropped, ok := pluginStats["events_dropped"].(int64); ok {
			totalDropped += dropped
		}
		if errors, ok := pluginStats["errors"].(int64); ok {
			totalErrors += errors
		}
	}

	report["totals"] = map[string]interface{}{
		"events_received": totalReceived,
		"events_sent":     totalSent,
		"events_dropped":  totalDropped,
		"errors":          totalErrors,
	}

	return report
}
