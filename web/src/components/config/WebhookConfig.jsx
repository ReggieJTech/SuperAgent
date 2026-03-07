import React, { useState, useEffect } from 'react';
import { Save, Plus, Trash2, Edit2, RefreshCw } from 'lucide-react';
import api from '../../api';

function WebhookConfig({ onSave }) {
  const [config, setConfig] = useState({
    listen_address: '0.0.0.0:8080',
    tls: {
      enabled: false,
      cert_file: '',
      key_file: '',
      auto_generate: false,
    },
    global: {
      timeout: '30s',
      max_body_size: 10485760,
      rate_limit: {
        enabled: true,
        requests_per_second: 100,
        burst: 200,
      },
    },
    sources: [],
  });

  const [loading, setLoading] = useState(true);
  const [modified, setModified] = useState(false);
  const [editingEndpoint, setEditingEndpoint] = useState(null);
  const [showEndpointForm, setShowEndpointForm] = useState(false);
  const [availableEndpoints, setAvailableEndpoints] = useState([]);

  useEffect(() => {
    loadConfig();
    loadEndpoints();
  }, []);

  const loadConfig = async () => {
    try {
      const data = await api.getWebhookConfig();
      if (data && data.config) {
        setConfig(data.config);
      }
      setLoading(false);
    } catch (err) {
      console.error('Failed to load Webhook config:', err);
      setLoading(false);
    }
  };

  const loadEndpoints = async () => {
    try {
      const data = await api.getBigPandaEndpoints();
      const endpoints = data?.endpoints || data || [];
      setAvailableEndpoints(endpoints.filter(e => e.enabled).map(e => e.name));
    } catch (err) {
      console.error('Failed to load BigPanda endpoints:', err);
    }
  };

  const handleChange = (path, value) => {
    setModified(true);
    setConfig(prev => {
      const newConfig = { ...prev };
      const keys = path.split('.');
      let current = newConfig;

      for (let i = 0; i < keys.length - 1; i++) {
        if (!current[keys[i]]) {
          current[keys[i]] = {};
        }
        current = current[keys[i]];
      }

      current[keys[keys.length - 1]] = value;
      return newConfig;
    });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onSave(config);
    setModified(false);
  };

  const addEndpoint = () => {
    setEditingEndpoint({
      name: '',
      enabled: true,
      path: '/webhook/',
      method: 'POST',
      auth: {
        type: 'bearer',
        token: '',
      },
      allowed_ips: [],
      endpoints: availableEndpoints.length > 0 ? [availableEndpoints[0]] : ['default'],
      transform: {
        field_map: {},
        status_map: {},
        primary_key: 'host',
        secondary_key: 'check',
        set: {},
      },
    });
    setShowEndpointForm(true);
  };

  const editEndpoint = (endpoint) => {
    setEditingEndpoint({ ...endpoint });
    setShowEndpointForm(true);
  };

  const saveEndpoint = () => {
    const sources = [...(config.sources || [])];
    const existingIndex = sources.findIndex(s => s.name === editingEndpoint.name);

    if (existingIndex >= 0) {
      sources[existingIndex] = editingEndpoint;
    } else {
      sources.push(editingEndpoint);
    }

    handleChange('sources', sources);
    setEditingEndpoint(null);
    setShowEndpointForm(false);
  };

  const deleteEndpoint = (name) => {
    if (!confirm(`Delete endpoint "${name}"?`)) return;
    const sources = config.sources.filter(s => s.name !== name);
    handleChange('sources', sources);
  };

  const updateEndpointField = (path, value) => {
    setEditingEndpoint(prev => {
      const updated = { ...prev };
      const keys = path.split('.');
      let current = updated;

      for (let i = 0; i < keys.length - 1; i++) {
        if (!current[keys[i]]) {
          current[keys[i]] = {};
        }
        current = current[keys[i]];
      }

      current[keys[keys.length - 1]] = value;
      return updated;
    });
  };

  const toggleWebhookEndpoint = (bpEndpoint) => {
    setEditingEndpoint(prev => {
      const currentEndpoints = prev.endpoints || [];
      let newEndpoints;

      if (currentEndpoints.includes(bpEndpoint)) {
        // Don't allow removing if it's the last one
        if (currentEndpoints.length > 1) {
          newEndpoints = currentEndpoints.filter(e => e !== bpEndpoint);
        } else {
          return prev; // Keep at least one endpoint
        }
      } else {
        newEndpoints = [...currentEndpoints, bpEndpoint];
      }

      return { ...prev, endpoints: newEndpoints };
    });
  };

  if (loading) {
    return (
      <div className="card">
        <div className="loading">
          <div className="spinner"></div>
        </div>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit}>
      {/* Global Configuration */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">Global Webhook Configuration</div>
            <div className="card-subtitle">Configure webhook receiver settings</div>
          </div>
          {modified && (
            <span className="badge badge-warning">Unsaved Changes</span>
          )}
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          <div className="grid grid-2">
            <div>
              <label className="form-label">Listen Address</label>
              <input
                type="text"
                className="form-input"
                value={config.listen_address || '0.0.0.0:8080'}
                onChange={(e) => handleChange('listen_address', e.target.value)}
                placeholder="0.0.0.0:8080"
                required
              />
            </div>
            <div>
              <label className="form-label">Request Timeout</label>
              <input
                type="text"
                className="form-input"
                value={config.global?.timeout || '30s'}
                onChange={(e) => handleChange('global.timeout', e.target.value)}
                placeholder="30s"
              />
            </div>
          </div>

          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <input
              type="checkbox"
              id="tls-enabled"
              checked={config.tls?.enabled || false}
              onChange={(e) => handleChange('tls.enabled', e.target.checked)}
            />
            <label htmlFor="tls-enabled" style={{ cursor: 'pointer', fontSize: '14px' }}>
              Enable TLS/HTTPS
            </label>
          </div>

          {config.tls?.enabled && (
            <div className="grid grid-2">
              <div>
                <label className="form-label">Certificate File</label>
                <input
                  type="text"
                  className="form-input"
                  value={config.tls?.cert_file || ''}
                  onChange={(e) => handleChange('tls.cert_file', e.target.value)}
                  placeholder="/path/to/cert.pem"
                />
              </div>
              <div>
                <label className="form-label">Key File</label>
                <input
                  type="text"
                  className="form-input"
                  value={config.tls?.key_file || ''}
                  onChange={(e) => handleChange('tls.key_file', e.target.value)}
                  placeholder="/path/to/key.pem"
                />
              </div>
            </div>
          )}

          <div>
            <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '12px' }}>
              <input
                type="checkbox"
                id="rate-limit-enabled"
                checked={config.global?.rate_limit?.enabled || false}
                onChange={(e) => handleChange('global.rate_limit.enabled', e.target.checked)}
              />
              <label htmlFor="rate-limit-enabled" style={{ cursor: 'pointer', fontSize: '14px', fontWeight: 600 }}>
                Enable Global Rate Limiting
              </label>
            </div>

            {config.global?.rate_limit?.enabled && (
              <div className="grid grid-2">
                <div>
                  <label className="form-label">Requests Per Second</label>
                  <input
                    type="number"
                    className="form-input"
                    value={config.global?.rate_limit?.requests_per_second || 100}
                    onChange={(e) => handleChange('global.rate_limit.requests_per_second', parseInt(e.target.value))}
                    min="1"
                  />
                </div>
                <div>
                  <label className="form-label">Burst</label>
                  <input
                    type="number"
                    className="form-input"
                    value={config.global?.rate_limit?.burst || 200}
                    onChange={(e) => handleChange('global.rate_limit.burst', parseInt(e.target.value))}
                    min="1"
                  />
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Webhook Endpoints */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">Webhook Endpoints</div>
            <div className="card-subtitle">{config.sources?.length || 0} configured endpoint(s)</div>
          </div>
          <button
            type="button"
            onClick={addEndpoint}
            className="btn btn-primary"
            style={{ fontSize: '14px' }}
          >
            <Plus size={18} />
            Add Endpoint
          </button>
        </div>

        {config.sources && config.sources.length > 0 ? (
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            {config.sources.map((endpoint, index) => (
              <div key={index} style={{
                padding: '16px',
                backgroundColor: 'var(--bg-tertiary)',
                borderRadius: '8px',
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
              }}>
                <div style={{ flex: 1 }}>
                  <div style={{ fontWeight: 600, marginBottom: '4px' }}>
                    {endpoint.name}
                    {endpoint.enabled ? (
                      <span className="badge badge-success" style={{ marginLeft: '8px' }}>Active</span>
                    ) : (
                      <span className="badge badge-error" style={{ marginLeft: '8px' }}>Disabled</span>
                    )}
                  </div>
                  <code style={{ fontSize: '13px', color: 'var(--text-secondary)' }}>
                    {endpoint.method} {endpoint.path}
                  </code>
                  <div style={{ fontSize: '12px', marginTop: '8px', color: 'var(--text-tertiary)' }}>
                    Auth: {endpoint.auth?.type || 'none'} |
                    IP Whitelist: {endpoint.allowed_ips?.length > 0 ? `${endpoint.allowed_ips.length} IP(s)` : 'All IPs'}
                  </div>
                  {endpoint.endpoints && endpoint.endpoints.length > 0 && (
                    <div style={{ marginTop: '8px', display: 'flex', flexWrap: 'wrap', gap: '4px', alignItems: 'center' }}>
                      <span style={{ fontSize: '11px', color: 'var(--text-tertiary)' }}>→</span>
                      {endpoint.endpoints.map(ep => (
                        <span key={ep} className="badge badge-info" style={{ fontSize: '10px' }}>
                          {ep}
                        </span>
                      ))}
                    </div>
                  )}
                </div>
                <div style={{ display: 'flex', gap: '8px' }}>
                  <button
                    type="button"
                    onClick={() => editEndpoint(endpoint)}
                    className="btn btn-secondary"
                    style={{ padding: '8px 12px' }}
                  >
                    <Edit2 size={16} />
                  </button>
                  <button
                    type="button"
                    onClick={() => deleteEndpoint(endpoint.name)}
                    className="btn btn-secondary"
                    style={{ padding: '8px 12px' }}
                  >
                    <Trash2 size={16} />
                  </button>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="empty-state">
            <div>No webhook endpoints configured</div>
            <div style={{ marginTop: '8px', fontSize: '14px', color: 'var(--text-secondary)' }}>
              Add endpoints to start receiving webhook events
            </div>
          </div>
        )}
      </div>

      {/* Endpoint Form Modal */}
      {showEndpointForm && editingEndpoint && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0, 0, 0, 0.7)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 1000,
          padding: '20px',
        }}>
          <div className="card" style={{
            maxWidth: '800px',
            width: '100%',
            maxHeight: '90vh',
            overflow: 'auto',
          }}>
            <div className="card-header">
              <div className="card-title">
                {config.sources?.some(s => s.name === editingEndpoint.name) ? 'Edit' : 'Add'} Webhook Endpoint
              </div>
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
              <div className="grid grid-2">
                <div>
                  <label className="form-label">Endpoint Name *</label>
                  <input
                    type="text"
                    className="form-input"
                    value={editingEndpoint.name || ''}
                    onChange={(e) => updateEndpointField('name', e.target.value)}
                    placeholder="prometheus"
                    required
                  />
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '12px', paddingTop: '28px' }}>
                  <input
                    type="checkbox"
                    id="endpoint-enabled"
                    checked={editingEndpoint.enabled || false}
                    onChange={(e) => updateEndpointField('enabled', e.target.checked)}
                  />
                  <label htmlFor="endpoint-enabled" style={{ cursor: 'pointer' }}>
                    Enabled
                  </label>
                </div>
              </div>

              <div className="grid grid-2">
                <div>
                  <label className="form-label">Path *</label>
                  <input
                    type="text"
                    className="form-input"
                    value={editingEndpoint.path || ''}
                    onChange={(e) => updateEndpointField('path', e.target.value)}
                    placeholder="/webhook/prometheus"
                    required
                  />
                </div>
                <div>
                  <label className="form-label">HTTP Method</label>
                  <select
                    className="form-input"
                    value={editingEndpoint.method || 'POST'}
                    onChange={(e) => updateEndpointField('method', e.target.value)}
                  >
                    <option value="POST">POST</option>
                    <option value="PUT">PUT</option>
                    <option value="GET">GET</option>
                  </select>
                </div>
              </div>

              <div className="grid grid-2">
                <div>
                  <label className="form-label">Authentication Type</label>
                  <select
                    className="form-input"
                    value={editingEndpoint.auth?.type || 'none'}
                    onChange={(e) => updateEndpointField('auth.type', e.target.value)}
                  >
                    <option value="none">None</option>
                    <option value="bearer">Bearer Token</option>
                    <option value="apikey">API Key</option>
                    <option value="basic">Basic Auth</option>
                    <option value="hmac">HMAC</option>
                  </select>
                </div>
                <div>
                  {editingEndpoint.auth?.type === 'bearer' && (
                    <>
                      <label className="form-label">Token</label>
                      <input
                        type="password"
                        className="form-input"
                        value={editingEndpoint.auth?.token || ''}
                        onChange={(e) => updateEndpointField('auth.token', e.target.value)}
                        placeholder="Bearer token"
                      />
                    </>
                  )}
                  {editingEndpoint.auth?.type === 'apikey' && (
                    <>
                      <label className="form-label">API Key</label>
                      <input
                        type="password"
                        className="form-input"
                        value={editingEndpoint.auth?.key || ''}
                        onChange={(e) => updateEndpointField('auth.key', e.target.value)}
                        placeholder="API key"
                      />
                    </>
                  )}
                </div>
              </div>

              <div>
                <label className="form-label">IP Whitelist (one per line)</label>
                <textarea
                  className="form-input"
                  rows="3"
                  value={(editingEndpoint.allowed_ips || []).join('\n')}
                  onChange={(e) => updateEndpointField('allowed_ips', e.target.value.split('\n').filter(ip => ip.trim()))}
                  placeholder="10.0.0.0/8&#10;172.16.0.0/12"
                  style={{ fontFamily: 'monospace', fontSize: '13px' }}
                />
                <div className="form-help">CIDR notation supported (e.g., 10.0.0.0/8)</div>
              </div>

              {/* BigPanda Endpoint Routing */}
              <div>
                <label className="form-label">BigPanda Endpoints (select one or more) *</label>
                <div style={{ fontSize: '13px', color: 'var(--text-secondary)', marginBottom: '12px' }}>
                  Events from this webhook will be sent to the selected BigPanda endpoint(s)
                </div>
                {availableEndpoints.length === 0 ? (
                  <div style={{
                    padding: '12px',
                    backgroundColor: 'rgba(245, 158, 11, 0.1)',
                    borderRadius: '6px',
                    border: '1px solid var(--warning)',
                    fontSize: '13px',
                  }}>
                    No BigPanda endpoints configured. Please configure at least one endpoint in the BigPanda tab.
                  </div>
                ) : (
                  <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
                    {availableEndpoints.map(bpEndpoint => {
                      const isSelected = (editingEndpoint.endpoints || []).includes(bpEndpoint);
                      return (
                        <button
                          key={bpEndpoint}
                          type="button"
                          onClick={() => toggleWebhookEndpoint(bpEndpoint)}
                          style={{
                            padding: '8px 16px',
                            borderRadius: '6px',
                            border: isSelected ? '2px solid var(--accent-primary)' : '1px solid var(--border)',
                            backgroundColor: isSelected ? 'rgba(59, 130, 246, 0.1)' : 'transparent',
                            color: isSelected ? 'var(--accent-primary)' : 'var(--text-primary)',
                            cursor: 'pointer',
                            fontSize: '14px',
                            fontWeight: isSelected ? 600 : 400,
                          }}
                        >
                          {isSelected && '✓ '}{bpEndpoint}
                        </button>
                      );
                    })}
                  </div>
                )}
              </div>

              <div>
                <label className="form-label">Primary Key Field</label>
                <input
                  type="text"
                  className="form-input"
                  value={editingEndpoint.transform?.primary_key || ''}
                  onChange={(e) => updateEndpointField('transform.primary_key', e.target.value)}
                  placeholder="host"
                />
              </div>

              <div>
                <label className="form-label">Secondary Key Field</label>
                <input
                  type="text"
                  className="form-input"
                  value={editingEndpoint.transform?.secondary_key || ''}
                  onChange={(e) => updateEndpointField('transform.secondary_key', e.target.value)}
                  placeholder="check"
                />
              </div>
            </div>

            <div style={{ display: 'flex', gap: '12px', marginTop: '24px', justifyContent: 'flex-end' }}>
              <button
                type="button"
                onClick={() => {
                  setEditingEndpoint(null);
                  setShowEndpointForm(false);
                }}
                className="btn btn-secondary"
              >
                Cancel
              </button>
              <button
                type="button"
                onClick={saveEndpoint}
                className="btn btn-primary"
                disabled={!editingEndpoint.name || !editingEndpoint.path}
              >
                <Save size={18} />
                Save Endpoint
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Save Button */}
      <div className="card" style={{ padding: '16px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div style={{ fontSize: '14px', color: 'var(--text-secondary)' }}>
            {modified ? 'You have unsaved changes' : 'All changes saved'}
          </div>
          <div style={{ display: 'flex', gap: '12px' }}>
            <button
              type="button"
              onClick={loadConfig}
              className="btn btn-secondary"
              disabled={!modified}
            >
              <RefreshCw size={18} />
              Reset
            </button>
            <button
              type="submit"
              className="btn btn-primary"
              disabled={!modified}
            >
              <Save size={18} />
              Save Configuration
            </button>
          </div>
        </div>
      </div>
    </form>
  );
}

export default WebhookConfig;
