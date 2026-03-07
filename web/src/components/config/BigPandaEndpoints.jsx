import React, { useState, useEffect } from 'react';
import { Plus, Edit2, Trash2, Save, X, Eye, EyeOff, Globe, TestTube, CheckCircle } from 'lucide-react';
import api from '../../api';

function BigPandaEndpoints({ onSave }) {
  const [endpoints, setEndpoints] = useState([]);
  const [loading, setLoading] = useState(true);
  const [editingEndpoint, setEditingEndpoint] = useState(null);
  const [showForm, setShowForm] = useState(false);
  const [showTokens, setShowTokens] = useState({});

  useEffect(() => {
    loadEndpoints();
  }, []);

  const loadEndpoints = async () => {
    try {
      const data = await api.getBigPandaEndpoints();
      setEndpoints(data?.endpoints || data || []);
      setLoading(false);
    } catch (err) {
      console.error('Failed to load BigPanda endpoints:', err);
      setLoading(false);
    }
  };

  const handleAddEndpoint = () => {
    setEditingEndpoint({
      name: '',
      description: '',
      enabled: true,
      api_url: 'https://integrations.bigpanda.io/oim/api/alerts',
      stream_url: '',
      heartbeat_url: '',
      token: '',
      app_key: '',
      batching: {
        enabled: true,
        max_size: 100,
        max_wait: '5s',
        max_bytes: 1048576,
      },
      retry: {
        max_attempts: 5,
        initial_backoff: '1s',
        max_backoff: '60s',
        backoff_multiplier: 2.0,
      },
      rate_limit: {
        events_per_second: 1000,
        burst: 2000,
      },
      timeout: {
        connect: '10s',
        request: '30s',
        idle: '90s',
      },
      tags: {},
    });
    setShowForm(true);
  };

  const handleEditEndpoint = (endpoint) => {
    setEditingEndpoint({ ...endpoint });
    setShowForm(true);
  };

  const handleSaveEndpoint = async () => {
    // Update or add endpoint to the list
    let updatedEndpoints;
    const existingIndex = endpoints.findIndex(e => e.name === editingEndpoint.name);

    if (existingIndex >= 0) {
      updatedEndpoints = [...endpoints];
      updatedEndpoints[existingIndex] = editingEndpoint;
    } else {
      updatedEndpoints = [...endpoints, editingEndpoint];
    }

    setEndpoints(updatedEndpoints);
    setShowForm(false);
    setEditingEndpoint(null);

    // Save to backend
    await onSave({ endpoints: updatedEndpoints });
  };

  const handleDeleteEndpoint = async (name) => {
    if (!confirm(`Delete endpoint "${name}"? This will affect routing for plugins using this endpoint.`)) {
      return;
    }

    const updatedEndpoints = endpoints.filter(e => e.name !== name);
    setEndpoints(updatedEndpoints);
    await onSave({ endpoints: updatedEndpoints });
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

  const toggleShowToken = (endpointName) => {
    setShowTokens(prev => ({ ...prev, [endpointName]: !prev[endpointName] }));
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
    <div>
      {/* Endpoints List */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">BigPanda Endpoints</div>
            <div className="card-subtitle">
              {endpoints.length} endpoint{endpoints.length !== 1 ? 's' : ''} configured
            </div>
          </div>
          <button onClick={handleAddEndpoint} className="btn btn-primary">
            <Plus size={18} />
            Add Endpoint
          </button>
        </div>

        {endpoints.length > 0 ? (
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            {endpoints.map((endpoint, index) => (
              <div
                key={index}
                style={{
                  padding: '20px',
                  backgroundColor: 'var(--bg-tertiary)',
                  borderRadius: '8px',
                  border: endpoint.enabled ? '2px solid var(--success)' : '2px solid var(--border)',
                }}
              >
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start', marginBottom: '16px' }}>
                  <div style={{ flex: 1 }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '8px' }}>
                      <h3 style={{ fontSize: '18px', fontWeight: 600, margin: 0 }}>
                        {endpoint.name}
                      </h3>
                      {endpoint.enabled ? (
                        <span className="badge badge-success">
                          <CheckCircle size={14} style={{ marginRight: '4px' }} />
                          Enabled
                        </span>
                      ) : (
                        <span className="badge badge-error">Disabled</span>
                      )}
                      {endpoint.name === 'default' && (
                        <span className="badge badge-info">Default</span>
                      )}
                    </div>
                    {endpoint.description && (
                      <div style={{ fontSize: '14px', color: 'var(--text-secondary)', marginBottom: '12px' }}>
                        {endpoint.description}
                      </div>
                    )}
                    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '12px', fontSize: '14px' }}>
                      <div>
                        <div style={{ color: 'var(--text-tertiary)', fontSize: '12px', marginBottom: '4px' }}>
                          API URL
                        </div>
                        <code style={{ fontSize: '12px' }}>{endpoint.api_url}</code>
                      </div>
                      <div>
                        <div style={{ color: 'var(--text-tertiary)', fontSize: '12px', marginBottom: '4px' }}>
                          App Key
                        </div>
                        <code style={{ fontSize: '12px' }}>{endpoint.app_key}</code>
                      </div>
                    </div>
                    {endpoint.tags && Object.keys(endpoint.tags).length > 0 && (
                      <div style={{ marginTop: '12px', display: 'flex', flexWrap: 'wrap', gap: '6px' }}>
                        {Object.entries(endpoint.tags).map(([key, value]) => (
                          <span key={key} className="badge badge-info" style={{ fontSize: '11px' }}>
                            {key}: {value}
                          </span>
                        ))}
                      </div>
                    )}
                  </div>
                  <div style={{ display: 'flex', gap: '8px' }}>
                    <button
                      onClick={() => handleEditEndpoint(endpoint)}
                      className="btn btn-secondary"
                      style={{ padding: '8px 12px' }}
                    >
                      <Edit2 size={16} />
                    </button>
                    <button
                      onClick={() => handleDeleteEndpoint(endpoint.name)}
                      className="btn btn-secondary"
                      style={{ padding: '8px 12px' }}
                      disabled={endpoint.name === 'default' && endpoints.length === 1}
                    >
                      <Trash2 size={16} />
                    </button>
                  </div>
                </div>

                {/* Quick Stats */}
                <div style={{
                  display: 'grid',
                  gridTemplateColumns: 'repeat(4, 1fr)',
                  gap: '12px',
                  paddingTop: '12px',
                  borderTop: '1px solid var(--border)',
                  fontSize: '12px',
                }}>
                  <div>
                    <div style={{ color: 'var(--text-tertiary)' }}>Batching</div>
                    <div style={{ fontWeight: 600 }}>
                      {endpoint.batching?.enabled ? `${endpoint.batching.max_size} events` : 'Disabled'}
                    </div>
                  </div>
                  <div>
                    <div style={{ color: 'var(--text-tertiary)' }}>Retry</div>
                    <div style={{ fontWeight: 600 }}>
                      {endpoint.retry?.max_attempts || 5} attempts
                    </div>
                  </div>
                  <div>
                    <div style={{ color: 'var(--text-tertiary)' }}>Rate Limit</div>
                    <div style={{ fontWeight: 600 }}>
                      {endpoint.rate_limit?.events_per_second || 1000}/sec
                    </div>
                  </div>
                  <div>
                    <div style={{ color: 'var(--text-tertiary)' }}>Timeout</div>
                    <div style={{ fontWeight: 600 }}>
                      {endpoint.timeout?.request || '30s'}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="empty-state">
            <div className="empty-state-icon">
              <Globe size={48} />
            </div>
            <div>No BigPanda endpoints configured</div>
            <div style={{ marginTop: '12px', fontSize: '14px', color: 'var(--text-secondary)' }}>
              Add at least one BigPanda endpoint to start sending events
            </div>
          </div>
        )}
      </div>

      {/* Endpoint Form Modal */}
      {showForm && editingEndpoint && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0, 0, 0, 0.8)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 1000,
          padding: '20px',
          overflow: 'auto',
        }}>
          <div className="card" style={{
            maxWidth: '900px',
            width: '100%',
            maxHeight: '90vh',
            overflow: 'auto',
            margin: 'auto',
          }}>
            <div className="card-header">
              <div className="card-title">
                {endpoints.some(e => e.name === editingEndpoint.name) ? 'Edit' : 'Add'} BigPanda Endpoint
              </div>
              <button onClick={() => setShowForm(false)} className="btn btn-secondary" style={{ padding: '8px' }}>
                <X size={18} />
              </button>
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
              {/* Basic Info */}
              <div>
                <h4 style={{ fontSize: '16px', fontWeight: 600, marginBottom: '16px' }}>Basic Information</h4>
                <div className="grid grid-2" style={{ gap: '16px' }}>
                  <div>
                    <label className="form-label">Endpoint Name *</label>
                    <input
                      type="text"
                      className="form-input"
                      value={editingEndpoint.name || ''}
                      onChange={(e) => updateEndpointField('name', e.target.value)}
                      placeholder="prod-network, test-all, etc."
                      required
                      disabled={endpoints.some(e => e.name === editingEndpoint.name)}
                    />
                    <div className="form-help">Unique identifier for this endpoint</div>
                  </div>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '12px', paddingTop: '28px' }}>
                    <input
                      type="checkbox"
                      id="endpoint-enabled"
                      checked={editingEndpoint.enabled || false}
                      onChange={(e) => updateEndpointField('enabled', e.target.checked)}
                    />
                    <label htmlFor="endpoint-enabled" style={{ cursor: 'pointer', margin: 0 }}>
                      Enabled
                    </label>
                  </div>
                </div>

                <div style={{ marginTop: '16px' }}>
                  <label className="form-label">Description</label>
                  <input
                    type="text"
                    className="form-input"
                    value={editingEndpoint.description || ''}
                    onChange={(e) => updateEndpointField('description', e.target.value)}
                    placeholder="Production BigPanda - Network Team"
                  />
                </div>
              </div>

              {/* API Configuration */}
              <div>
                <h4 style={{ fontSize: '16px', fontWeight: 600, marginBottom: '16px' }}>API Configuration</h4>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                  <div>
                    <label className="form-label">API URL *</label>
                    <select
                      className="form-input"
                      value={editingEndpoint.api_url || ''}
                      onChange={(e) => updateEndpointField('api_url', e.target.value)}
                    >
                      <option value="https://integrations.bigpanda.io/oim/api/alerts">US - integrations.bigpanda.io</option>
                      <option value="https://eu.integrations.bigpanda.io/oim/api/alerts">EU - eu.integrations.bigpanda.io</option>
                      <option value="">Custom...</option>
                    </select>
                    {editingEndpoint.api_url === '' && (
                      <input
                        type="text"
                        className="form-input"
                        style={{ marginTop: '8px' }}
                        placeholder="https://your-bigpanda-url.com/data/v2/alerts"
                        onChange={(e) => updateEndpointField('api_url', e.target.value)}
                      />
                    )}
                  </div>

                  <div className="grid grid-2" style={{ gap: '16px' }}>
                    <div>
                      <label className="form-label">API Token *</label>
                      <div style={{ position: 'relative' }}>
                        <input
                          type={showTokens[editingEndpoint.name] ? 'text' : 'password'}
                          className="form-input"
                          style={{ paddingRight: '40px' }}
                          value={editingEndpoint.token || ''}
                          onChange={(e) => updateEndpointField('token', e.target.value)}
                          placeholder="BigPanda API token"
                          required
                        />
                        <button
                          type="button"
                          onClick={() => toggleShowToken(editingEndpoint.name)}
                          style={{
                            position: 'absolute',
                            right: '10px',
                            top: '50%',
                            transform: 'translateY(-50%)',
                            background: 'none',
                            border: 'none',
                            cursor: 'pointer',
                            color: 'var(--text-tertiary)',
                          }}
                        >
                          {showTokens[editingEndpoint.name] ? <EyeOff size={18} /> : <Eye size={18} />}
                        </button>
                      </div>
                    </div>
                    <div>
                      <label className="form-label">App Key *</label>
                      <input
                        type="text"
                        className="form-input"
                        value={editingEndpoint.app_key || ''}
                        onChange={(e) => updateEndpointField('app_key', e.target.value)}
                        placeholder="Integration app key"
                        required
                      />
                      <div className="form-help">Identifies this integration in BigPanda</div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Batching */}
              <div>
                <h4 style={{ fontSize: '16px', fontWeight: 600, marginBottom: '16px' }}>Batching Configuration</h4>
                <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '16px' }}>
                  <input
                    type="checkbox"
                    id="batching-enabled"
                    checked={editingEndpoint.batching?.enabled || false}
                    onChange={(e) => updateEndpointField('batching.enabled', e.target.checked)}
                  />
                  <label htmlFor="batching-enabled" style={{ cursor: 'pointer', margin: 0 }}>
                    Enable event batching
                  </label>
                </div>
                {editingEndpoint.batching?.enabled && (
                  <div className="grid grid-3" style={{ gap: '16px' }}>
                    <div>
                      <label className="form-label">Max Size</label>
                      <input
                        type="number"
                        className="form-input"
                        value={editingEndpoint.batching?.max_size || 100}
                        onChange={(e) => updateEndpointField('batching.max_size', parseInt(e.target.value))}
                        min="1"
                      />
                      <div className="form-help">Events per batch</div>
                    </div>
                    <div>
                      <label className="form-label">Max Wait</label>
                      <input
                        type="text"
                        className="form-input"
                        value={editingEndpoint.batching?.max_wait || '5s'}
                        onChange={(e) => updateEndpointField('batching.max_wait', e.target.value)}
                        placeholder="5s"
                      />
                    </div>
                    <div>
                      <label className="form-label">Max Bytes</label>
                      <input
                        type="number"
                        className="form-input"
                        value={editingEndpoint.batching?.max_bytes || 1048576}
                        onChange={(e) => updateEndpointField('batching.max_bytes', parseInt(e.target.value))}
                      />
                    </div>
                  </div>
                )}
              </div>

              {/* Retry */}
              <div>
                <h4 style={{ fontSize: '16px', fontWeight: 600, marginBottom: '16px' }}>Retry Configuration</h4>
                <div className="grid grid-4" style={{ gap: '16px' }}>
                  <div>
                    <label className="form-label">Max Attempts</label>
                    <input
                      type="number"
                      className="form-input"
                      value={editingEndpoint.retry?.max_attempts || 5}
                      onChange={(e) => updateEndpointField('retry.max_attempts', parseInt(e.target.value))}
                      min="1"
                      max="10"
                    />
                  </div>
                  <div>
                    <label className="form-label">Initial Backoff</label>
                    <input
                      type="text"
                      className="form-input"
                      value={editingEndpoint.retry?.initial_backoff || '1s'}
                      onChange={(e) => updateEndpointField('retry.initial_backoff', e.target.value)}
                    />
                  </div>
                  <div>
                    <label className="form-label">Max Backoff</label>
                    <input
                      type="text"
                      className="form-input"
                      value={editingEndpoint.retry?.max_backoff || '60s'}
                      onChange={(e) => updateEndpointField('retry.max_backoff', e.target.value)}
                    />
                  </div>
                  <div>
                    <label className="form-label">Multiplier</label>
                    <input
                      type="number"
                      className="form-input"
                      value={editingEndpoint.retry?.backoff_multiplier || 2.0}
                      onChange={(e) => updateEndpointField('retry.backoff_multiplier', parseFloat(e.target.value))}
                      step="0.1"
                    />
                  </div>
                </div>
              </div>

              {/* Rate Limit & Timeout */}
              <div className="grid grid-2" style={{ gap: '24px' }}>
                <div>
                  <h4 style={{ fontSize: '16px', fontWeight: 600, marginBottom: '16px' }}>Rate Limiting</h4>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                    <div>
                      <label className="form-label">Events/Second</label>
                      <input
                        type="number"
                        className="form-input"
                        value={editingEndpoint.rate_limit?.events_per_second || 1000}
                        onChange={(e) => updateEndpointField('rate_limit.events_per_second', parseInt(e.target.value))}
                        min="1"
                      />
                    </div>
                    <div>
                      <label className="form-label">Burst</label>
                      <input
                        type="number"
                        className="form-input"
                        value={editingEndpoint.rate_limit?.burst || 2000}
                        onChange={(e) => updateEndpointField('rate_limit.burst', parseInt(e.target.value))}
                        min="1"
                      />
                    </div>
                  </div>
                </div>
                <div>
                  <h4 style={{ fontSize: '16px', fontWeight: 600, marginBottom: '16px' }}>Timeouts</h4>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                    <div>
                      <label className="form-label">Connect</label>
                      <input
                        type="text"
                        className="form-input"
                        value={editingEndpoint.timeout?.connect || '10s'}
                        onChange={(e) => updateEndpointField('timeout.connect', e.target.value)}
                      />
                    </div>
                    <div>
                      <label className="form-label">Request</label>
                      <input
                        type="text"
                        className="form-input"
                        value={editingEndpoint.timeout?.request || '30s'}
                        onChange={(e) => updateEndpointField('timeout.request', e.target.value)}
                      />
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Save/Cancel Buttons */}
            <div style={{ display: 'flex', gap: '12px', marginTop: '24px', justifyContent: 'flex-end' }}>
              <button onClick={() => setShowForm(false)} className="btn btn-secondary">
                Cancel
              </button>
              <button
                onClick={handleSaveEndpoint}
                className="btn btn-primary"
                disabled={!editingEndpoint.name || !editingEndpoint.api_url || !editingEndpoint.token || !editingEndpoint.app_key}
              >
                <Save size={18} />
                Save Endpoint
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default BigPandaEndpoints;
