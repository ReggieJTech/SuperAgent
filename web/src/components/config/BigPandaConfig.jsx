import React, { useState, useEffect } from 'react';
import { Save, RefreshCw, Eye, EyeOff } from 'lucide-react';
import api from '../../api';

function BigPandaConfig({ onSave }) {
  const [config, setConfig] = useState({
    api_url: '',
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
  });

  const [loading, setLoading] = useState(true);
  const [showToken, setShowToken] = useState(false);
  const [showAppKey, setShowAppKey] = useState(false);
  const [modified, setModified] = useState(false);

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    try {
      const data = await api.getBigPandaConfig();
      if (data) {
        setConfig(data);
      }
      setLoading(false);
    } catch (err) {
      console.error('Failed to load BigPanda config:', err);
      setLoading(false);
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
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">BigPanda API Configuration</div>
            <div className="card-subtitle">Configure connection to BigPanda</div>
          </div>
          {modified && (
            <span className="badge badge-warning">Unsaved Changes</span>
          )}
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          {/* API Endpoints */}
          <div>
            <label className="form-label">API URL</label>
            <input
              type="text"
              className="form-input"
              value={config.api_url || ''}
              onChange={(e) => handleChange('api_url', e.target.value)}
              placeholder="https://api.bigpanda.io/data/v2/alerts"
              required
            />
            <div className="form-help">BigPanda alerts API endpoint</div>
          </div>

          <div className="grid grid-2">
            <div>
              <label className="form-label">Stream URL</label>
              <input
                type="text"
                className="form-input"
                value={config.stream_url || ''}
                onChange={(e) => handleChange('stream_url', e.target.value)}
                placeholder="https://api.bigpanda.io/data/v2/stream"
              />
            </div>
            <div>
              <label className="form-label">Heartbeat URL</label>
              <input
                type="text"
                className="form-input"
                value={config.heartbeat_url || ''}
                onChange={(e) => handleChange('heartbeat_url', e.target.value)}
                placeholder="https://api.bigpanda.io/data/v2/heartbeat"
              />
            </div>
          </div>

          {/* Authentication */}
          <div className="grid grid-2">
            <div>
              <label className="form-label">API Token</label>
              <div style={{ position: 'relative' }}>
                <input
                  type={showToken ? 'text' : 'password'}
                  className="form-input"
                  style={{ paddingRight: '40px' }}
                  value={config.token || ''}
                  onChange={(e) => handleChange('token', e.target.value)}
                  placeholder="Enter BigPanda API token"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowToken(!showToken)}
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
                  {showToken ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
              </div>
            </div>
            <div>
              <label className="form-label">App Key</label>
              <div style={{ position: 'relative' }}>
                <input
                  type={showAppKey ? 'text' : 'password'}
                  className="form-input"
                  style={{ paddingRight: '40px' }}
                  value={config.app_key || ''}
                  onChange={(e) => handleChange('app_key', e.target.value)}
                  placeholder="Enter BigPanda app key"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowAppKey(!showAppKey)}
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
                  {showAppKey ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Batching Configuration */}
      <div className="card">
        <div className="card-header">
          <div className="card-title">Batching Configuration</div>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <input
              type="checkbox"
              id="batching-enabled"
              checked={config.batching?.enabled || false}
              onChange={(e) => handleChange('batching.enabled', e.target.checked)}
            />
            <label htmlFor="batching-enabled" style={{ cursor: 'pointer', fontSize: '14px' }}>
              Enable event batching
            </label>
          </div>

          {config.batching?.enabled && (
            <div className="grid grid-3">
              <div>
                <label className="form-label">Max Batch Size</label>
                <input
                  type="number"
                  className="form-input"
                  value={config.batching?.max_size || 100}
                  onChange={(e) => handleChange('batching.max_size', parseInt(e.target.value))}
                  min="1"
                  max="1000"
                />
                <div className="form-help">Events per batch</div>
              </div>
              <div>
                <label className="form-label">Max Wait Time</label>
                <input
                  type="text"
                  className="form-input"
                  value={config.batching?.max_wait || '5s'}
                  onChange={(e) => handleChange('batching.max_wait', e.target.value)}
                  placeholder="5s"
                />
                <div className="form-help">e.g., 5s, 1m</div>
              </div>
              <div>
                <label className="form-label">Max Bytes</label>
                <input
                  type="number"
                  className="form-input"
                  value={config.batching?.max_bytes || 1048576}
                  onChange={(e) => handleChange('batching.max_bytes', parseInt(e.target.value))}
                  min="1024"
                />
                <div className="form-help">Bytes per batch</div>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Retry Configuration */}
      <div className="card">
        <div className="card-header">
          <div className="card-title">Retry Configuration</div>
        </div>

        <div className="grid grid-4">
          <div>
            <label className="form-label">Max Attempts</label>
            <input
              type="number"
              className="form-input"
              value={config.retry?.max_attempts || 5}
              onChange={(e) => handleChange('retry.max_attempts', parseInt(e.target.value))}
              min="1"
              max="10"
            />
          </div>
          <div>
            <label className="form-label">Initial Backoff</label>
            <input
              type="text"
              className="form-input"
              value={config.retry?.initial_backoff || '1s'}
              onChange={(e) => handleChange('retry.initial_backoff', e.target.value)}
              placeholder="1s"
            />
          </div>
          <div>
            <label className="form-label">Max Backoff</label>
            <input
              type="text"
              className="form-input"
              value={config.retry?.max_backoff || '60s'}
              onChange={(e) => handleChange('retry.max_backoff', e.target.value)}
              placeholder="60s"
            />
          </div>
          <div>
            <label className="form-label">Backoff Multiplier</label>
            <input
              type="number"
              className="form-input"
              value={config.retry?.backoff_multiplier || 2.0}
              onChange={(e) => handleChange('retry.backoff_multiplier', parseFloat(e.target.value))}
              min="1"
              max="10"
              step="0.1"
            />
          </div>
        </div>
      </div>

      {/* Rate Limiting */}
      <div className="card">
        <div className="card-header">
          <div className="card-title">Rate Limiting</div>
        </div>

        <div className="grid grid-2">
          <div>
            <label className="form-label">Events Per Second</label>
            <input
              type="number"
              className="form-input"
              value={config.rate_limit?.events_per_second || 1000}
              onChange={(e) => handleChange('rate_limit.events_per_second', parseInt(e.target.value))}
              min="1"
            />
          </div>
          <div>
            <label className="form-label">Burst</label>
            <input
              type="number"
              className="form-input"
              value={config.rate_limit?.burst || 2000}
              onChange={(e) => handleChange('rate_limit.burst', parseInt(e.target.value))}
              min="1"
            />
          </div>
        </div>
      </div>

      {/* Timeout Configuration */}
      <div className="card">
        <div className="card-header">
          <div className="card-title">Timeout Configuration</div>
        </div>

        <div className="grid grid-3">
          <div>
            <label className="form-label">Connect Timeout</label>
            <input
              type="text"
              className="form-input"
              value={config.timeout?.connect || '10s'}
              onChange={(e) => handleChange('timeout.connect', e.target.value)}
              placeholder="10s"
            />
          </div>
          <div>
            <label className="form-label">Request Timeout</label>
            <input
              type="text"
              className="form-input"
              value={config.timeout?.request || '30s'}
              onChange={(e) => handleChange('timeout.request', e.target.value)}
              placeholder="30s"
            />
          </div>
          <div>
            <label className="form-label">Idle Timeout</label>
            <input
              type="text"
              className="form-input"
              value={config.timeout?.idle || '90s'}
              onChange={(e) => handleChange('timeout.idle', e.target.value)}
              placeholder="90s"
            />
          </div>
        </div>
      </div>

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

export default BigPandaConfig;
