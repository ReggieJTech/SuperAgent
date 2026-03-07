package snmp

import (
	"time"
)

// Config contains SNMP module configuration
type Config struct {
	ListenAddress   string         `yaml:"listen_address"`
	SNMPVersion     string         `yaml:"snmp_version"`
	Community       string         `yaml:"community"`
	V3              V3Config       `yaml:"v3"`
	EventConfigsDir string         `yaml:"event_configs_dir"`
	MIBsDir         string         `yaml:"mibs_dir"`
	AutoReload      bool           `yaml:"auto_reload"`
	ReloadInterval  time.Duration  `yaml:"reload_interval"`
	Filtering       FilterConfig   `yaml:"filtering"`
	UnknownTraps    UnknownConfig  `yaml:"unknown_traps"`
	RateLimiting    RateLimitConfig `yaml:"rate_limiting"`
	Performance     PerfConfig     `yaml:"performance"`
	Logging         LogConfig      `yaml:"logging"`
	MIB             MIBConfig      `yaml:"mib"`
	Routing         RoutingConfig  `yaml:"routing"` // NEW: Route traps to BigPanda endpoints
}

// V3Config contains SNMPv3 security settings
type V3Config struct {
	SecurityLevel string `yaml:"security_level"`
	AuthProtocol  string `yaml:"auth_protocol"`
	AuthPassword  string `yaml:"auth_password"`
	PrivProtocol  string `yaml:"priv_protocol"`
	PrivPassword  string `yaml:"priv_password"`
	SecurityName  string `yaml:"security_name"`
}

// FilterConfig contains filtering configuration
type FilterConfig struct {
	Enabled bool          `yaml:"enabled"`
	Rules   []FilterRule  `yaml:"rules"`
}

// FilterRule represents a single filter rule
type FilterRule struct {
	Action  string `yaml:"action"`  // drop, accept
	Type    string `yaml:"type"`    // oid, source, source_network
	Pattern string `yaml:"pattern"` // OID pattern, IP, or CIDR
}

// UnknownConfig contains unknown trap handling configuration
type UnknownConfig struct {
	Action            string `yaml:"action"`              // log, drop, forward
	LogDetails        bool   `yaml:"log_details"`
	ForwardAsCritical bool   `yaml:"forward_as_critical"`
}

// RateLimitConfig contains rate limiting configuration
type RateLimitConfig struct {
	Enabled   bool `yaml:"enabled"`
	PerSource int  `yaml:"per_source"` // Max traps/sec per source
	Global    int  `yaml:"global"`     // Max traps/sec total
	Burst     int  `yaml:"burst"`      // Burst allowance
}

// PerfConfig contains performance tuning configuration
type PerfConfig struct {
	Workers    int `yaml:"workers"`
	BufferSize int `yaml:"buffer_size"`
	BatchSize  int `yaml:"batch_size"`
}

// LogConfig contains logging configuration
type LogConfig struct {
	Debug             bool `yaml:"debug"`
	LogReceivedTraps  bool `yaml:"log_received_traps"`
	LogFilteredTraps  bool `yaml:"log_filtered_traps"`
	LogUnknownTraps   bool `yaml:"log_unknown_traps"`
}

// MIBConfig contains MIB compilation settings
type MIBConfig struct {
	PreferNoMIB bool   `yaml:"prefer_no_mib"`
	AutoCompile bool   `yaml:"auto_compile"`
	CacheDir    string `yaml:"cache_dir"`
}

// RoutingConfig contains BigPanda endpoint routing configuration
type RoutingConfig struct {
	DefaultEndpoints []string       `yaml:"default_endpoints"` // Default endpoint(s) for all traps
	Rules            []RoutingRule  `yaml:"rules"`            // Conditional routing rules
}

// RoutingRule defines conditional routing to BigPanda endpoints
type RoutingRule struct {
	Name        string   `yaml:"name"`        // Rule name/description
	MatchType   string   `yaml:"match_type"`  // oid, oid_prefix, vendor, source, source_network
	MatchValue  string   `yaml:"match_value"` // Value to match against
	Endpoints   []string `yaml:"endpoints"`   // Target endpoint name(s)
	Priority    int      `yaml:"priority"`    // Higher priority rules evaluated first
}

// SetDefaults sets default values for the config
func (c *Config) SetDefaults() {
	if c.ListenAddress == "" {
		c.ListenAddress = "0.0.0.0:162"
	}
	if c.SNMPVersion == "" {
		c.SNMPVersion = "2c"
	}
	if c.Community == "" {
		c.Community = "public"
	}
	if c.EventConfigsDir == "" {
		c.EventConfigsDir = "/etc/bigpanda-agent/snmp/event_configs"
	}
	if c.MIBsDir == "" {
		c.MIBsDir = "/etc/bigpanda-agent/snmp/mibs"
	}
	if c.ReloadInterval == 0 {
		c.ReloadInterval = 60 * time.Second
	}
	if c.RateLimiting.PerSource == 0 {
		c.RateLimiting.PerSource = 100
	}
	if c.RateLimiting.Global == 0 {
		c.RateLimiting.Global = 1000
	}
	if c.RateLimiting.Burst == 0 {
		c.RateLimiting.Burst = 200
	}
	if c.Performance.Workers == 0 {
		c.Performance.Workers = 4
	}
	if c.Performance.BufferSize == 0 {
		c.Performance.BufferSize = 1000
	}
	if c.Performance.BatchSize == 0 {
		c.Performance.BatchSize = 50
	}
	if c.UnknownTraps.Action == "" {
		c.UnknownTraps.Action = "log"
	}
	if c.MIB.CacheDir == "" {
		c.MIB.CacheDir = "/var/lib/bigpanda-agent/mib_cache"
	}
	// Initialize routing with default endpoint if not specified
	if len(c.Routing.DefaultEndpoints) == 0 {
		c.Routing.DefaultEndpoints = []string{"default"}
	}
}
