import React, { useState, useEffect } from 'react';
import { AlertCircle, RefreshCw, Trash2, Eye } from 'lucide-react';
import api from '../api';

function DLQ() {
  const [dlqEvents, setDlqEvents] = useState([]);
  const [selectedEvent, setSelectedEvent] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadDLQEvents();
    const interval = setInterval(loadDLQEvents, 10000); // Refresh every 10 seconds
    return () => clearInterval(interval);
  }, []);

  const loadDLQEvents = async () => {
    try {
      const data = await api.getDLQEvents();
      setDlqEvents(data || []);
      setLoading(false);
      setError(null);
    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  };

  const viewEventDetails = (event) => {
    setSelectedEvent(event);
  };

  const closeDetails = () => {
    setSelectedEvent(null);
  };

  if (loading) {
    return (
      <div className="loading">
        <div className="spinner"></div>
      </div>
    );
  }

  return (
    <div>
      <div className="page-header">
        <h2>Dead Letter Queue</h2>
        <p>Failed events that require attention</p>
      </div>

      {error && (
        <div className="error-message">
          Failed to load DLQ events: {error}
        </div>
      )}

      {/* Summary */}
      <div className="card" style={{
        backgroundColor: dlqEvents.length > 0 ? 'rgba(239, 68, 68, 0.1)' : 'var(--bg-secondary)',
        borderColor: dlqEvents.length > 0 ? 'var(--error)' : 'var(--border)'
      }}>
        <div className="card-header">
          <div>
            <div className="card-title">
              {dlqEvents.length > 0 ? (
                <>
                  <AlertCircle size={20} style={{ marginRight: '8px', verticalAlign: 'middle' }} />
                  {dlqEvents.length} Failed Event{dlqEvents.length !== 1 ? 's' : ''}
                </>
              ) : (
                'No Failed Events'
              )}
            </div>
            <div className="card-subtitle">
              {dlqEvents.length > 0
                ? 'Events that failed to forward to BigPanda after multiple retries'
                : 'All events are being processed successfully'}
            </div>
          </div>
          {dlqEvents.length > 0 && (
            <button className="btn btn-secondary" onClick={loadDLQEvents}>
              <RefreshCw size={16} />
              Refresh
            </button>
          )}
        </div>
      </div>

      {/* DLQ Events Table */}
      {dlqEvents.length > 0 ? (
        <div className="card">
          <div className="card-header">
            <div>
              <div className="card-title">Failed Events</div>
              <div className="card-subtitle">
                Events that could not be delivered to BigPanda
              </div>
            </div>
          </div>

          <div className="table-container">
            <table>
              <thead>
                <tr>
                  <th>Timestamp</th>
                  <th>Source</th>
                  <th>Title</th>
                  <th>Host</th>
                  <th>Retry Count</th>
                  <th>Last Error</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {dlqEvents.map((event, index) => (
                  <tr key={index}>
                    <td style={{ fontSize: '12px' }}>
                      {event.timestamp
                        ? new Date(event.timestamp).toLocaleString()
                        : 'N/A'}
                    </td>
                    <td>
                      <code>{event.source || 'unknown'}</code>
                    </td>
                    <td style={{ maxWidth: '300px' }}>
                      {event.title || event.check || 'No title'}
                    </td>
                    <td>{event.host || '-'}</td>
                    <td>
                      <span className="badge badge-error">
                        {event.retry_count || 0} retries
                      </span>
                    </td>
                    <td style={{
                      maxWidth: '250px',
                      fontSize: '12px',
                      color: 'var(--error)'
                    }}>
                      {event.last_error || 'Unknown error'}
                    </td>
                    <td>
                      <button
                        className="btn btn-secondary"
                        style={{ padding: '6px 12px', fontSize: '12px' }}
                        onClick={() => viewEventDetails(event)}
                      >
                        <Eye size={14} />
                        View
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      ) : (
        <div className="card">
          <div className="empty-state">
            <div className="empty-state-icon">
              <AlertCircle size={48} style={{ opacity: 0.5, color: 'var(--success)' }} />
            </div>
            <div>No Failed Events</div>
            <div style={{ marginTop: '12px', fontSize: '14px', color: 'var(--text-secondary)' }}>
              All events are being processed and forwarded successfully
            </div>
          </div>
        </div>
      )}

      {/* Common Issues and Recommendations */}
      {dlqEvents.length > 0 && (
        <div className="card">
          <div className="card-header">
            <div>
              <div className="card-title">Troubleshooting</div>
              <div className="card-subtitle">Common causes and solutions</div>
            </div>
          </div>

          <div style={{ fontSize: '14px', lineHeight: '1.8' }}>
            <h4 style={{ fontSize: '16px', marginBottom: '12px' }}>Common Causes:</h4>
            <ul style={{ marginLeft: '20px', marginBottom: '20px' }}>
              <li style={{ marginBottom: '8px' }}>
                <strong>Authentication failure:</strong> Check your BigPanda API token and app key
              </li>
              <li style={{ marginBottom: '8px' }}>
                <strong>Network connectivity:</strong> Verify connection to BigPanda API endpoint
              </li>
              <li style={{ marginBottom: '8px' }}>
                <strong>Invalid event format:</strong> Event may be missing required fields
              </li>
              <li style={{ marginBottom: '8px' }}>
                <strong>Rate limiting:</strong> Too many requests to BigPanda API
              </li>
            </ul>

            <h4 style={{ fontSize: '16px', marginBottom: '12px' }}>Recommended Actions:</h4>
            <ul style={{ marginLeft: '20px' }}>
              <li style={{ marginBottom: '8px' }}>
                Review the agent configuration and BigPanda credentials
              </li>
              <li style={{ marginBottom: '8px' }}>
                Check network connectivity and firewall rules
              </li>
              <li style={{ marginBottom: '8px' }}>
                Examine the forwarder logs for detailed error messages
              </li>
              <li style={{ marginBottom: '8px' }}>
                Contact BigPanda support if authentication issues persist
              </li>
            </ul>
          </div>
        </div>
      )}

      {/* Event Details Modal */}
      {selectedEvent && (
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
        }}>
          <div style={{
            backgroundColor: 'var(--bg-secondary)',
            border: '1px solid var(--border)',
            borderRadius: '12px',
            padding: '24px',
            maxWidth: '800px',
            width: '100%',
            maxHeight: '80vh',
            overflow: 'auto',
          }}>
            <div style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              marginBottom: '20px',
            }}>
              <h3 style={{ fontSize: '20px', fontWeight: 600 }}>Event Details</h3>
              <button
                className="btn btn-secondary"
                onClick={closeDetails}
                style={{ padding: '8px 16px' }}
              >
                Close
              </button>
            </div>

            <div style={{ marginBottom: '20px' }}>
              <div style={{ marginBottom: '12px' }}>
                <strong>Error Information:</strong>
              </div>
              <div style={{
                backgroundColor: 'rgba(239, 68, 68, 0.1)',
                border: '1px solid var(--error)',
                borderRadius: '8px',
                padding: '12px',
                fontSize: '14px',
                color: 'var(--error)',
              }}>
                {selectedEvent.last_error || 'Unknown error'}
              </div>
            </div>

            <div style={{ marginBottom: '20px' }}>
              <div style={{ marginBottom: '12px' }}>
                <strong>Event Data:</strong>
              </div>
              <pre style={{
                backgroundColor: 'var(--bg-tertiary)',
                padding: '16px',
                borderRadius: '8px',
                overflow: 'auto',
                fontSize: '12px',
                lineHeight: '1.6',
                maxHeight: '400px',
              }}>
                {JSON.stringify(selectedEvent, null, 2)}
              </pre>
            </div>

            <div style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(2, 1fr)',
              gap: '16px',
              fontSize: '14px',
            }}>
              <div>
                <div style={{ color: 'var(--text-secondary)', marginBottom: '4px' }}>
                  Retry Count:
                </div>
                <div style={{ fontWeight: 600 }}>
                  {selectedEvent.retry_count || 0}
                </div>
              </div>
              <div>
                <div style={{ color: 'var(--text-secondary)', marginBottom: '4px' }}>
                  Timestamp:
                </div>
                <div style={{ fontWeight: 600 }}>
                  {selectedEvent.timestamp
                    ? new Date(selectedEvent.timestamp).toLocaleString()
                    : 'N/A'}
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default DLQ;
