const API_BASE = '';

class ApiClient {
  constructor() {
    this.token = localStorage.getItem('auth_token');
  }

  async fetch(endpoint, options = {}) {
    const headers = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
      ...options,
      headers,
    });

    if (!response.ok) {
      if (response.status === 401) {
        this.token = null;
        localStorage.removeItem('auth_token');
        window.location.href = '/login';
      }
      throw new Error(`API error: ${response.status} ${response.statusText}`);
    }

    const json = await response.json();
    // Backend wraps responses in {success, data, time} - extract the data
    return json.data || json;
  }

  // Health endpoints
  async getHealth() {
    return this.fetch('/health');
  }

  async getHealthLive() {
    return this.fetch('/health/live');
  }

  async getHealthReady() {
    return this.fetch('/health/ready');
  }

  // Agent endpoints
  async getAgentInfo() {
    return this.fetch('/api/v1/agent/info');
  }

  async getAgentConfig() {
    return this.fetch('/api/v1/agent/config');
  }

  async getStats() {
    return this.fetch('/api/v1/stats');
  }

  // Queue endpoints
  async getQueueStats() {
    return this.fetch('/api/v1/queue/stats');
  }

  async getQueueSize() {
    return this.fetch('/api/v1/queue/size');
  }

  // Plugin endpoints
  async getPlugins() {
    return this.fetch('/api/v1/plugins');
  }

  async getPlugin(name) {
    return this.fetch(`/api/v1/plugins/${name}`);
  }

  async getPluginStats(name) {
    return this.fetch(`/api/v1/plugins/${name}/stats`);
  }

  // Event endpoints
  async getRecentEvents() {
    return this.fetch('/api/v1/events/recent');
  }

  async getDLQEvents() {
    return this.fetch('/api/v1/events/dlq');
  }

  // SNMP endpoints
  async getSNMPConfigs() {
    return this.fetch('/api/v1/snmp/configs');
  }

  async getUnknownSNMPTraps() {
    return this.fetch('/api/v1/snmp/unknown');
  }

  // Configuration Management endpoints
  // BigPanda Endpoints configuration
  async getBigPandaEndpoints() {
    return this.fetch('/api/v1/config/bigpanda/endpoints');
  }

  async updateBigPandaEndpoints(endpoints) {
    return this.fetch('/api/v1/config/bigpanda/endpoints', {
      method: 'PUT',
      body: JSON.stringify(endpoints),
    });
  }

  async createBigPandaEndpoint(endpoint) {
    return this.fetch('/api/v1/config/bigpanda/endpoints', {
      method: 'POST',
      body: JSON.stringify(endpoint),
    });
  }

  async deleteBigPandaEndpoint(name) {
    return this.fetch(`/api/v1/config/bigpanda/endpoints/${name}`, {
      method: 'DELETE',
    });
  }

  // Legacy single BigPanda config (backward compatible)
  async getBigPandaConfig() {
    return this.fetch('/api/v1/config/bigpanda');
  }

  async updateBigPandaConfig(config) {
    return this.fetch('/api/v1/config/bigpanda', {
      method: 'PUT',
      body: JSON.stringify(config),
    });
  }

  // SNMP configuration
  async getSNMPConfig() {
    return this.fetch('/api/v1/config/snmp');
  }

  async updateSNMPConfig(config) {
    return this.fetch('/api/v1/config/snmp', {
      method: 'PUT',
      body: JSON.stringify(config),
    });
  }

  // Webhook configuration
  async getWebhookConfig() {
    return this.fetch('/api/v1/config/webhook');
  }

  async updateWebhookConfig(config) {
    return this.fetch('/api/v1/config/webhook', {
      method: 'PUT',
      body: JSON.stringify(config),
    });
  }

  // SNMP MIB Management
  async uploadMIB(file) {
    const formData = new FormData();
    formData.append('mib', file);

    const response = await fetch(`${API_BASE}/api/v1/snmp/mibs/upload`, {
      method: 'POST',
      headers: {
        ...(this.token && { 'Authorization': `Bearer ${this.token}` }),
      },
      body: formData,
    });

    if (!response.ok) {
      throw new Error(`Upload failed: ${response.status} ${response.statusText}`);
    }

    const json = await response.json();
    return json.data || json;
  }

  async generateEventConfig(mibName, vendor) {
    return this.fetch('/api/v1/snmp/events/generate', {
      method: 'POST',
      body: JSON.stringify({ mib_name: mibName, vendor }),
    });
  }

  // SNMP Event Config Management
  async listEventConfigs() {
    return this.fetch('/api/v1/snmp/events');
  }

  async getEventConfig(name) {
    return this.fetch(`/api/v1/snmp/events/${name}`);
  }

  async updateEventConfig(name, config) {
    return this.fetch(`/api/v1/snmp/events/${name}`, {
      method: 'PUT',
      body: JSON.stringify(config),
    });
  }

  async deleteEventConfig(name) {
    return this.fetch(`/api/v1/snmp/events/${name}`, {
      method: 'DELETE',
    });
  }

  // Webhook Endpoint Management
  async createWebhookEndpoint(endpoint) {
    return this.fetch('/api/v1/webhook/endpoints', {
      method: 'POST',
      body: JSON.stringify(endpoint),
    });
  }

  async updateWebhookEndpoint(name, endpoint) {
    return this.fetch(`/api/v1/webhook/endpoints/${name}`, {
      method: 'PUT',
      body: JSON.stringify(endpoint),
    });
  }

  async deleteWebhookEndpoint(name) {
    return this.fetch(`/api/v1/webhook/endpoints/${name}`, {
      method: 'DELETE',
    });
  }

  // Auth endpoints
  async login(username, password) {
    const response = await this.fetch('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    });
    this.token = response.token;
    localStorage.setItem('auth_token', response.token);
    return response;
  }

  async logout() {
    await this.fetch('/api/v1/auth/logout', { method: 'POST' });
    this.token = null;
    localStorage.removeItem('auth_token');
  }

  async refreshToken() {
    const response = await this.fetch('/api/v1/auth/refresh', {
      method: 'POST',
    });
    this.token = response.token;
    localStorage.setItem('auth_token', response.token);
    return response;
  }

  // WebSocket connection for real-time events
  connectEventStream(onMessage, onError) {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/v1/events/stream`;

    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      console.log('WebSocket connected');
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        onMessage(data);
      } catch (err) {
        console.error('Failed to parse WebSocket message:', err);
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      if (onError) onError(error);
    };

    ws.onclose = () => {
      console.log('WebSocket disconnected');
    };

    return ws;
  }
}

export default new ApiClient();
