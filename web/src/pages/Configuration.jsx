import React, { useState, useEffect } from 'react';
import { Settings, Save, AlertCircle, CheckCircle, RefreshCw } from 'lucide-react';
import api from '../api';
import BigPandaEndpoints from '../components/config/BigPandaEndpoints';
import SNMPConfig from '../components/config/SNMPConfig';
import WebhookConfig from '../components/config/WebhookConfig';
import GeneralConfig from '../components/config/GeneralConfig';

function Configuration() {
  const [activeTab, setActiveTab] = useState('bigpanda');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);

  const tabs = [
    { id: 'bigpanda', label: 'BigPanda' },
    { id: 'snmp', label: 'SNMP' },
    { id: 'webhook', label: 'Webhooks' },
    { id: 'general', label: 'General' },
  ];

  const clearMessages = () => {
    setError(null);
    setSuccess(null);
  };

  useEffect(() => {
    // Clear messages when switching tabs
    clearMessages();
  }, [activeTab]);

  const handleSave = async (configType, data) => {
    setLoading(true);
    clearMessages();

    try {
      let response;
      switch (configType) {
        case 'bigpanda':
          response = await api.updateBigPandaConfig(data);
          break;
        case 'snmp':
          response = await api.updateSNMPConfig(data);
          break;
        case 'webhook':
          response = await api.updateWebhookConfig(data);
          break;
        default:
          throw new Error('Unknown configuration type');
      }

      setSuccess(response.message || 'Configuration saved successfully');

      if (response.restart_required) {
        setSuccess(prev => `${prev}. Agent restart required for changes to take effect.`);
      }
    } catch (err) {
      setError(err.message || 'Failed to save configuration');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <div className="page-header">
        <div>
          <h2>Configuration</h2>
          <p>Manage agent and plugin configurations</p>
        </div>
      </div>

      {/* Global Messages */}
      {error && (
        <div className="card" style={{
          backgroundColor: 'rgba(239, 68, 68, 0.1)',
          borderColor: 'var(--error)',
          marginBottom: '20px'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <AlertCircle size={20} color="var(--error)" />
            <div style={{ flex: 1 }}>
              <div style={{ fontWeight: 600, color: 'var(--error)' }}>Error</div>
              <div style={{ fontSize: '14px', marginTop: '4px' }}>{error}</div>
            </div>
          </div>
        </div>
      )}

      {success && (
        <div className="card" style={{
          backgroundColor: 'rgba(34, 197, 94, 0.1)',
          borderColor: 'var(--success)',
          marginBottom: '20px'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <CheckCircle size={20} color="var(--success)" />
            <div style={{ flex: 1 }}>
              <div style={{ fontWeight: 600, color: 'var(--success)' }}>Success</div>
              <div style={{ fontSize: '14px', marginTop: '4px' }}>{success}</div>
            </div>
          </div>
        </div>
      )}

      {/* Tabs */}
      <div className="card" style={{ marginBottom: '20px' }}>
        <div style={{
          display: 'flex',
          gap: '4px',
          borderBottom: '1px solid var(--border)',
          padding: '0 20px'
        }}>
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              style={{
                padding: '12px 24px',
                background: 'none',
                border: 'none',
                borderBottom: activeTab === tab.id ? '2px solid var(--accent-primary)' : '2px solid transparent',
                color: activeTab === tab.id ? 'var(--accent-primary)' : 'var(--text-secondary)',
                fontWeight: activeTab === tab.id ? 600 : 400,
                cursor: 'pointer',
                fontSize: '14px',
                transition: 'all 0.2s',
              }}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {/* Tab Content */}
      <div style={{ position: 'relative' }}>
        {loading && (
          <div style={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundColor: 'rgba(0, 0, 0, 0.5)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            zIndex: 10,
            borderRadius: '12px',
          }}>
            <div style={{ textAlign: 'center' }}>
              <RefreshCw size={32} className="spinner" style={{ marginBottom: '12px' }} />
              <div style={{ color: 'white', fontWeight: 600 }}>Saving configuration...</div>
            </div>
          </div>
        )}

        {activeTab === 'bigpanda' && (
          <BigPandaEndpoints onSave={(data) => handleSave('bigpanda', data)} />
        )}

        {activeTab === 'snmp' && (
          <SNMPConfig onSave={(data) => handleSave('snmp', data)} />
        )}

        {activeTab === 'webhook' && (
          <WebhookConfig onSave={(data) => handleSave('webhook', data)} />
        )}

        {activeTab === 'general' && (
          <GeneralConfig />
        )}
      </div>
    </div>
  );
}

export default Configuration;
