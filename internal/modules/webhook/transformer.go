package webhook

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ReggieJTech/SuperAgent/internal/queue"
)

// Transformer handles webhook payload transformation
type Transformer struct {
	config TransformConfig
}

// NewTransformer creates a new transformer
func NewTransformer(config TransformConfig) *Transformer {
	return &Transformer{
		config: config,
	}
}

// Transform converts a webhook payload to BigPanda events
func (t *Transformer) Transform(payload []byte, sourceName string) ([]*queue.Event, error) {
	// Parse JSON payload
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// If template is provided, use template-based transformation
	if t.config.Template != "" {
		return t.transformWithTemplate(data, sourceName)
	}

	// Otherwise use field mapping transformation
	return t.transformWithFieldMap(data, sourceName)
}

// transformWithFieldMap uses field mapping for transformation
func (t *Transformer) transformWithFieldMap(data map[string]interface{}, sourceName string) ([]*queue.Event, error) {
	event := &queue.Event{
		Timestamp:    time.Now(),
		CustomFields: make(map[string]interface{}),
	}

	// Map fields
	for srcField, dstField := range t.config.FieldMap {
		value := extractValue(data, srcField)
		if value != nil {
			switch dstField {
			case "host":
				event.PrimaryKey = fmt.Sprintf("%v", value)
			case "check":
				event.Check = fmt.Sprintf("%v", value)
				event.SecondaryKey = fmt.Sprintf("%v", value)
			case "status":
				event.Status = fmt.Sprintf("%v", value)
			case "description":
				event.Description = fmt.Sprintf("%v", value)
			case "severity":
				// Map severity to status if status not set
				if event.Status == "" {
					event.Status = fmt.Sprintf("%v", value)
				} else {
					event.CustomFields[dstField] = value
				}
			default:
				event.CustomFields[dstField] = value
			}
		}
	}

	// Apply status mapping
	if t.config.StatusMap != nil && event.Status != "" {
		if mappedStatus, ok := t.config.StatusMap[event.Status]; ok {
			event.Status = mappedStatus
		}
	}

	// Set additional fields
	for key, value := range t.config.Set {
		event.CustomFields[key] = value
	}

	// Add source information
	event.CustomFields["webhook_source"] = sourceName

	// Set primary and secondary keys
	if t.config.PrimaryKey != "" {
		primaryValue := extractValue(data, t.config.PrimaryKey)
		if primaryValue != nil {
			event.PrimaryKey = fmt.Sprintf("%v", primaryValue)
		}
	}

	if t.config.SecondaryKey != "" {
		secondaryValue := extractValue(data, t.config.SecondaryKey)
		if secondaryValue != nil {
			event.SecondaryKey = fmt.Sprintf("%v", secondaryValue)
			event.Check = fmt.Sprintf("%v", secondaryValue)
		}
	}

	// Default values if not set
	if event.PrimaryKey == "" {
		event.PrimaryKey = "unknown"
	}
	if event.SecondaryKey == "" {
		event.SecondaryKey = "webhook_event"
	}
	if event.Check == "" {
		event.Check = "webhook_event"
	}
	if event.Status == "" {
		event.Status = "warning"
	}

	return []*queue.Event{event}, nil
}

// transformWithTemplate uses Go template for transformation (simplified version)
func (t *Transformer) transformWithTemplate(data map[string]interface{}, sourceName string) ([]*queue.Event, error) {
	// For now, return basic transformation
	// Full template support would require text/template package
	return t.transformWithFieldMap(data, sourceName)
}

// extractValue extracts a value from nested map using dot notation (e.g., "result.host")
func extractValue(data map[string]interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	var current interface{} = data

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil
			}
		default:
			return nil
		}
	}

	return current
}
