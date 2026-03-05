package snmp

import (
	"net"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

// Filter handles trap filtering
type Filter struct {
	config FilterConfig
	rules  []*compiledRule
}

// compiledRule is a compiled filter rule
type compiledRule struct {
	original FilterRule
	pattern  *regexp.Regexp
	network  *net.IPNet
}

// NewFilter creates a new filter
func NewFilter(config FilterConfig) (*Filter, error) {
	f := &Filter{
		config: config,
		rules:  make([]*compiledRule, 0, len(config.Rules)),
	}

	// Compile rules
	for _, rule := range config.Rules {
		compiled, err := f.compileRule(rule)
		if err != nil {
			log.Error().Err(err).Str("pattern", rule.Pattern).Msg("Failed to compile filter rule")
			continue
		}
		f.rules = append(f.rules, compiled)
	}

	log.Info().Int("rules", len(f.rules)).Msg("Filter initialized")
	return f, nil
}

// compileRule compiles a filter rule
func (f *Filter) compileRule(rule FilterRule) (*compiledRule, error) {
	compiled := &compiledRule{
		original: rule,
	}

	switch rule.Type {
	case "oid":
		// Convert OID pattern to regex
		// Replace * with .* and escape dots
		pattern := strings.ReplaceAll(rule.Pattern, ".", "\\.")
		pattern = strings.ReplaceAll(pattern, "*", ".*")
		pattern = "^" + pattern + "$"

		regex, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		compiled.pattern = regex

	case "source":
		// Exact IP match (we'll use string comparison)
		// No compilation needed

	case "source_network":
		// Parse CIDR
		_, network, err := net.ParseCIDR(rule.Pattern)
		if err != nil {
			return nil, err
		}
		compiled.network = network
	}

	return compiled, nil
}

// ShouldDrop returns true if the trap should be dropped
func (f *Filter) ShouldDrop(trapOID, sourceIP string) bool {
	if !f.config.Enabled {
		return false
	}

	for _, rule := range f.rules {
		matches := f.matchesRule(rule, trapOID, sourceIP)

		if matches {
			if rule.original.Action == "drop" {
				log.Debug().
					Str("type", rule.original.Type).
					Str("pattern", rule.original.Pattern).
					Str("oid", trapOID).
					Str("source", sourceIP).
					Msg("Trap dropped by filter")
				return true
			} else if rule.original.Action == "accept" {
				// Explicit accept, don't drop
				return false
			}
		}
	}

	// Default: don't drop
	return false
}

// matchesRule checks if a trap matches a filter rule
func (f *Filter) matchesRule(rule *compiledRule, trapOID, sourceIP string) bool {
	switch rule.original.Type {
	case "oid":
		// Match OID pattern
		if rule.pattern != nil {
			return rule.pattern.MatchString(trapOID)
		}
		return false

	case "source":
		// Exact IP match
		return sourceIP == rule.original.Pattern

	case "source_network":
		// Network CIDR match
		if rule.network != nil {
			ip := net.ParseIP(sourceIP)
			if ip != nil {
				return rule.network.Contains(ip)
			}
		}
		return false

	default:
		return false
	}
}

// Stats returns filter statistics
func (f *Filter) Stats() map[string]interface{} {
	return map[string]interface{}{
		"enabled": f.config.Enabled,
		"rules":   len(f.rules),
	}
}
