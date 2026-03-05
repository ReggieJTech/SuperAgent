import React, { useState, useEffect } from 'react';
import { Database, TrendingUp, AlertTriangle } from 'lucide-react';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import api from '../api';

function Queue() {
  const [queueStats, setQueueStats] = useState(null);
  const [queueSize, setQueueSize] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [history, setHistory] = useState([]);

  useEffect(() => {
    loadData();
    const interval = setInterval(loadData, 3000); // Refresh every 3 seconds
    return () => clearInterval(interval);
  }, []);

  const loadData = async () => {
    try {
      const [stats, size] = await Promise.all([
        api.getQueueStats(),
        api.getQueueSize(),
      ]);

      setQueueStats(stats);
      setQueueSize(size.size || 0);

      // Update history for chart
      setHistory(prev => {
        const newHistory = [...prev, {
          time: new Date().toLocaleTimeString(),
          size: size.size || 0,
          enqueued: stats.enqueued || 0,
          dequeued: stats.dequeued || 0,
        }].slice(-30); // Keep last 30 data points
        return newHistory;
      });

      setLoading(false);
      setError(null);
    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  };

  if (loading && !queueStats) {
    return (
      <div className="loading">
        <div className="spinner"></div>
      </div>
    );
  }

  if (error && !queueStats) {
    return (
      <div>
        <div className="page-header">
          <h2>Queue</h2>
          <p>Event queue monitoring and statistics</p>
        </div>
        <div className="error-message">
          Failed to load queue data: {error}
        </div>
      </div>
    );
  }

  const throughputRate = queueStats ?
    (queueStats.dequeued / Math.max(1, queueStats.enqueued) * 100).toFixed(1) : 0;

  const queueHealth = queueSize < 100 ? 'healthy' : queueSize < 1000 ? 'warning' : 'critical';

  return (
    <div>
      <div className="page-header">
        <h2>Queue</h2>
        <p>Event queue monitoring and statistics</p>
      </div>

      {/* Health Status */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">Queue Health</div>
            <div className="card-subtitle">Current queue status</div>
          </div>
          <span className={`badge ${
            queueHealth === 'healthy' ? 'badge-success' :
            queueHealth === 'warning' ? 'badge-warning' :
            'badge-error'
          }`}>
            {queueHealth === 'healthy' && '✓ Healthy'}
            {queueHealth === 'warning' && '⚠ Warning'}
            {queueHealth === 'critical' && '✗ Critical'}
          </span>
        </div>
      </div>

      {/* Key Metrics */}
      <div className="grid grid-4">
        <div className="stat-card">
          <div className="stat-label">Current Queue Size</div>
          <div className="stat-value" style={{
            color: queueHealth === 'critical' ? 'var(--error)' :
                   queueHealth === 'warning' ? 'var(--warning)' :
                   'var(--text-primary)'
          }}>
            {queueSize.toLocaleString()}
          </div>
          <div className="stat-change">
            <Database size={14} />
            Events pending
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-label">Events Enqueued</div>
          <div className="stat-value">{queueStats?.enqueued?.toLocaleString() || 0}</div>
          <div className="stat-change positive">
            <TrendingUp size={14} />
            Total added
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-label">Events Dequeued</div>
          <div className="stat-value">{queueStats?.dequeued?.toLocaleString() || 0}</div>
          <div className="stat-change positive">
            <TrendingUp size={14} />
            Total processed
          </div>
        </div>

        <div className="stat-card">
          <div className="stat-label">Throughput Rate</div>
          <div className="stat-value">{throughputRate}%</div>
          <div className="stat-change">
            Processing efficiency
          </div>
        </div>
      </div>

      {/* Queue Size Chart */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">Queue Size Over Time</div>
            <div className="card-subtitle">Real-time queue depth monitoring</div>
          </div>
        </div>
        <ResponsiveContainer width="100%" height={300}>
          <AreaChart data={history}>
            <defs>
              <linearGradient id="colorSize" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.8}/>
                <stop offset="95%" stopColor="#3b82f6" stopOpacity={0}/>
              </linearGradient>
            </defs>
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
            <Area
              type="monotone"
              dataKey="size"
              stroke="#3b82f6"
              fillOpacity={1}
              fill="url(#colorSize)"
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>

      {/* Additional Stats */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">Detailed Statistics</div>
            <div className="card-subtitle">Queue performance metrics</div>
          </div>
        </div>

        <div className="grid grid-2">
          <div>
            <h3 style={{ fontSize: '16px', marginBottom: '16px' }}>Queue Operations</h3>
            <div style={{ fontSize: '14px', lineHeight: '2' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: 'var(--text-secondary)' }}>Total Enqueued:</span>
                <strong>{queueStats?.enqueued?.toLocaleString() || 0}</strong>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: 'var(--text-secondary)' }}>Total Dequeued:</span>
                <strong>{queueStats?.dequeued?.toLocaleString() || 0}</strong>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: 'var(--text-secondary)' }}>Current Size:</span>
                <strong>{queueSize.toLocaleString()}</strong>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: 'var(--text-secondary)' }}>Failed Events:</span>
                <strong style={{ color: queueStats?.failed > 0 ? 'var(--error)' : 'inherit' }}>
                  {queueStats?.failed?.toLocaleString() || 0}
                </strong>
              </div>
            </div>
          </div>

          <div>
            <h3 style={{ fontSize: '16px', marginBottom: '16px' }}>Performance</h3>
            <div style={{ fontSize: '14px', lineHeight: '2' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: 'var(--text-secondary)' }}>Processing Rate:</span>
                <strong>{throughputRate}%</strong>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: 'var(--text-secondary)' }}>Backend:</span>
                <strong>{queueStats?.backend || 'BadgerDB'}</strong>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: 'var(--text-secondary)' }}>Persistence:</span>
                <strong>{queueStats?.persistent ? 'Enabled' : 'Disabled'}</strong>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: 'var(--text-secondary)' }}>Status:</span>
                <span className={`badge ${queueHealth === 'healthy' ? 'badge-success' :
                                         queueHealth === 'warning' ? 'badge-warning' :
                                         'badge-error'}`}>
                  {queueHealth}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {queueSize > 1000 && (
        <div className="card" style={{
          backgroundColor: 'rgba(239, 68, 68, 0.1)',
          borderColor: 'var(--error)'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <AlertTriangle size={24} color="var(--error)" />
            <div>
              <div style={{ fontWeight: 600, marginBottom: '4px' }}>High Queue Backlog Detected</div>
              <div style={{ fontSize: '14px', color: 'var(--text-secondary)' }}>
                The queue has {queueSize.toLocaleString()} pending events. Consider checking the forwarder
                status and BigPanda API connectivity.
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default Queue;
