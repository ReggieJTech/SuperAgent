import React, { useState, useEffect, useRef } from 'react';
import { Activity, PlayCircle, PauseCircle, Trash2 } from 'lucide-react';
import api from '../api';

function Events() {
  const [events, setEvents] = useState([]);
  const [recentEvents, setRecentEvents] = useState([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const wsRef = useRef(null);

  useEffect(() => {
    loadRecentEvents();
    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, []);

  const loadRecentEvents = async () => {
    try {
      const data = await api.getRecentEvents();
      // Backend returns {events: [], count: 0}
      const eventList = data?.events || data || [];
      setRecentEvents(Array.isArray(eventList) ? eventList : []);
      setLoading(false);
      setError(null);
    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  };

  const toggleStreaming = () => {
    if (isStreaming) {
      // Stop streaming
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
      setIsStreaming(false);
    } else {
      // Start streaming
      wsRef.current = api.connectEventStream(
        (event) => {
          setEvents((prev) => [event, ...prev].slice(0, 100)); // Keep last 100 events
        },
        (error) => {
          console.error('WebSocket error:', error);
          setIsStreaming(false);
        }
      );
      setIsStreaming(true);
    }
  };

  const clearEvents = () => {
    setEvents([]);
  };

  const getStatusBadge = (status) => {
    const badgeMap = {
      critical: 'badge-error',
      warning: 'badge-warning',
      ok: 'badge-success',
      info: 'badge-info',
    };
    return badgeMap[status?.toLowerCase()] || 'badge-info';
  };

  const formatTimestamp = (timestamp) => {
    if (!timestamp) return 'N/A';
    try {
      return new Date(timestamp).toLocaleString();
    } catch {
      return timestamp;
    }
  };

  if (loading) {
    return (
      <div className="loading">
        <div className="spinner"></div>
      </div>
    );
  }

  const displayEvents = isStreaming && events.length > 0 ? events : recentEvents;

  return (
    <div>
      <div className="page-header">
        <h2>Events</h2>
        <p>Monitor real-time event stream and recent activity</p>
      </div>

      {error && (
        <div className="error-message">
          Failed to load events: {error}
        </div>
      )}

      {/* Controls */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">Event Stream</div>
            <div className="card-subtitle">
              {isStreaming ? 'Live streaming active' : 'Showing recent events'}
            </div>
          </div>
          <div style={{ display: 'flex', gap: '12px' }}>
            {isStreaming && events.length > 0 && (
              <button className="btn btn-secondary" onClick={clearEvents}>
                <Trash2 size={16} />
                Clear
              </button>
            )}
            <button
              className={`btn ${isStreaming ? 'btn-secondary' : 'btn-primary'}`}
              onClick={toggleStreaming}
            >
              {isStreaming ? (
                <>
                  <PauseCircle size={16} />
                  Stop Stream
                </>
              ) : (
                <>
                  <PlayCircle size={16} />
                  Start Stream
                </>
              )}
            </button>
          </div>
        </div>

        {isStreaming && (
          <div style={{
            padding: '12px',
            backgroundColor: 'rgba(16, 185, 129, 0.1)',
            borderRadius: '8px',
            marginBottom: '16px',
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
          }}>
            <span className="status-indicator online"></span>
            <span style={{ fontSize: '14px', color: 'var(--success)' }}>
              Live stream active - {events.length} events received
            </span>
          </div>
        )}

        <div className="table-container">
          <table>
            <thead>
              <tr>
                <th>Timestamp</th>
                <th>Source</th>
                <th>Status</th>
                <th>Title</th>
                <th>Host</th>
                <th>Service</th>
                <th>Tags</th>
              </tr>
            </thead>
            <tbody>
              {displayEvents.map((event, index) => (
                <tr key={`${event.timestamp}-${index}`}>
                  <td style={{ fontSize: '12px' }}>
                    {formatTimestamp(event.timestamp)}
                  </td>
                  <td>
                    <code>{event.source || 'unknown'}</code>
                  </td>
                  <td>
                    <span className={`badge ${getStatusBadge(event.status)}`}>
                      {event.status || 'info'}
                    </span>
                  </td>
                  <td style={{ maxWidth: '300px' }}>
                    {event.title || event.check || 'No title'}
                  </td>
                  <td>{event.host || '-'}</td>
                  <td>{event.service || '-'}</td>
                  <td style={{ fontSize: '12px' }}>
                    {event.tags && Object.keys(event.tags).length > 0 ? (
                      <div style={{ display: 'flex', flexWrap: 'wrap', gap: '4px' }}>
                        {Object.entries(event.tags).slice(0, 3).map(([key, value]) => (
                          <span key={key} className="badge badge-info" style={{ fontSize: '10px' }}>
                            {key}:{value}
                          </span>
                        ))}
                        {Object.keys(event.tags).length > 3 && (
                          <span className="badge badge-info" style={{ fontSize: '10px' }}>
                            +{Object.keys(event.tags).length - 3}
                          </span>
                        )}
                      </div>
                    ) : (
                      '-'
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {displayEvents.length === 0 && (
          <div className="empty-state">
            <div className="empty-state-icon">
              <Activity size={48} />
            </div>
            <div>
              {isStreaming
                ? 'Waiting for events...'
                : 'No recent events. Start streaming to see live events.'}
            </div>
          </div>
        )}
      </div>

      {/* Event Details Modal would go here */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">Event Statistics</div>
            <div className="card-subtitle">Current session summary</div>
          </div>
        </div>
        <div className="grid grid-4">
          <div className="stat-card">
            <div className="stat-label">Total Events</div>
            <div className="stat-value">{displayEvents.length}</div>
          </div>
          <div className="stat-card">
            <div className="stat-label">Critical</div>
            <div className="stat-value" style={{ color: 'var(--error)' }}>
              {displayEvents.filter(e => e.status?.toLowerCase() === 'critical').length}
            </div>
          </div>
          <div className="stat-card">
            <div className="stat-label">Warning</div>
            <div className="stat-value" style={{ color: 'var(--warning)' }}>
              {displayEvents.filter(e => e.status?.toLowerCase() === 'warning').length}
            </div>
          </div>
          <div className="stat-card">
            <div className="stat-label">OK</div>
            <div className="stat-value" style={{ color: 'var(--success)' }}>
              {displayEvents.filter(e => e.status?.toLowerCase() === 'ok').length}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Events;
