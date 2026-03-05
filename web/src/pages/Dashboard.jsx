import React, { useState, useEffect } from 'react';
import { Activity, TrendingUp, AlertTriangle, CheckCircle } from 'lucide-react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import api from '../api';

function Dashboard() {
  const [stats, setStats] = useState(null);
  const [health, setHealth] = useState(null);
  const [agentInfo, setAgentInfo] = useState(null);
  const [queueStats, setQueueStats] = useState(null);
  const [plugins, setPlugins] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [eventHistory, setEventHistory] = useState([]);

  useEffect(() => {
    loadData();
    const interval = setInterval(loadData, 5000); // Refresh every 5 seconds
    return () => clearInterval(interval);
  }, []);

  const loadData = async () => {
    try {
      const [statsData, healthData, infoData, queueData] = await Promise.all([
        api.getStats(),
        api.getHealth(),
        api.getAgentInfo(),
        api.getQueueStats(),
      ]);

      setStats(statsData);
      setHealth(healthData);
      setAgentInfo(infoData);
      setQueueStats(queueData);

      // Extract plugin info from health data (has status) and stats data (has metrics)
      const pluginsList = [];
      if (healthData?.plugins?.plugins) {
        Object.entries(healthData.plugins.plugins).forEach(([name, pluginHealth]) => {
          const pluginStats = statsData?.plugins?.plugins?.[name] || {};
          pluginsList.push({
            name,
            status: pluginHealth.status === 'healthy' ? 'running' : 'stopped',
            health: pluginHealth,
            stats: pluginStats,
          });
        });
      }
      setPlugins(pluginsList);

      // Update event history for chart
      if (statsData) {
        setEventHistory(prev => {
          const newHistory = [...prev, {
            time: new Date().toLocaleTimeString(),
            events: statsData.plugins?.totals?.events_received || 0,
          }].slice(-20); // Keep last 20 data points
          return newHistory;
        });
      }

      setLoading(false);
      setError(null);
    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  };

  if (loading && !stats) {
    return (
      <div className="loading">
        <div className="spinner"></div>
      </div>
    );
  }

  if (error && !stats) {
    return (
      <div>
        <div className="page-header">
          <h2>Dashboard</h2>
          <p>Real-time monitoring and statistics</p>
        </div>
        <div className="error-message">
          Failed to load data: {error}
        </div>
      </div>
    );
  }

  const uptime = agentInfo?.uptime || 0;
  const uptimeHours = Math.floor(uptime / 3600);
  const uptimeMinutes = Math.floor((uptime % 3600) / 60);

  return (
    <div>
      <div className="page-header">
        <h2>Dashboard</h2>
        <p>Real-time monitoring and statistics</p>
      </div>

      {/* Health Status */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">System Health</div>
            <div className="card-subtitle">Current status of the agent</div>
          </div>
          {health?.status === 'healthy' ? (
            <span className="badge badge-success">
              <CheckCircle size={14} /> Healthy
            </span>
          ) : (
            <span className="badge badge-error">
              <AlertTriangle size={14} /> Degraded
            </span>
          )}
        </div>
        <div className="grid grid-4">
          <div>
            <div className="stat-label">Version</div>
            <div style={{ fontSize: '18px', fontWeight: 600 }}>
              {agentInfo?.version || 'N/A'}
            </div>
          </div>
          <div>
            <div className="stat-label">Uptime</div>
            <div style={{ fontSize: '18px', fontWeight: 600 }}>
              {uptimeHours}h {uptimeMinutes}m
            </div>
          </div>
          <div>
            <div className="stat-label">Active Plugins</div>
            <div style={{ fontSize: '18px', fontWeight: 600 }}>
              {plugins.filter(p => p.status === 'running').length} / {plugins.length}
            </div>
          </div>
          <div>
            <div className="stat-label">Last Updated</div>
            <div style={{ fontSize: '18px', fontWeight: 600 }}>
              {new Date().toLocaleTimeString()}
            </div>
          </div>
        </div>
      </div>

      {/* Key Metrics */}
      <div className="grid grid-4">
        <div className="stat-card">
          <div className="stat-label">Events Received</div>
          <div className="stat-value">{stats?.plugins?.totals?.events_received?.toLocaleString() || 0}</div>
          <div className="stat-change positive">
            <TrendingUp size={14} />
            Total processed
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Events Forwarded</div>
          <div className="stat-value">{stats?.forwarder?.sent?.toLocaleString() || 0}</div>
          <div className="stat-change positive">
            <TrendingUp size={14} />
            To BigPanda
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Queue Size</div>
          <div className="stat-value">{queueStats?.maxSize?.toLocaleString() || 0}</div>
          <div className="stat-change">
            <Activity size={14} />
            Max capacity
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Failed Events</div>
          <div className="stat-value">{(stats?.plugins?.totals?.errors || 0)?.toLocaleString()}</div>
          {(stats?.plugins?.totals?.errors || 0) > 0 ? (
            <div className="stat-change negative">
              <AlertTriangle size={14} />
              Check DLQ
            </div>
          ) : (
            <div className="stat-change">
              All good
            </div>
          )}
        </div>
      </div>

      {/* Event Rate Chart */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">Event Rate</div>
            <div className="card-subtitle">Events received over time</div>
          </div>
        </div>
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={eventHistory}>
            <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
            <XAxis dataKey="time" stroke="#94a3b8" />
            <YAxis stroke="#94a3b8" />
            <Tooltip
              contentStyle={{
                backgroundColor: '#1e293b',
                border: '1px solid #334155',
                borderRadius: '8px',
              }}
            />
            <Line
              type="monotone"
              dataKey="events"
              stroke="#3b82f6"
              strokeWidth={2}
              dot={false}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>

      {/* Active Plugins */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">Active Plugins</div>
            <div className="card-subtitle">Currently running receiver modules</div>
          </div>
        </div>
        <div className="grid grid-3">
          {plugins.map((plugin) => (
            <div key={plugin.name} className="card">
              <div style={{ marginBottom: '12px' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <h3 style={{ fontSize: '16px', fontWeight: 600 }}>{plugin.name}</h3>
                  <span className={`badge ${plugin.status === 'running' ? 'badge-success' : 'badge-error'}`}>
                    {plugin.status}
                  </span>
                </div>
              </div>
              <div style={{ fontSize: '14px', color: 'var(--text-secondary)' }}>
                <div>Events: {plugin.stats?.events_received?.toLocaleString() || 0}</div>
                <div>Errors: {plugin.stats?.errors?.toLocaleString() || 0}</div>
              </div>
            </div>
          ))}
        </div>
        {plugins.length === 0 && (
          <div className="empty-state">
            <div className="empty-state-icon">📦</div>
            <div>No plugins active</div>
          </div>
        )}
      </div>
    </div>
  );
}

export default Dashboard;
