package webui

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Time    string      `json:"time"`
}

// writeJSON writes a JSON response
func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success: status >= 200 && status < 300,
		Data:    data,
		Time:    time.Now().UTC().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Err(err).Msg("Failed to write JSON response")
	}
}

// writeError writes an error response
func (s *Server) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success: false,
		Error:   message,
		Time:    time.Now().UTC().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

// handleHealth handles the health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := s.agent.Health()
	s.writeJSON(w, http.StatusOK, health)
}

// handleHealthLive handles the liveness probe
func (s *Server) handleHealthLive(w http.ResponseWriter, r *http.Request) {
	// Always return OK if the server is responding
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "alive",
	})
}

// handleHealthReady handles the readiness probe
func (s *Server) handleHealthReady(w http.ResponseWriter, r *http.Request) {
	if !s.agent.IsStarted() {
		s.writeError(w, http.StatusServiceUnavailable, "agent not started")
		return
	}

	health := s.agent.Health()
	status, ok := health["status"].(string)
	if !ok || status != "healthy" {
		s.writeError(w, http.StatusServiceUnavailable, "agent not ready")
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ready",
	})
}

// handleStats handles the stats endpoint
func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats := s.agent.Stats()
	s.writeJSON(w, http.StatusOK, stats)
}

// handleAgentInfo handles the agent info endpoint
func (s *Server) handleAgentInfo(w http.ResponseWriter, r *http.Request) {
	info := map[string]interface{}{
		"started": s.agent.IsStarted(),
		"uptime":  s.agent.Uptime().String(),
	}
	s.writeJSON(w, http.StatusOK, info)
}

// handleAgentConfig handles the agent config endpoint
func (s *Server) handleAgentConfig(w http.ResponseWriter, r *http.Request) {
	config := s.agent.Config()
	s.writeJSON(w, http.StatusOK, config)
}

// handleQueueStats handles the queue stats endpoint
func (s *Server) handleQueueStats(w http.ResponseWriter, r *http.Request) {
	stats := s.agent.Stats()
	queueStats, ok := stats["queue"].(map[string]interface{})
	if !ok {
		s.writeError(w, http.StatusNotFound, "queue stats not available")
		return
	}
	s.writeJSON(w, http.StatusOK, queueStats)
}

// handleQueueSize handles the queue size endpoint
func (s *Server) handleQueueSize(w http.ResponseWriter, r *http.Request) {
	stats := s.agent.Stats()
	queueStats, ok := stats["queue"].(map[string]interface{})
	if !ok {
		s.writeError(w, http.StatusNotFound, "queue stats not available")
		return
	}

	size, _ := queueStats["size"].(int)
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"size": size,
	})
}

// handlePluginsList handles the plugins list endpoint
func (s *Server) handlePluginsList(w http.ResponseWriter, r *http.Request) {
	stats := s.agent.Stats()
	pluginStats, ok := stats["plugins"].(map[string]interface{})
	if !ok {
		s.writeJSON(w, http.StatusOK, map[string]interface{}{"plugins": []string{}})
		return
	}

	// Extract plugin names
	plugins := make([]string, 0)
	if report, ok := pluginStats["plugins"].(map[string]map[string]interface{}); ok {
		for name := range report {
			plugins = append(plugins, name)
		}
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"plugins": plugins,
	})
}

// handlePluginInfo handles the plugin info endpoint
func (s *Server) handlePluginInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	// Get plugin health from agent health
	health := s.agent.Health()
	pluginHealth, ok := health["plugins"].(map[string]interface{})
	if !ok {
		s.writeError(w, http.StatusNotFound, "plugin not found")
		return
	}

	healthInfo, ok := pluginHealth["plugins"].(map[string]interface{})
	if !ok {
		s.writeError(w, http.StatusNotFound, "plugin not found")
		return
	}

	pluginInfo, ok := healthInfo[name]
	if !ok {
		s.writeError(w, http.StatusNotFound, "plugin not found")
		return
	}

	s.writeJSON(w, http.StatusOK, pluginInfo)
}

// handlePluginStats handles the plugin stats endpoint
func (s *Server) handlePluginStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	stats := s.agent.Stats()
	pluginStats, ok := stats["plugins"].(map[string]interface{})
	if !ok {
		s.writeError(w, http.StatusNotFound, "plugin stats not available")
		return
	}

	report, ok := pluginStats["plugins"].(map[string]map[string]interface{})
	if !ok {
		s.writeError(w, http.StatusNotFound, "plugin stats not available")
		return
	}

	pluginInfo, ok := report[name]
	if !ok {
		s.writeError(w, http.StatusNotFound, "plugin not found")
		return
	}

	s.writeJSON(w, http.StatusOK, pluginInfo)
}

// handleRecentEvents handles the recent events endpoint
func (s *Server) handleRecentEvents(w http.ResponseWriter, r *http.Request) {
	// This would query recent events from the queue or a cache
	// For now, return empty array
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"events": []interface{}{},
		"count":  0,
	})
}

// handleDLQEvents handles the DLQ events endpoint
func (s *Server) handleDLQEvents(w http.ResponseWriter, r *http.Request) {
	// This would query DLQ events
	// For now, return empty array
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"events": []interface{}{},
		"count":  0,
	})
}

// handleSNMPConfigs handles the SNMP configs endpoint
func (s *Server) handleSNMPConfigs(w http.ResponseWriter, r *http.Request) {
	// This would query SNMP event configs
	// For now, return placeholder
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"configs": []string{},
		"count":   0,
	})
}

// handleSNMPUnknown handles the unknown traps endpoint
func (s *Server) handleSNMPUnknown(w http.ResponseWriter, r *http.Request) {
	// This would query unknown traps
	// For now, return placeholder
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"traps": []interface{}{},
		"count": 0,
	})
}

// handleEventStream handles the WebSocket event stream
func (s *Server) handleEventStream(w http.ResponseWriter, r *http.Request) {
	// WebSocket implementation would go here
	// For now, return not implemented
	s.writeError(w, http.StatusNotImplemented, "WebSocket support coming soon")
}
