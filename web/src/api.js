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
