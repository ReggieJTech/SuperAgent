package agent

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the main agent configuration
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	BigPanda   BigPandaConfig   `yaml:"bigpanda"` // Legacy: single endpoint (backward compatible)
	BigPandaEndpoints []BigPandaEndpoint `yaml:"bigpanda_endpoints"` // New: multiple endpoints
	Queue      QueueConfig      `yaml:"queue"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
	Logging    LoggingConfig    `yaml:"logging"`
	Auth       AuthConfig       `yaml:"auth"`
	Modules    []ModuleConfig   `yaml:"modules"`
	Security   SecurityConfig   `yaml:"security"`
}

// ServerConfig contains web server settings
type ServerConfig struct {
	ListenAddress string    `yaml:"listen_address"`
	TLS           TLSConfig `yaml:"tls"`
}

// TLSConfig contains TLS settings
type TLSConfig struct {
	Enabled      bool   `yaml:"enabled"`
	CertFile     string `yaml:"cert_file"`
	KeyFile      string `yaml:"key_file"`
	AutoGenerate bool   `yaml:"auto_generate"`
}

// BigPandaConfig contains BigPanda API settings (legacy single endpoint)
// Deprecated: Use BigPandaEndpoints for multi-endpoint support
type BigPandaConfig struct {
	APIURL       string         `yaml:"api_url"`
	StreamURL    string         `yaml:"stream_url"`
	HeartbeatURL string         `yaml:"heartbeat_url"`
	Token        string         `yaml:"token"`
	AppKey       string         `yaml:"app_key"`
	Batching     BatchingConfig `yaml:"batching"`
	Retry        RetryConfig    `yaml:"retry"`
	RateLimit    RateLimitConfig `yaml:"rate_limit"`
	Timeout      TimeoutConfig   `yaml:"timeout"`
}

// BigPandaEndpoint represents a single BigPanda destination
type BigPandaEndpoint struct {
	Name         string          `yaml:"name"`           // Unique identifier (e.g., "prod-network", "test-servers")
	Description  string          `yaml:"description"`    // Human-readable description
	Enabled      bool            `yaml:"enabled"`        // Enable/disable this endpoint
	APIURL       string          `yaml:"api_url"`        // BigPanda API URL
	StreamURL    string          `yaml:"stream_url"`     // Stream API URL (optional)
	HeartbeatURL string          `yaml:"heartbeat_url"`  // Heartbeat URL (optional)
	Token        string          `yaml:"token"`          // BigPanda API token
	AppKey       string          `yaml:"app_key"`        // BigPanda app key (integration identifier)
	Batching     BatchingConfig  `yaml:"batching"`       // Batching settings for this endpoint
	Retry        RetryConfig     `yaml:"retry"`          // Retry settings for this endpoint
	RateLimit    RateLimitConfig `yaml:"rate_limit"`     // Rate limiting for this endpoint
	Timeout      TimeoutConfig   `yaml:"timeout"`        // Timeout settings for this endpoint
	Tags         map[string]string `yaml:"tags"`         // Additional tags to add to all events
}

// BatchingConfig contains batching settings
type BatchingConfig struct {
	Enabled  bool          `yaml:"enabled"`
	MaxSize  int           `yaml:"max_size"`
	MaxWait  time.Duration `yaml:"max_wait"`
	MaxBytes int           `yaml:"max_bytes"`
}

// RetryConfig contains retry settings
type RetryConfig struct {
	MaxAttempts        int           `yaml:"max_attempts"`
	InitialBackoff     time.Duration `yaml:"initial_backoff"`
	MaxBackoff         time.Duration `yaml:"max_backoff"`
	BackoffMultiplier  float64       `yaml:"backoff_multiplier"`
}

// RateLimitConfig contains rate limiting settings
type RateLimitConfig struct {
	EventsPerSecond int `yaml:"events_per_second"`
	Burst           int `yaml:"burst"`
}

// TimeoutConfig contains timeout settings
type TimeoutConfig struct {
	Connect time.Duration `yaml:"connect"`
	Request time.Duration `yaml:"request"`
	Idle    time.Duration `yaml:"idle"`
}

// QueueConfig contains queue settings
type QueueConfig struct {
	Path        string    `yaml:"path"`
	MaxSize     int       `yaml:"max_size"`
	Persistence bool      `yaml:"persistence"`
	SyncWrites  bool      `yaml:"sync_writes"`
	DLQ         DLQConfig `yaml:"dlq"`
}

// DLQConfig contains dead letter queue settings
type DLQConfig struct {
	Enabled   bool          `yaml:"enabled"`
	MaxSize   int           `yaml:"max_size"`
	Retention time.Duration `yaml:"retention"`
}

// MonitoringConfig contains monitoring settings
type MonitoringConfig struct {
	Enabled           bool          `yaml:"enabled"`
	HeartbeatInterval time.Duration `yaml:"heartbeat_interval"`
	SelfMonitoring    bool          `yaml:"self_monitoring"`
	Metrics           MetricsConfig `yaml:"metrics"`
	Health            HealthConfig  `yaml:"health"`
}

// MetricsConfig contains metrics endpoint settings
type MetricsConfig struct {
	Enabled       bool   `yaml:"enabled"`
	ListenAddress string `yaml:"listen_address"`
	Path          string `yaml:"path"`
}

// HealthConfig contains health check settings
type HealthConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

// AuthConfig contains authentication settings
type AuthConfig struct {
	Local LocalAuthConfig `yaml:"local"`
	LDAP  LDAPAuthConfig  `yaml:"ldap"`
	SSO   SSOAuthConfig   `yaml:"sso"`
}

// LocalAuthConfig contains local authentication settings
type LocalAuthConfig struct {
	Enabled         bool          `yaml:"enabled"`
	UsersFile       string        `yaml:"users_file"`
	SessionDuration time.Duration `yaml:"session_duration"`
	JWTSecret       string        `yaml:"jwt_secret"`
}

// LDAPAuthConfig contains LDAP authentication settings
type LDAPAuthConfig struct {
	Enabled           bool   `yaml:"enabled"`
	Server            string `yaml:"server"`
	BindDN            string `yaml:"bind_dn"`
	BindPassword      string `yaml:"bind_password"`
	UserSearchBase    string `yaml:"user_search_base"`
	UserSearchFilter  string `yaml:"user_search_filter"`
	GroupSearchBase   string `yaml:"group_search_base"`
	GroupSearchFilter string `yaml:"group_search_filter"`
	RequiredGroup     string `yaml:"required_group"`
}

// SSOAuthConfig contains SSO authentication settings
type SSOAuthConfig struct {
	Enabled          bool   `yaml:"enabled"`
	Provider         string `yaml:"provider"`
	EntityID         string `yaml:"entity_id"`
	MetadataURL      string `yaml:"metadata_url"`
	ACSURL           string `yaml:"acs_url"`
	CertificateFile  string `yaml:"certificate_file"`
	PrivateKeyFile   string `yaml:"private_key_file"`
}

// ModuleConfig contains module settings
type ModuleConfig struct {
	Name       string `yaml:"name"`
	Enabled    bool   `yaml:"enabled"`
	ConfigFile string `yaml:"config_file"`
}

// SecurityConfig contains security settings
type SecurityConfig struct {
	Encryption EncryptionConfig `yaml:"encryption"`
	TLS        TLSSecurityConfig `yaml:"tls"`
	APIRateLimit APIRateLimitConfig `yaml:"api_rate_limit"`
}

// EncryptionConfig contains encryption settings
type EncryptionConfig struct {
	Enabled bool   `yaml:"enabled"`
	KeyFile string `yaml:"key_file"`
}

// TLSSecurityConfig contains TLS security settings
type TLSSecurityConfig struct {
	MinVersion   string   `yaml:"min_version"`
	CipherSuites []string `yaml:"cipher_suites"`
}

// APIRateLimitConfig contains API rate limiting settings
type APIRateLimitConfig struct {
	Enabled            bool `yaml:"enabled"`
	RequestsPerMinute  int  `yaml:"requests_per_minute"`
	Burst              int  `yaml:"burst"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Expand environment variables
	expanded := expandEnvVars(string(data))

	// Parse YAML
	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Set defaults
	cfg.SetDefaults()

	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Convert legacy single endpoint to new format if needed
	c.migrateLegacyConfig()

	// Validate BigPanda endpoints
	if len(c.BigPandaEndpoints) == 0 {
		return fmt.Errorf("at least one BigPanda endpoint is required")
	}

	endpointNames := make(map[string]bool)
	for i, endpoint := range c.BigPandaEndpoints {
		// Check for unique names
		if endpoint.Name == "" {
			return fmt.Errorf("bigpanda_endpoints[%d].name is required", i)
		}
		if endpointNames[endpoint.Name] {
			return fmt.Errorf("duplicate endpoint name: %s", endpoint.Name)
		}
		endpointNames[endpoint.Name] = true

		// Validate enabled endpoints
		if endpoint.Enabled {
			if endpoint.APIURL == "" {
				return fmt.Errorf("bigpanda_endpoints[%s].api_url is required", endpoint.Name)
			}
			if endpoint.Token == "" {
				return fmt.Errorf("bigpanda_endpoints[%s].token is required", endpoint.Name)
			}
			if endpoint.AppKey == "" {
				return fmt.Errorf("bigpanda_endpoints[%s].app_key is required", endpoint.Name)
			}
		}
	}

	// Validate queue settings
	if c.Queue.Path == "" {
		return fmt.Errorf("queue.path is required")
	}
	if c.Queue.MaxSize <= 0 {
		return fmt.Errorf("queue.max_size must be positive")
	}

	// Validate server settings
	if c.Server.ListenAddress == "" {
		return fmt.Errorf("server.listen_address is required")
	}

	// Validate log level
	validLevels := map[string]bool{
		"trace": true, "debug": true, "info": true,
		"warn": true, "error": true, "fatal": true,
	}
	if c.Logging.Level != "" && !validLevels[c.Logging.Level] {
		return fmt.Errorf("invalid log level: %s", c.Logging.Level)
	}

	return nil
}

// migrateLegacyConfig converts legacy single BigPanda config to new endpoints format
func (c *Config) migrateLegacyConfig() {
	// If new endpoints are already configured, nothing to do
	if len(c.BigPandaEndpoints) > 0 {
		return
	}

	// If legacy config exists, convert it
	if c.BigPanda.APIURL != "" && c.BigPanda.Token != "" && c.BigPanda.AppKey != "" {
		c.BigPandaEndpoints = []BigPandaEndpoint{
			{
				Name:         "default",
				Description:  "Default BigPanda endpoint (migrated from legacy config)",
				Enabled:      true,
				APIURL:       c.BigPanda.APIURL,
				StreamURL:    c.BigPanda.StreamURL,
				HeartbeatURL: c.BigPanda.HeartbeatURL,
				Token:        c.BigPanda.Token,
				AppKey:       c.BigPanda.AppKey,
				Batching:     c.BigPanda.Batching,
				Retry:        c.BigPanda.Retry,
				RateLimit:    c.BigPanda.RateLimit,
				Timeout:      c.BigPanda.Timeout,
			},
		}
	}
}

// SetDefaults sets default values for optional fields
func (c *Config) SetDefaults() {
	// Set defaults for all BigPanda endpoints
	for i := range c.BigPandaEndpoints {
		endpoint := &c.BigPandaEndpoints[i]

		// Default URLs
		if endpoint.StreamURL == "" && endpoint.APIURL != "" {
			endpoint.StreamURL = strings.Replace(endpoint.APIURL, "/alerts", "/stream", 1)
		}
		if endpoint.HeartbeatURL == "" && endpoint.APIURL != "" {
			endpoint.HeartbeatURL = strings.Replace(endpoint.APIURL, "/alerts", "/heartbeat", 1)
		}

		// Default batching
		if endpoint.Batching.MaxSize == 0 {
			endpoint.Batching.MaxSize = 100
		}
		if endpoint.Batching.MaxWait == 0 {
			endpoint.Batching.MaxWait = 5 * time.Second
		}
		if endpoint.Batching.MaxBytes == 0 {
			endpoint.Batching.MaxBytes = 1048576 // 1MB
		}

		// Default retry
		if endpoint.Retry.MaxAttempts == 0 {
			endpoint.Retry.MaxAttempts = 5
		}
		if endpoint.Retry.InitialBackoff == 0 {
			endpoint.Retry.InitialBackoff = 1 * time.Second
		}
		if endpoint.Retry.MaxBackoff == 0 {
			endpoint.Retry.MaxBackoff = 60 * time.Second
		}
		if endpoint.Retry.BackoffMultiplier == 0 {
			endpoint.Retry.BackoffMultiplier = 2.0
		}

		// Default rate limit
		if endpoint.RateLimit.EventsPerSecond == 0 {
			endpoint.RateLimit.EventsPerSecond = 1000
		}
		if endpoint.RateLimit.Burst == 0 {
			endpoint.RateLimit.Burst = 2000
		}

		// Default timeout
		if endpoint.Timeout.Connect == 0 {
			endpoint.Timeout.Connect = 10 * time.Second
		}
		if endpoint.Timeout.Request == 0 {
			endpoint.Timeout.Request = 30 * time.Second
		}
		if endpoint.Timeout.Idle == 0 {
			endpoint.Timeout.Idle = 90 * time.Second
		}

		// Initialize tags map if nil
		if endpoint.Tags == nil {
			endpoint.Tags = make(map[string]string)
		}
	}

	// Logging defaults
	if c.Logging.Level == "" {
		c.Logging.Level = "info"
	}
	if c.Logging.Format == "" {
		c.Logging.Format = "json"
	}

	// Monitoring defaults
	if c.Monitoring.HeartbeatInterval == 0 {
		c.Monitoring.HeartbeatInterval = 60 * time.Second
	}
	if c.Monitoring.Metrics.Path == "" {
		c.Monitoring.Metrics.Path = "/metrics"
	}
	if c.Monitoring.Health.Path == "" {
		c.Monitoring.Health.Path = "/health"
	}

	// Auth defaults
	if c.Auth.Local.SessionDuration == 0 {
		c.Auth.Local.SessionDuration = 24 * time.Hour
	}
}

// GetEndpointByName returns a BigPanda endpoint by name
func (c *Config) GetEndpointByName(name string) *BigPandaEndpoint {
	for i := range c.BigPandaEndpoints {
		if c.BigPandaEndpoints[i].Name == name {
			return &c.BigPandaEndpoints[i]
		}
	}
	return nil
}

// GetEnabledEndpoints returns all enabled BigPanda endpoints
func (c *Config) GetEnabledEndpoints() []BigPandaEndpoint {
	var enabled []BigPandaEndpoint
	for _, endpoint := range c.BigPandaEndpoints {
		if endpoint.Enabled {
			enabled = append(enabled, endpoint)
		}
	}
	return enabled
}

// expandEnvVars expands environment variables in the format ${VAR_NAME}
func expandEnvVars(s string) string {
	return os.Expand(s, func(key string) string {
		return os.Getenv(key)
	})
}
