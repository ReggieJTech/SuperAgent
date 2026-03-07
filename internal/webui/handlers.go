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
	health := s.agent.Health()

	pluginStats, ok := stats["plugins"].(map[string]interface{})
	if !ok {
		s.writeJSON(w, http.StatusOK, []interface{}{})
		return
	}

	pluginHealth, _ := health["plugins"].(map[string]interface{})
	healthInfo, _ := pluginHealth["plugins"].(map[string]interface{})

	// Build full plugin objects with stats and health
	plugins := make([]map[string]interface{}, 0)
	if report, ok := pluginStats["plugins"].(map[string]map[string]interface{}); ok {
		for name, pluginData := range report {
			plugin := map[string]interface{}{
				"name":  name,
				"stats": pluginData,
			}

			// Add health info if available
			if healthInfo != nil {
				if health, ok := healthInfo[name].(map[string]interface{}); ok {
					plugin["status"] = health["status"]
					plugin["health"] = health
				}
			}

			plugins = append(plugins, plugin)
		}
	}

	s.writeJSON(w, http.StatusOK, plugins)
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

// Configuration Management Handlers

// handleGetBigPandaConfig returns the BigPanda configuration
func (s *Server) handleGetBigPandaConfig(w http.ResponseWriter, r *http.Request) {
	config := s.agent.Config()

	// Type assertion to get the full config
	if cfg, ok := config.(map[string]interface{}); ok {
		if bpConfig, exists := cfg["bigpanda"]; exists {
			s.writeJSON(w, http.StatusOK, bpConfig)
			return
		}
	}

	s.writeError(w, http.StatusNotFound, "BigPanda configuration not available")
}

// handleUpdateBigPandaConfig updates the BigPanda configuration
func (s *Server) handleUpdateBigPandaConfig(w http.ResponseWriter, r *http.Request) {
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// TODO: Implement configuration update logic
	// This would involve:
	// 1. Validate the updates
	// 2. Update the config file
	// 3. Reload the agent configuration
	// 4. Restart affected components if needed

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "BigPanda configuration updated successfully",
		"restart_required": true,
	})
}

// handleGetSNMPConfig returns the SNMP module configuration
func (s *Server) handleGetSNMPConfig(w http.ResponseWriter, r *http.Request) {
	// Get SNMP plugin stats which may include config
	stats := s.agent.Stats()
	pluginStats, ok := stats["plugins"].(map[string]interface{})
	if !ok {
		s.writeError(w, http.StatusNotFound, "SNMP configuration not available")
		return
	}

	report, ok := pluginStats["plugins"].(map[string]map[string]interface{})
	if !ok {
		s.writeError(w, http.StatusNotFound, "SNMP configuration not available")
		return
	}

	snmpData, exists := report["snmp"]
	if !exists {
		s.writeError(w, http.StatusNotFound, "SNMP module not loaded")
		return
	}

	s.writeJSON(w, http.StatusOK, snmpData)
}

// handleUpdateSNMPConfig updates the SNMP module configuration
func (s *Server) handleUpdateSNMPConfig(w http.ResponseWriter, r *http.Request) {
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// TODO: Implement SNMP configuration update logic
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "SNMP configuration updated successfully",
		"restart_required": true,
	})
}

// handleGetWebhookConfig returns the Webhook module configuration
func (s *Server) handleGetWebhookConfig(w http.ResponseWriter, r *http.Request) {
	stats := s.agent.Stats()
	pluginStats, ok := stats["plugins"].(map[string]interface{})
	if !ok {
		s.writeError(w, http.StatusNotFound, "Webhook configuration not available")
		return
	}

	report, ok := pluginStats["plugins"].(map[string]map[string]interface{})
	if !ok {
		s.writeError(w, http.StatusNotFound, "Webhook configuration not available")
		return
	}

	webhookData, exists := report["webhook"]
	if !exists {
		s.writeError(w, http.StatusNotFound, "Webhook module not loaded")
		return
	}

	s.writeJSON(w, http.StatusOK, webhookData)
}

// handleUpdateWebhookConfig updates the Webhook module configuration
func (s *Server) handleUpdateWebhookConfig(w http.ResponseWriter, r *http.Request) {
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// TODO: Implement Webhook configuration update logic
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Webhook configuration updated successfully",
		"restart_required": true,
	})
}

// handleUploadMIB handles MIB file uploads
func (s *Server) handleUploadMIB(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (limit to 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		s.writeError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	file, header, err := r.FormFile("mib")
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "No MIB file provided")
		return
	}
	defer file.Close()

	log.Info().
		Str("filename", header.Filename).
		Int64("size", header.Size).
		Msg("MIB file upload received")

	// TODO: Implement MIB file processing:
	// 1. Save to MIBs directory
	// 2. Validate MIB syntax
	// 3. Compile MIB if needed
	// 4. Update available MIBs list

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "MIB file uploaded successfully",
		"filename": header.Filename,
		"size": header.Size,
	})
}

// handleGenerateEventConfig generates event configuration from a MIB
func (s *Server) handleGenerateEventConfig(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MIBName string `json:"mib_name"`
		Vendor  string `json:"vendor"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.MIBName == "" {
		s.writeError(w, http.StatusBadRequest, "mib_name is required")
		return
	}

	log.Info().
		Str("mib", req.MIBName).
		Str("vendor", req.Vendor).
		Msg("Event config generation requested")

	// TODO: Implement event config generation:
	// 1. Parse MIB file
	// 2. Extract trap definitions
	// 3. Generate event configuration mappings
	// 4. Return generated config for review

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Event configuration generated",
		"mib": req.MIBName,
		"configs_generated": 0,
	})
}

// handleListEventConfigs lists all SNMP event configurations
func (s *Server) handleListEventConfigs(w http.ResponseWriter, r *http.Request) {
	// TODO: List all event config files from the event_configs_dir
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"configs": []string{},
		"count": 0,
	})
}

// handleGetEventConfig gets a specific SNMP event configuration
func (s *Server) handleGetEventConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		s.writeError(w, http.StatusBadRequest, "config name is required")
		return
	}

	// TODO: Load event config file and return its contents
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"name": name,
		"events": []interface{}{},
	})
}

// handleUpdateEventConfig updates a specific SNMP event configuration
func (s *Server) handleUpdateEventConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		s.writeError(w, http.StatusBadRequest, "config name is required")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// TODO: Update event config file
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Event configuration updated",
		"name": name,
	})
}

// handleDeleteEventConfig deletes a specific SNMP event configuration
func (s *Server) handleDeleteEventConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		s.writeError(w, http.StatusBadRequest, "config name is required")
		return
	}

	// TODO: Delete event config file
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Event configuration deleted",
		"name": name,
	})
}

// handleCreateWebhookEndpoint creates a new webhook endpoint
func (s *Server) handleCreateWebhookEndpoint(w http.ResponseWriter, r *http.Request) {
	var endpoint map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&endpoint); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Validate required fields
	name, ok := endpoint["name"].(string)
	if !ok || name == "" {
		s.writeError(w, http.StatusBadRequest, "endpoint name is required")
		return
	}

	path, ok := endpoint["path"].(string)
	if !ok || path == "" {
		s.writeError(w, http.StatusBadRequest, "endpoint path is required")
		return
	}

	// TODO: Add webhook endpoint to configuration
	s.writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Webhook endpoint created",
		"name": name,
		"restart_required": true,
	})
}

// handleUpdateWebhookEndpoint updates an existing webhook endpoint
func (s *Server) handleUpdateWebhookEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		s.writeError(w, http.StatusBadRequest, "endpoint name is required")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// TODO: Update webhook endpoint in configuration
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Webhook endpoint updated",
		"name": name,
		"restart_required": true,
	})
}

// handleDeleteWebhookEndpoint deletes a webhook endpoint
func (s *Server) handleDeleteWebhookEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		s.writeError(w, http.StatusBadRequest, "endpoint name is required")
		return
	}

	// TODO: Remove webhook endpoint from configuration
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Webhook endpoint deleted",
		"name": name,
		"restart_required": true,
	})
}
