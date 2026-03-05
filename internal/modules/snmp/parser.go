package snmp

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ReggieJTech/SuperAgent/internal/queue"
	"github.com/gosnmp/gosnmp"
	"github.com/rs/zerolog/log"
)

// TrapData represents parsed SNMP trap data
type TrapData struct {
	TrapOID    string
	TrapName   string
	SourceIP   string
	Timestamp  time.Time
	VarBinds   map[string]interface{}
	RawVarBinds []gosnmp.SnmpPDU
}

// ParseTrap parses an SNMP trap packet
func ParseTrap(packet *gosnmp.SnmpPacket, sourceIP string) (*TrapData, error) {
	if packet == nil {
		return nil, fmt.Errorf("packet is nil")
	}

	data := &TrapData{
		SourceIP:    sourceIP,
		Timestamp:   time.Now(),
		VarBinds:    make(map[string]interface{}),
		RawVarBinds: packet.Variables,
	}

	// Extract trap OID (usually the first or second varbind)
	for _, vb := range packet.Variables {
		oidStr := vb.Name

		// SNMPv2-MIB::snmpTrapOID
		if strings.HasPrefix(oidStr, ".1.3.6.1.6.3.1.1.4.1") ||
		   strings.HasPrefix(oidStr, "1.3.6.1.6.3.1.1.4.1") {
			data.TrapOID = fmt.Sprintf("%v", vb.Value)
			break
		}
	}

	// If no trap OID found in varbinds, use enterprise OID
	if data.TrapOID == "" && packet.Enterprise != "" {
		data.TrapOID = packet.Enterprise
	}

	// Extract trap name from OID (last part)
	if data.TrapOID != "" {
		parts := strings.Split(strings.Trim(data.TrapOID, "."), ".")
		if len(parts) > 0 {
			data.TrapName = parts[len(parts)-1]
		}
	}

	// Parse varbinds
	for _, vb := range packet.Variables {
		oidStr := vb.Name
		value := vb.Value

		// Convert value to appropriate type
		var convertedValue interface{}
		switch v := value.(type) {
		case []byte:
			convertedValue = string(v)
		case int, int32, int64, uint, uint32, uint64:
			convertedValue = v
		default:
			convertedValue = fmt.Sprintf("%v", v)
		}

		data.VarBinds[oidStr] = convertedValue
	}

	log.Debug().
		Str("trap_oid", data.TrapOID).
		Str("trap_name", data.TrapName).
		Str("source", sourceIP).
		Int("varbinds", len(data.VarBinds)).
		Msg("Trap parsed")

	return data, nil
}

// TransformTrap transforms trap data using event configuration
func TransformTrap(trap *TrapData, config *EventConfig) (*queue.Event, error) {
	event := queue.NewEvent()

	// Set basic fields
	event.Timestamp = trap.Timestamp
	event.SourceIP = trap.SourceIP
	event.SourceModule = "snmp"

	// Add standard SNMP fields
	event.Tags["snmp_trap_oid"] = trap.TrapOID
	event.Tags["snmp_trap_name"] = trap.TrapName
	event.Tags["snmp_source_ip"] = trap.SourceIP

	// Map varbinds for NO_MIB type
	mappedVars := make(map[string]interface{})
	if config.Type == "NO_MIB" && len(config.TrapVarBinds) > 0 {
		// Map varbinds by position
		for name, position := range config.TrapVarBinds {
			if position > 0 && position <= len(trap.RawVarBinds) {
				vb := trap.RawVarBinds[position-1]
				mappedVars[name] = formatValue(vb.Value)
			}
		}
	} else {
		// Use OID-based varbinds
		mappedVars = trap.VarBinds
	}

	// Apply transformations in order:
	// 1. Copy fields
	for src, dst := range config.Copy {
		if val, ok := mappedVars[src]; ok {
			event.CustomFields[dst] = val
		}
	}

	// 2. Rename fields
	for src, dst := range config.Rename {
		if val, ok := mappedVars[src]; ok {
			event.CustomFields[dst] = val
			delete(mappedVars, src)
		}
	}

	// 3. Set static values
	for key, val := range config.Set {
		if key == "status" {
			event.Status = fmt.Sprintf("%v", val)
		} else if key == "description" {
			event.Description = fmt.Sprintf("%v", val)
		} else {
			event.CustomFields[key] = val
		}
	}

	// 4. Map status values
	if len(config.MapStatus) > 0 {
		for field, valueMap := range config.MapStatus {
			if val, ok := mappedVars[field]; ok {
				valStr := fmt.Sprintf("%v", val)
				if mappedStatus, ok := valueMap[valStr]; ok {
					event.Status = mappedStatus
				}
			}
		}
	}

	// If status not set, try to determine from trap name
	if event.Status == "" {
		event.Status = DetermineStatus(trap.TrapName, trap.TrapOID)
	}

	// 5. Set primary and secondary keys
	primaryKey := evaluateKey(config.Primary, event, mappedVars, trap)
	secondaryKey := evaluateKey(config.Secondary, event, mappedVars, trap)

	if primaryKey == "" {
		primaryKey = trap.SourceIP
	}
	if secondaryKey == "" {
		secondaryKey = trap.TrapName
	}

	event.PrimaryKey = primaryKey
	event.SecondaryKey = secondaryKey

	// Add all remaining mapped vars as custom fields
	for key, val := range mappedVars {
		if _, exists := event.CustomFields[key]; !exists {
			event.CustomFields[key] = val
		}
	}

	// Validate event
	if err := event.Validate(); err != nil {
		return nil, fmt.Errorf("invalid transformed event: %w", err)
	}

	log.Debug().
		Str("trap", trap.TrapName).
		Str("status", event.Status).
		Str("primary", event.PrimaryKey).
		Str("secondary", event.SecondaryKey).
		Msg("Trap transformed")

	return event, nil
}

// evaluateKey evaluates a key expression
func evaluateKey(expr string, event *queue.Event, vars map[string]interface{}, trap *TrapData) string {
	if expr == "" {
		return ""
	}

	// Handle special values
	switch expr {
	case "snmp_source_ip":
		return trap.SourceIP
	case "snmp_trap_name":
		return trap.TrapName
	case "snmp_trap_oid":
		return trap.TrapOID
	}

	// Check if it's a field reference
	if val, ok := vars[expr]; ok {
		return fmt.Sprintf("%v", val)
	}

	// Check custom fields
	if val, ok := event.CustomFields[expr]; ok {
		return fmt.Sprintf("%v", val)
	}

	// Check tags
	if val, ok := event.Tags[expr]; ok {
		return val
	}

	// Return as-is if not found
	return expr
}

// formatValue formats a varbind value to a string
func formatValue(value interface{}) interface{} {
	switch v := value.(type) {
	case []byte:
		return string(v)
	case nil:
		return ""
	default:
		return v
	}
}

// DetermineStatus determines event status from trap name using pattern matching
func DetermineStatus(trapName, trapOID string) string {
	name := strings.ToLower(trapName)

	// Critical patterns
	criticalPatterns := []string{
		"critical", "error", "fail", "failure", "down", "outage",
		"alarm", "alert", "emergency", "fatal", "severe",
		"offline", "unreachable", "dead", "lost", "disconnect",
	}

	for _, pattern := range criticalPatterns {
		if strings.Contains(name, pattern) {
			return "critical"
		}
	}

	// Warning patterns
	warningPatterns := []string{
		"warning", "warn", "degraded", "threshold", "high",
		"low", "minor", "caution",
	}

	for _, pattern := range warningPatterns {
		if strings.Contains(name, pattern) {
			return "warning"
		}
	}

	// OK/Clear patterns
	okPatterns := []string{
		"clear", "cleared", "ok", "normal", "up", "online",
		"available", "recovered", "restore", "resolved",
	}

	for _, pattern := range okPatterns {
		if strings.Contains(name, pattern) {
			return "ok"
		}
	}

	// Default to critical for unknown traps
	return "critical"
}

// ValidateTrapOID validates a trap OID format
func ValidateTrapOID(oid string) bool {
	if oid == "" {
		return false
	}

	// OID should be numeric parts separated by dots
	pattern := regexp.MustCompile(`^\.?\d+(\.\d+)*$`)
	return pattern.MatchString(oid)
}
