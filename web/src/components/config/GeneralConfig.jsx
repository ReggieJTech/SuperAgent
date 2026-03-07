import React from 'react';
import { Info, Server, Shield, AlertCircle } from 'lucide-react';

function GeneralConfig() {
  return (
    <div>
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">General Settings</div>
            <div className="card-subtitle">System-wide agent configuration</div>
          </div>
        </div>

        <div style={{
          padding: '20px',
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          borderRadius: '8px',
          display: 'flex',
          alignItems: 'start',
          gap: '12px',
        }}>
          <Info size={20} style={{ marginTop: '2px', color: 'var(--accent-primary)' }} />
          <div style={{ flex: 1 }}>
            <div style={{ fontWeight: 600, marginBottom: '8px' }}>Configuration Information</div>
            <div style={{ fontSize: '14px', lineHeight: '1.6', color: 'var(--text-secondary)' }}>
              General agent settings such as logging, monitoring, authentication, and security
              are configured in the main agent configuration file. These settings require
              an agent restart to take effect and are not currently editable through the UI.
            </div>
          </div>
        </div>
      </div>

      <div className="card">
        <div className="card-header">
          <div className="card-title">Configuration Files</div>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
          <div style={{
            padding: '16px',
            backgroundColor: 'var(--bg-tertiary)',
            borderRadius: '8px',
          }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '12px' }}>
              <Server size={20} />
              <div style={{ fontWeight: 600 }}>Main Configuration</div>
            </div>
            <code style={{ fontSize: '13px' }}>/etc/bigpanda-agent/config.yaml</code>
            <div style={{ fontSize: '14px', marginTop: '8px', color: 'var(--text-secondary)' }}>
              Contains server, logging, monitoring, authentication, and security settings
            </div>
          </div>

          <div style={{
            padding: '16px',
            backgroundColor: 'var(--bg-tertiary)',
            borderRadius: '8px',
          }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '12px' }}>
              <Server size={20} />
              <div style={{ fontWeight: 600 }}>Module Configurations</div>
            </div>
            <div style={{ fontSize: '14px', lineHeight: '1.8' }}>
              <div><code>/etc/bigpanda-agent/modules/snmp.yaml</code></div>
              <div><code>/etc/bigpanda-agent/modules/webhook.yaml</code></div>
              <div><code>/etc/bigpanda-agent/modules/automation.yaml</code></div>
            </div>
          </div>
        </div>
      </div>

      <div className="card">
        <div className="card-header">
          <div className="card-title">Available Settings</div>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          <div>
            <h4 style={{ fontSize: '16px', fontWeight: 600, marginBottom: '12px', display: 'flex', alignItems: 'center', gap: '8px' }}>
              <Server size={18} />
              Server Settings
            </h4>
            <ul style={{ marginLeft: '20px', lineHeight: '1.8', fontSize: '14px', color: 'var(--text-secondary)' }}>
              <li>Listen address and port</li>
              <li>TLS/SSL configuration</li>
              <li>Timeouts (read, write, idle)</li>
            </ul>
          </div>

          <div>
            <h4 style={{ fontSize: '16px', fontWeight: 600, marginBottom: '12px', display: 'flex', alignItems: 'center', gap: '8px' }}>
              <Shield size={18} />
              Authentication & Security
            </h4>
            <ul style={{ marginLeft: '20px', lineHeight: '1.8', fontSize: '14px', color: 'var(--text-secondary)' }}>
              <li>Local user authentication</li>
              <li>LDAP integration</li>
              <li>SSO/SAML support</li>
              <li>Credential encryption</li>
              <li>API rate limiting</li>
            </ul>
          </div>

          <div>
            <h4 style={{ fontSize: '16px', fontWeight: 600, marginBottom: '12px', display: 'flex', alignItems: 'center', gap: '8px' }}>
              <AlertCircle size={18} />
              Monitoring & Logging
            </h4>
            <ul style={{ marginLeft: '20px', lineHeight: '1.8', fontSize: '14px', color: 'var(--text-secondary)' }}>
              <li>Log level and format</li>
              <li>Log rotation settings</li>
              <li>Metrics endpoint configuration</li>
              <li>Health check settings</li>
              <li>Heartbeat interval</li>
            </ul>
          </div>
        </div>
      </div>

      <div className="card">
        <div className="card-header">
          <div className="card-title">Documentation</div>
        </div>

        <div style={{ fontSize: '14px', lineHeight: '1.8' }}>
          <p style={{ marginBottom: '12px' }}>
            For detailed configuration options and examples, refer to the documentation:
          </p>
          <ul style={{ marginLeft: '20px' }}>
            <li style={{ marginBottom: '8px' }}>
              <a href="/docs/deployment-guide.md" style={{ color: 'var(--accent-primary)' }}>
                Deployment Guide
              </a> - Full installation and configuration guide
            </li>
            <li style={{ marginBottom: '8px' }}>
              <a href="/docs/snmp-guide.md" style={{ color: 'var(--accent-primary)' }}>
                SNMP Module Guide
              </a> - SNMP trap configuration
            </li>
            <li style={{ marginBottom: '8px' }}>
              <a href="/docs/api-reference.md" style={{ color: 'var(--accent-primary)' }}>
                API Reference
              </a> - REST API documentation
            </li>
          </ul>
        </div>
      </div>
    </div>
  );
}

export default GeneralConfig;
