package webhook

import (
	"time"
)

// Config represents the webhook module configuration
type Config struct {
	ListenAddress string    `yaml:"listen_address"`
	TLS           TLSConfig `yaml:"tls"`
	Global        GlobalConfig `yaml:"global"`
	Sources       []SourceConfig `yaml:"sources"`
	Response      ResponseConfig `yaml:"response"`
	Logging       LoggingConfig `yaml:"logging"`
}

// TLSConfig represents TLS configuration
type TLSConfig struct {
	Enabled      bool   `yaml:"enabled"`
	CertFile     string `yaml:"cert_file"`
	KeyFile      string `yaml:"key_file"`
	AutoGenerate bool   `yaml:"auto_generate"`
}

// GlobalConfig represents global webhook settings
type GlobalConfig struct {
	Timeout      time.Duration `yaml:"timeout"`
	MaxBodySize  int64         `yaml:"max_body_size"`
	RateLimit    RateLimitConfig `yaml:"rate_limit"`
}

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	Enabled            bool `yaml:"enabled"`
	RequestsPerSecond  int  `yaml:"requests_per_second"`
	Burst              int  `yaml:"burst"`
}

// SourceConfig represents a webhook source configuration
type SourceConfig struct {
	Name       string            `yaml:"name"`
	Enabled    bool              `yaml:"enabled"`
	Path       string            `yaml:"path"`
	Method     string            `yaml:"method"`
	Auth       AuthConfig        `yaml:"auth"`
	AllowedIPs []string          `yaml:"allowed_ips"`
	Transform  TransformConfig   `yaml:"transform"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Type      string `yaml:"type"` // bearer, apikey, basic, hmac, none
	Token     string `yaml:"token"`
	Header    string `yaml:"header"`
	Key       string `yaml:"key"`
	Secret    string `yaml:"secret"`
	Algorithm string `yaml:"algorithm"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

// TransformConfig represents transformation configuration
type TransformConfig struct {
	FieldMap      map[string]string            `yaml:"field_map"`
	StatusMap     map[string]string            `yaml:"status_map"`
	PrimaryKey    string                       `yaml:"primary_key"`
	SecondaryKey  string                       `yaml:"secondary_key"`
	Set           map[string]interface{}       `yaml:"set"`
	Template      string                       `yaml:"template"`
}

// ResponseConfig represents response configuration
type ResponseConfig struct {
	Async   bool                 `yaml:"async"`
	Success SuccessResponse      `yaml:"success"`
	Error   ErrorResponses       `yaml:"error"`
}

// SuccessResponse represents success response configuration
type SuccessResponse struct {
	StatusCode int    `yaml:"status_code"`
	Body       string `yaml:"body"`
}

// ErrorResponses represents error response configurations
type ErrorResponses struct {
	Authentication ErrorResponse `yaml:"authentication"`
	RateLimit      ErrorResponse `yaml:"rate_limit"`
	InvalidPayload ErrorResponse `yaml:"invalid_payload"`
}

// ErrorResponse represents an error response configuration
type ErrorResponse struct {
	StatusCode int    `yaml:"status_code"`
	Body       string `yaml:"body"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Debug        bool `yaml:"debug"`
	LogRequests  bool `yaml:"log_requests"`
	LogResponses bool `yaml:"log_responses"`
	LogPayload   bool `yaml:"log_payload"`
}

// SetDefaults sets default configuration values
func (c *Config) SetDefaults() {
	if c.ListenAddress == "" {
		c.ListenAddress = "0.0.0.0:8080"
	}
	if c.Global.Timeout == 0 {
		c.Global.Timeout = 30 * time.Second
	}
	if c.Global.MaxBodySize == 0 {
		c.Global.MaxBodySize = 10485760 // 10MB
	}
	if c.Global.RateLimit.RequestsPerSecond == 0 {
		c.Global.RateLimit.RequestsPerSecond = 100
	}
	if c.Global.RateLimit.Burst == 0 {
		c.Global.RateLimit.Burst = 200
	}
	if c.Response.Success.StatusCode == 0 {
		c.Response.Success.StatusCode = 200
	}
	if c.Response.Success.Body == "" {
		c.Response.Success.Body = `{"status": "accepted"}`
	}
}
