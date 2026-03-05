package plugin

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// Registry manages all plugins
type Registry struct {
	plugins  map[string]Plugin
	factories map[string]PluginFactory
	mu       sync.RWMutex
}

// NewRegistry creates a new plugin registry
func NewRegistry() *Registry {
	return &Registry{
		plugins:   make(map[string]Plugin),
		factories: make(map[string]PluginFactory),
	}
}

// Register registers a plugin factory
func (r *Registry) Register(name string, factory PluginFactory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.factories[name]; exists {
		return fmt.Errorf("plugin %s already registered", name)
	}

	r.factories[name] = factory
	log.Debug().Str("plugin", name).Msg("Plugin factory registered")
	return nil
}

// Load loads and initializes a plugin
func (r *Registry) Load(ctx context.Context, config PluginConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if plugin already loaded
	if _, exists := r.plugins[config.Name]; exists {
		return fmt.Errorf("plugin %s already loaded", config.Name)
	}

	// Get factory
	factory, exists := r.factories[config.Name]
	if !exists {
		return fmt.Errorf("plugin %s not registered", config.Name)
	}

	// Create plugin instance
	plugin := factory()
	if plugin == nil {
		return fmt.Errorf("plugin factory returned nil for %s", config.Name)
	}

	// Load plugin configuration if config file specified
	if config.ConfigFile != "" {
		pluginConfig, err := r.loadPluginConfig(config.ConfigFile)
		if err != nil {
			return fmt.Errorf("failed to load plugin config: %w", err)
		}
		config.Config = pluginConfig
	}

	// Initialize plugin
	if err := plugin.Init(ctx, config); err != nil {
		return fmt.Errorf("failed to initialize plugin %s: %w", config.Name, err)
	}

	// Store plugin
	r.plugins[config.Name] = plugin

	log.Info().
		Str("plugin", config.Name).
		Str("version", plugin.Version()).
		Msg("Plugin loaded successfully")

	return nil
}

// Start starts a plugin
func (r *Registry) Start(ctx context.Context, name string) error {
	r.mu.RLock()
	plugin, exists := r.plugins[name]
	r.mu.RUnlock()

	if !exists {
		return fmt.Errorf("plugin %s not loaded", name)
	}

	log.Info().Str("plugin", name).Msg("Starting plugin")

	if err := plugin.Start(ctx); err != nil {
		return fmt.Errorf("failed to start plugin %s: %w", name, err)
	}

	return nil
}

// Stop stops a plugin
func (r *Registry) Stop(ctx context.Context, name string) error {
	r.mu.RLock()
	plugin, exists := r.plugins[name]
	r.mu.RUnlock()

	if !exists {
		return fmt.Errorf("plugin %s not loaded", name)
	}

	log.Info().Str("plugin", name).Msg("Stopping plugin")

	if err := plugin.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop plugin %s: %w", name, err)
	}

	return nil
}

// StartAll starts all loaded plugins
func (r *Registry) StartAll(ctx context.Context) error {
	r.mu.RLock()
	plugins := make([]Plugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		plugins = append(plugins, p)
	}
	r.mu.RUnlock()

	var errors []error
	for _, plugin := range plugins {
		if err := plugin.Start(ctx); err != nil {
			log.Error().
				Err(err).
				Str("plugin", plugin.Name()).
				Msg("Failed to start plugin")
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to start %d plugin(s)", len(errors))
	}

	return nil
}

// StopAll stops all loaded plugins
func (r *Registry) StopAll(ctx context.Context) error {
	r.mu.RLock()
	plugins := make([]Plugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		plugins = append(plugins, p)
	}
	r.mu.RUnlock()

	var errors []error
	for _, plugin := range plugins {
		if err := plugin.Stop(ctx); err != nil {
			log.Error().
				Err(err).
				Str("plugin", plugin.Name()).
				Msg("Failed to stop plugin")
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to stop %d plugin(s)", len(errors))
	}

	return nil
}

// Get returns a plugin by name
func (r *Registry) Get(name string) (Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return plugin, nil
}

// List returns all loaded plugins
func (r *Registry) List() []Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make([]Plugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		plugins = append(plugins, p)
	}

	return plugins
}

// ListNames returns all loaded plugin names
func (r *Registry) ListNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.plugins))
	for name := range r.plugins {
		names = append(names, name)
	}

	return names
}

// Health returns health status for all plugins
func (r *Registry) Health() map[string]HealthStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()

	health := make(map[string]HealthStatus)
	for name, plugin := range r.plugins {
		health[name] = plugin.Health()
	}

	return health
}

// Stats returns statistics for all plugins
func (r *Registry) Stats() map[string]map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := make(map[string]map[string]interface{})
	for name, plugin := range r.plugins {
		stats[name] = plugin.Stats()
	}

	return stats
}

// Unload unloads a plugin
func (r *Registry) Unload(ctx context.Context, name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not loaded", name)
	}

	// Stop plugin first
	if err := plugin.Stop(ctx); err != nil {
		log.Error().Err(err).Str("plugin", name).Msg("Error stopping plugin during unload")
	}

	// Remove from registry
	delete(r.plugins, name)

	log.Info().Str("plugin", name).Msg("Plugin unloaded")
	return nil
}

// Reload reloads a plugin
func (r *Registry) Reload(ctx context.Context, config PluginConfig) error {
	// Unload existing plugin
	if err := r.Unload(ctx, config.Name); err != nil {
		log.Warn().Err(err).Str("plugin", config.Name).Msg("Error unloading plugin during reload")
	}

	// Load plugin again
	if err := r.Load(ctx, config); err != nil {
		return fmt.Errorf("failed to reload plugin %s: %w", config.Name, err)
	}

	// Start plugin
	if err := r.Start(ctx, config.Name); err != nil {
		return fmt.Errorf("failed to start reloaded plugin %s: %w", config.Name, err)
	}

	log.Info().Str("plugin", config.Name).Msg("Plugin reloaded successfully")
	return nil
}

// loadPluginConfig loads plugin configuration from a file
func (r *Registry) loadPluginConfig(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Expand environment variables
	expanded := os.Expand(string(data), func(key string) string {
		return os.Getenv(key)
	})

	// Parse YAML
	var config map[string]interface{}
	if err := yaml.Unmarshal([]byte(expanded), &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return config, nil
}
