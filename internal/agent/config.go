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
	BigPanda   BigPandaConfig   `yaml:"bigpanda"`
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

// BigPandaConfig contains BigPanda API settings
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
	// Validate BigPanda settings
	if c.BigPanda.APIURL == "" {
		return fmt.Errorf("bigpanda.api_url is required")
	}
	if c.BigPanda.Token == "" {
		return fmt.Errorf("bigpanda.token is required")
	}
	if c.BigPanda.AppKey == "" {
		return fmt.Errorf("bigpanda.app_key is required")
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

// SetDefaults sets default values for optional fields
func (c *Config) SetDefaults() {
	// BigPanda defaults
	if c.BigPanda.StreamURL == "" {
		c.BigPanda.StreamURL = strings.Replace(c.BigPanda.APIURL, "/alerts", "/stream", 1)
	}
	if c.BigPanda.HeartbeatURL == "" {
		c.BigPanda.HeartbeatURL = strings.Replace(c.BigPanda.APIURL, "/alerts", "/heartbeat", 1)
	}
	if c.BigPanda.Batching.MaxSize == 0 {
		c.BigPanda.Batching.MaxSize = 100
	}
	if c.BigPanda.Batching.MaxWait == 0 {
		c.BigPanda.Batching.MaxWait = 5 * time.Second
	}
	if c.BigPanda.Retry.MaxAttempts == 0 {
		c.BigPanda.Retry.MaxAttempts = 5
	}
	if c.BigPanda.Retry.InitialBackoff == 0 {
		c.BigPanda.Retry.InitialBackoff = 1 * time.Second
	}
	if c.BigPanda.Retry.MaxBackoff == 0 {
		c.BigPanda.Retry.MaxBackoff = 60 * time.Second
	}
	if c.BigPanda.Retry.BackoffMultiplier == 0 {
		c.BigPanda.Retry.BackoffMultiplier = 2.0
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

// expandEnvVars expands environment variables in the format ${VAR_NAME}
func expandEnvVars(s string) string {
	return os.Expand(s, func(key string) string {
		return os.Getenv(key)
	})
}
