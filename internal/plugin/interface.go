package plugin

import (
	"context"
	"time"
)

// Plugin represents a receiver module plugin
type Plugin interface {
	// Name returns the plugin name
	Name() string

	// Version returns the plugin version
	Version() string

	// Description returns a short description of the plugin
	Description() string

	// Init initializes the plugin with its configuration
	Init(ctx context.Context, config PluginConfig) error

	// Start starts the plugin
	Start(ctx context.Context) error

	// Stop stops the plugin gracefully
	Stop(ctx context.Context) error

	// Health returns the plugin health status
	Health() HealthStatus

	// Stats returns plugin statistics
	Stats() map[string]interface{}
}

// PluginConfig contains plugin configuration
type PluginConfig struct {
	// Name is the plugin name
	Name string

	// Enabled indicates if the plugin is enabled
	Enabled bool

	// ConfigFile is the path to the plugin configuration file
	ConfigFile string

	// Config contains the parsed configuration
	Config interface{}

	// Queue is the event queue for sending events
	Queue interface{}
}

// HealthStatus represents plugin health
type HealthStatus struct {
	Status    string                 `json:"status"`    // healthy, degraded, unhealthy
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details"`
	Timestamp time.Time              `json:"timestamp"`
}

// NewHealthyStatus creates a healthy status
func NewHealthyStatus() HealthStatus {
	return HealthStatus{
		Status:    "healthy",
		Message:   "Plugin is operating normally",
		Details:   make(map[string]interface{}),
		Timestamp: time.Now(),
	}
}

// NewDegradedStatus creates a degraded status
func NewDegradedStatus(message string, details map[string]interface{}) HealthStatus {
	if details == nil {
		details = make(map[string]interface{})
	}
	return HealthStatus{
		Status:    "degraded",
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}
}

// NewUnhealthyStatus creates an unhealthy status
func NewUnhealthyStatus(message string, details map[string]interface{}) HealthStatus {
	if details == nil {
		details = make(map[string]interface{})
	}
	return HealthStatus{
		Status:    "unhealthy",
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}
}

// PluginFactory is a function that creates a new plugin instance
type PluginFactory func() Plugin

// PluginMetadata contains plugin metadata
type PluginMetadata struct {
	Name        string
	Version     string
	Description string
	Author      string
	Factory     PluginFactory
}
