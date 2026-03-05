package snmp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

// EventConfig represents a trap event configuration
type EventConfig struct {
	Type           string                 `json:"type"`            // MIB or NO_MIB
	Trap           string                 `json:"trap"`            // Trap name
	TrapOID        string                 `json:"trap-oid"`        // Trap OID (for NO_MIB)
	TrapVarBinds   map[string]int         `json:"trap-var-binds"`  // Varbind name to position mapping
	MIB            string                 `json:"mib"`             // MIB name
	Copy           map[string]string               `json:"copy"`            // Copy field mappings
	Rename         map[string]string               `json:"rename"`          // Rename field mappings
	Set            map[string]interface{}          `json:"set"`             // Static field values
	MapStatus      map[string]map[string]string    `json:"map-status"`      // Status value mappings (field -> value -> mapped_value)
	Conditions     []interface{}                   `json:"conditions"`      // Conditional logic
	CustomActions  map[string]interface{} `json:"custom-actions"`  // Custom actions
	Primary        string                 `json:"primary"`         // Primary key field
	Secondary      string                 `json:"secondary"`       // Secondary key field
}

// EventConfigManager manages event configurations
type EventConfigManager struct {
	configs     map[string]*EventConfig // Key: trap name or OID
	configsDir  string
	mu          sync.RWMutex
}

// NewEventConfigManager creates a new event config manager
func NewEventConfigManager(configsDir string) (*EventConfigManager, error) {
	mgr := &EventConfigManager{
		configs:    make(map[string]*EventConfig),
		configsDir: configsDir,
	}

	// Load all event configs
	if err := mgr.LoadAll(); err != nil {
		return nil, fmt.Errorf("failed to load event configs: %w", err)
	}

	return mgr, nil
}

// LoadAll loads all event configuration files from the directory
func (m *EventConfigManager) LoadAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Info().Str("dir", m.configsDir).Msg("Loading event configurations")

	// Check if directory exists
	if _, err := os.Stat(m.configsDir); os.IsNotExist(err) {
		log.Warn().Str("dir", m.configsDir).Msg("Event configs directory does not exist")
		return nil
	}

	// Find all .ec files
	files, err := filepath.Glob(filepath.Join(m.configsDir, "*.ec"))
	if err != nil {
		return fmt.Errorf("failed to glob config files: %w", err)
	}

	loaded := 0
	for _, file := range files {
		if err := m.loadFile(file); err != nil {
			log.Error().Err(err).Str("file", file).Msg("Failed to load event config file")
			continue
		}
		loaded++
	}

	log.Info().
		Int("loaded", loaded).
		Int("total", len(m.configs)).
		Msg("Event configurations loaded")

	return nil
}

// loadFile loads a single event configuration file
func (m *EventConfigManager) loadFile(path string) error {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON array
	var configs []EventConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Store configs
	for i := range configs {
		config := &configs[i]

		// Use trap name as key for MIB type
		if config.Type == "MIB" {
			m.configs[config.Trap] = config
		} else if config.Type == "NO_MIB" {
			// Use trap OID as key for NO_MIB type
			if config.TrapOID != "" {
				m.configs[config.TrapOID] = config
			}
			// Also store by trap name if available
			if config.Trap != "" {
				m.configs[config.Trap] = config
			}
		}
	}

	log.Debug().
		Str("file", filepath.Base(path)).
		Int("count", len(configs)).
		Msg("Loaded event config file")

	return nil
}

// Get retrieves an event config by trap name or OID
func (m *EventConfigManager) Get(key string) (*EventConfig, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	config, ok := m.configs[key]
	return config, ok
}

// GetByOID retrieves an event config by OID
func (m *EventConfigManager) GetByOID(oid string) (*EventConfig, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// First try exact match
	if config, ok := m.configs[oid]; ok {
		return config, true
	}

	// Try matching with trailing .0 (common in SNMP)
	oidWithZero := strings.TrimSuffix(oid, ".0")
	if config, ok := m.configs[oidWithZero]; ok {
		return config, true
	}

	// Try without trailing .0
	oidNoZero := oid + ".0"
	if config, ok := m.configs[oidNoZero]; ok {
		return config, true
	}

	return nil, false
}

// Count returns the number of loaded configs
func (m *EventConfigManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.configs)
}

// List returns all trap names
func (m *EventConfigManager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.configs))
	for key := range m.configs {
		names = append(names, key)
	}
	return names
}

// Reload reloads all event configurations
func (m *EventConfigManager) Reload() error {
	log.Info().Msg("Reloading event configurations")

	// Create new configs map
	newConfigs := make(map[string]*EventConfig)

	// Temporarily store old configs
	m.mu.Lock()
	oldConfigs := m.configs
	m.configs = newConfigs
	m.mu.Unlock()

	// Load all configs
	if err := m.LoadAll(); err != nil {
		// Restore old configs on error
		m.mu.Lock()
		m.configs = oldConfigs
		m.mu.Unlock()
		return fmt.Errorf("failed to reload configs: %w", err)
	}

	log.Info().
		Int("old_count", len(oldConfigs)).
		Int("new_count", len(m.configs)).
		Msg("Event configurations reloaded")

	return nil
}
