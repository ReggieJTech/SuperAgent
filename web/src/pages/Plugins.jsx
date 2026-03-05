import React, { useState, useEffect } from 'react';
import { Layers, Activity, AlertCircle, CheckCircle, XCircle } from 'lucide-react';
import api from '../api';

function Plugins() {
  const [plugins, setPlugins] = useState([]);
  const [selectedPlugin, setSelectedPlugin] = useState(null);
  const [pluginDetails, setPluginDetails] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadPlugins();
    const interval = setInterval(loadPlugins, 5000);
    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    if (selectedPlugin) {
      loadPluginDetails(selectedPlugin);
    }
  }, [selectedPlugin]);

  const loadPlugins = async () => {
    try {
      const data = await api.getPlugins();
      setPlugins(data);
      setLoading(false);
      setError(null);
    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  };

  const loadPluginDetails = async (name) => {
    try {
      const details = await api.getPlugin(name);
      setPluginDetails(details);
    } catch (err) {
      console.error('Failed to load plugin details:', err);
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
          <h2>Plugins</h2>
          <p>Manage receiver modules</p>
        </div>
        <div className="error-message">
          Failed to load plugins: {error}
        </div>
      </div>
    );
  }

  const getStatusIcon = (status) => {
    switch (status) {
      case 'running':
        return <CheckCircle size={20} color="#10b981" />;
      case 'stopped':
        return <XCircle size={20} color="#ef4444" />;
      default:
        return <AlertCircle size={20} color="#f59e0b" />;
    }
  };

  const getStatusBadge = (status) => {
    const badgeClass = {
      running: 'badge-success',
      stopped: 'badge-error',
      starting: 'badge-warning',
    }[status] || 'badge-info';

    return <span className={`badge ${badgeClass}`}>{status}</span>;
  };

  return (
    <div>
      <div className="page-header">
        <h2>Plugins</h2>
        <p>Manage and monitor receiver modules</p>
      </div>

      {/* Plugin List */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">Active Plugins</div>
            <div className="card-subtitle">{plugins.length} plugin(s) configured</div>
          </div>
        </div>

        <div className="table-container">
          <table>
            <thead>
              <tr>
                <th>Status</th>
                <th>Name</th>
                <th>Type</th>
                <th>Events Received</th>
                <th>Events Processed</th>
                <th>Errors</th>
                <th>Last Activity</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {plugins.map((plugin) => (
                <tr key={plugin.name}>
                  <td>{getStatusIcon(plugin.status)}</td>
                  <td style={{ fontWeight: 600 }}>{plugin.name}</td>
                  <td>{getStatusBadge(plugin.status)}</td>
                  <td>{plugin.stats?.events_received?.toLocaleString() || 0}</td>
                  <td>{plugin.stats?.events_processed?.toLocaleString() || 0}</td>
                  <td>
                    {(plugin.stats?.errors || 0) > 0 ? (
                      <span style={{ color: 'var(--error)' }}>
                        {plugin.stats.errors.toLocaleString()}
                      </span>
                    ) : (
                      '0'
                    )}
                  </td>
                  <td>
                    {plugin.stats?.last_event_time
                      ? new Date(plugin.stats.last_event_time).toLocaleString()
                      : 'Never'}
                  </td>
                  <td>
                    <button
                      className="btn btn-secondary"
                      style={{ padding: '6px 12px', fontSize: '12px' }}
                      onClick={() => setSelectedPlugin(plugin.name)}
                    >
                      Details
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {plugins.length === 0 && (
          <div className="empty-state">
            <div className="empty-state-icon">
              <Layers size={48} />
            </div>
            <div>No plugins configured</div>
          </div>
        )}
      </div>

      {/* Plugin Details */}
      {selectedPlugin && pluginDetails && (
        <div className="card">
          <div className="card-header">
            <div>
              <div className="card-title">Plugin Details: {selectedPlugin}</div>
              <div className="card-subtitle">Detailed information and configuration</div>
            </div>
            <button
              className="btn btn-secondary"
              onClick={() => {
                setSelectedPlugin(null);
                setPluginDetails(null);
              }}
            >
              Close
            </button>
          </div>

          <div className="grid grid-2">
            <div>
              <h3 style={{ fontSize: '16px', marginBottom: '12px' }}>Information</h3>
              <div style={{ fontSize: '14px', lineHeight: '1.8' }}>
                <div><strong>Name:</strong> {pluginDetails.name}</div>
                <div><strong>Status:</strong> {getStatusBadge(pluginDetails.status)}</div>
                <div><strong>Type:</strong> {pluginDetails.type || 'Receiver'}</div>
                <div>
                  <strong>Health:</strong>{' '}
                  {pluginDetails.health?.status === 'healthy' ? (
                    <span style={{ color: 'var(--success)' }}>Healthy</span>
                  ) : (
                    <span style={{ color: 'var(--error)' }}>Unhealthy</span>
                  )}
                </div>
              </div>
            </div>

            <div>
              <h3 style={{ fontSize: '16px', marginBottom: '12px' }}>Statistics</h3>
              <div style={{ fontSize: '14px', lineHeight: '1.8' }}>
                <div><strong>Events Received:</strong> {pluginDetails.stats?.events_received?.toLocaleString() || 0}</div>
                <div><strong>Events Processed:</strong> {pluginDetails.stats?.events_processed?.toLocaleString() || 0}</div>
                <div><strong>Events Dropped:</strong> {pluginDetails.stats?.events_dropped?.toLocaleString() || 0}</div>
                <div><strong>Errors:</strong> {pluginDetails.stats?.errors?.toLocaleString() || 0}</div>
              </div>
            </div>
          </div>

          {pluginDetails.config && (
            <div style={{ marginTop: '24px' }}>
              <h3 style={{ fontSize: '16px', marginBottom: '12px' }}>Configuration</h3>
              <pre style={{
                backgroundColor: 'var(--bg-tertiary)',
                padding: '16px',
                borderRadius: '8px',
                overflow: 'auto',
                fontSize: '12px',
                lineHeight: '1.6',
              }}>
                {JSON.stringify(pluginDetails.config, null, 2)}
              </pre>
            </div>
          )}

          {pluginDetails.health?.message && (
            <div style={{ marginTop: '24px' }}>
              <h3 style={{ fontSize: '16px', marginBottom: '12px' }}>Health Message</h3>
              <div className={pluginDetails.health.status === 'healthy' ? 'badge badge-success' : 'badge badge-error'}>
                {pluginDetails.health.message}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default Plugins;
