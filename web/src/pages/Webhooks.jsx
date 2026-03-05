import React, { useState, useEffect } from 'react';
import { Webhook, Activity, Shield, MapPin } from 'lucide-react';
import api from '../api';

function Webhooks() {
  const [pluginData, setPluginData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadData();
    const interval = setInterval(loadData, 5000);
    return () => clearInterval(interval);
  }, []);

  const loadData = async () => {
    try {
      // Try to get webhook plugin data
      const data = await api.getPlugin('webhook');
      setPluginData(data);
      setLoading(false);
      setError(null);
    } catch (err) {
      // If webhook plugin doesn't exist or isn't running
      if (err.message.includes('404') || err.message.includes('not found')) {
        setError('Webhook plugin is not enabled or configured');
      } else {
        setError(err.message);
      }
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="loading">
        <div className="spinner"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div>
        <div className="page-header">
          <h2>Webhooks</h2>
          <p>HTTP/HTTPS webhook receiver configuration</p>
        </div>
        <div className="error-message">
          {error}
        </div>
        <div className="card">
          <div className="empty-state">
            <div className="empty-state-icon">
              <Webhook size={48} />
            </div>
            <div>Webhook module is not currently active</div>
            <div style={{ marginTop: '12px', fontSize: '14px', color: 'var(--text-secondary)' }}>
              Enable the webhook module in your configuration to start receiving HTTP events
            </div>
          </div>
        </div>
      </div>
    );
  }

  const config = pluginData?.config || {};
  const stats = pluginData?.stats || {};
  const endpoints = config.endpoints || [];

  return (
    <div>
      <div className="page-header">
        <h2>Webhooks</h2>
        <p>HTTP/HTTPS webhook receiver configuration and monitoring</p>
      </div>

      {/* Status Card */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">Webhook Module Status</div>
            <div className="card-subtitle">Current operational status</div>
          </div>
          <span className={`badge ${pluginData.status === 'running' ? 'badge-success' : 'badge-error'}`}>
            {pluginData.status}
          </span>
        </div>

        <div className="grid grid-4">
          <div>
            <div className="stat-label">Listen Address</div>
            <div style={{ fontSize: '16px', fontWeight: 600 }}>
              {config.listen_address || 'Not configured'}
            </div>
          </div>
          <div>
            <div className="stat-label">TLS Enabled</div>
            <div style={{ fontSize: '16px', fontWeight: 600 }}>
              {config.tls?.enabled ? (
                <span style={{ color: 'var(--success)' }}>✓ Yes</span>
              ) : (
                <span style={{ color: 'var(--text-tertiary)' }}>✗ No</span>
              )}
            </div>
          </div>
          <div>
            <div className="stat-label">Active Endpoints</div>
            <div style={{ fontSize: '16px', fontWeight: 600 }}>
              {endpoints.length}
            </div>
          </div>
          <div>
            <div className="stat-label">Events Received</div>
            <div style={{ fontSize: '16px', fontWeight: 600 }}>
              {stats.events_received?.toLocaleString() || 0}
            </div>
          </div>
        </div>
      </div>

      {/* Statistics */}
      <div className="grid grid-4">
        <div className="stat-card">
          <div className="stat-label">Total Requests</div>
          <div className="stat-value">{stats.events_received?.toLocaleString() || 0}</div>
          <div className="stat-change">
            <Activity size={14} />
            All time
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Successful</div>
          <div className="stat-value" style={{ color: 'var(--success)' }}>
            {stats.events_processed?.toLocaleString() || 0}
          </div>
          <div className="stat-change positive">
            Processed successfully
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Failed</div>
          <div className="stat-value" style={{
            color: (stats.errors || 0) > 0 ? 'var(--error)' : 'var(--text-primary)'
          }}>
            {stats.errors?.toLocaleString() || 0}
          </div>
          <div className="stat-change">
            Authentication or validation errors
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Dropped</div>
          <div className="stat-value">{stats.events_dropped?.toLocaleString() || 0}</div>
          <div className="stat-change">
            Rate limited or filtered
          </div>
        </div>
      </div>

      {/* Endpoints Configuration */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">Configured Endpoints</div>
            <div className="card-subtitle">{endpoints.length} webhook endpoint(s)</div>
          </div>
        </div>

        {endpoints.length > 0 ? (
          <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
            {endpoints.map((endpoint, index) => (
              <div key={index} className="card">
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start', marginBottom: '16px' }}>
                  <div>
                    <h3 style={{ fontSize: '16px', fontWeight: 600, marginBottom: '4px' }}>
                      {endpoint.name || `Endpoint ${index + 1}`}
                    </h3>
                    <code style={{ fontSize: '13px' }}>{endpoint.path}</code>
                  </div>
                  <span className="badge badge-success">Active</span>
                </div>

                <div className="grid grid-3" style={{ gap: '20px' }}>
                  <div>
                    <h4 style={{ fontSize: '14px', color: 'var(--text-secondary)', marginBottom: '8px', display: 'flex', alignItems: 'center', gap: '6px' }}>
                      <Shield size={16} />
                      Authentication
                    </h4>
                    <div style={{ fontSize: '14px' }}>
                      {endpoint.auth?.type ? (
                        <span className="badge badge-info">{endpoint.auth.type}</span>
                      ) : (
                        <span style={{ color: 'var(--text-tertiary)' }}>None</span>
                      )}
                    </div>
                  </div>

                  <div>
                    <h4 style={{ fontSize: '14px', color: 'var(--text-secondary)', marginBottom: '8px', display: 'flex', alignItems: 'center', gap: '6px' }}>
                      <MapPin size={16} />
                      IP Whitelist
                    </h4>
                    <div style={{ fontSize: '14px' }}>
                      {endpoint.ip_whitelist && endpoint.ip_whitelist.length > 0 ? (
                        <span className="badge badge-warning">
                          {endpoint.ip_whitelist.length} IP(s)
                        </span>
                      ) : (
                        <span style={{ color: 'var(--text-tertiary)' }}>All IPs allowed</span>
                      )}
                    </div>
                  </div>

                  <div>
                    <h4 style={{ fontSize: '14px', color: 'var(--text-secondary)', marginBottom: '8px', display: 'flex', alignItems: 'center', gap: '6px' }}>
                      <Activity size={16} />
                      Rate Limit
                    </h4>
                    <div style={{ fontSize: '14px' }}>
                      {endpoint.rate_limit ? (
                        <span className="badge badge-info">
                          {endpoint.rate_limit.requests}/{endpoint.rate_limit.window}
                        </span>
                      ) : (
                        <span style={{ color: 'var(--text-tertiary)' }}>Unlimited</span>
                      )}
                    </div>
                  </div>
                </div>

                {endpoint.field_mapping && Object.keys(endpoint.field_mapping).length > 0 && (
                  <div style={{ marginTop: '16px' }}>
                    <h4 style={{ fontSize: '14px', color: 'var(--text-secondary)', marginBottom: '8px' }}>
                      Field Mapping
                    </h4>
                    <div style={{
                      backgroundColor: 'var(--bg-tertiary)',
                      padding: '12px',
                      borderRadius: '6px',
                      fontSize: '12px',
                      fontFamily: 'monospace'
                    }}>
                      {Object.entries(endpoint.field_mapping).slice(0, 5).map(([key, value]) => (
                        <div key={key} style={{ marginBottom: '4px' }}>
                          <span style={{ color: 'var(--accent-primary)' }}>{key}</span>
                          {' → '}
                          <span>{value}</span>
                        </div>
                      ))}
                      {Object.keys(endpoint.field_mapping).length > 5 && (
                        <div style={{ color: 'var(--text-tertiary)', marginTop: '8px' }}>
                          ...and {Object.keys(endpoint.field_mapping).length - 5} more
                        </div>
                      )}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        ) : (
          <div className="empty-state">
            <div className="empty-state-icon">
              <Webhook size={48} />
            </div>
            <div>No webhook endpoints configured</div>
            <div style={{ marginTop: '12px', fontSize: '14px', color: 'var(--text-secondary)' }}>
              Add webhook endpoints in your configuration file
            </div>
          </div>
        )}
      </div>

      {/* Example Usage */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">Example Usage</div>
            <div className="card-subtitle">How to send events to this webhook receiver</div>
          </div>
        </div>

        <div style={{ fontSize: '14px', lineHeight: '1.8' }}>
          <p style={{ marginBottom: '16px', color: 'var(--text-secondary)' }}>
            Send HTTP POST requests to the configured endpoint(s):
          </p>

          <pre style={{
            backgroundColor: 'var(--bg-tertiary)',
            padding: '16px',
            borderRadius: '8px',
            overflow: 'auto',
            fontSize: '12px',
            lineHeight: '1.6',
          }}>
{`curl -X POST http://${config.listen_address || 'localhost:8080'}/webhook/prometheus \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer YOUR_TOKEN" \\
  -d '{
    "status": "critical",
    "title": "High CPU Usage",
    "host": "server-01",
    "service": "application"
  }'`}
          </pre>
        </div>
      </div>
    </div>
  );
}

export default Webhooks;
