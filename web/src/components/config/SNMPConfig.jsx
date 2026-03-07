import React, { useState, useEffect } from 'react';
import { Save, Upload, FileText, Plus, Trash2, RefreshCw } from 'lucide-react';
import api from '../../api';

function SNMPConfig({ onSave }) {
  const [config, setConfig] = useState({
    listen_address: '0.0.0.0:162',
    snmp_version: '2c',
    community: 'public',
    v3: {
      security_level: 'authPriv',
      auth_protocol: 'SHA',
      auth_password: '',
      priv_protocol: 'AES',
      priv_password: '',
      security_name: 'bigpanda',
    },
    filtering: {
      enabled: false,
      rules: [],
    },
    rate_limiting: {
      enabled: true,
      per_source: 100,
      global: 1000,
      burst: 200,
    },
    performance: {
      workers: 4,
      buffer_size: 1000,
      batch_size: 50,
    },
    routing: {
      default_endpoints: ['default'],
      rules: [],
    },
  });

  const [availableEndpoints, setAvailableEndpoints] = useState([]);

  const [loading, setLoading] = useState(true);
  const [modified, setModified] = useState(false);
  const [uploadingMIB, setUploadingMIB] = useState(false);
  const [uploadSuccess, setUploadSuccess] = useState(null);

  useEffect(() => {
    loadConfig();
    loadEndpoints();
  }, []);

  const loadConfig = async () => {
    try {
      const data = await api.getSNMPConfig();
      if (data && data.config) {
        setConfig(data.config);
      }
      setLoading(false);
    } catch (err) {
      console.error('Failed to load SNMP config:', err);
      setLoading(false);
    }
  };

  const loadEndpoints = async () => {
    try {
      const data = await api.getBigPandaEndpoints();
      const endpoints = data?.endpoints || data || [];
      setAvailableEndpoints(endpoints.filter(e => e.enabled).map(e => e.name));
    } catch (err) {
      console.error('Failed to load BigPanda endpoints:', err);
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

  const handleMIBUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    setUploadingMIB(true);
    setUploadSuccess(null);

    try {
      const result = await api.uploadMIB(file);
      setUploadSuccess(`MIB file "${result.filename}" uploaded successfully`);
      setTimeout(() => setUploadSuccess(null), 5000);
    } catch (err) {
      alert(`Failed to upload MIB: ${err.message}`);
    } finally {
      setUploadingMIB(false);
      e.target.value = '';
    }
  };

  const addFilterRule = () => {
    const rules = [...(config.filtering?.rules || [])];
    rules.push({
      action: 'drop',
      type: 'oid',
      pattern: '',
    });
    handleChange('filtering.rules', rules);
  };

  const removeFilterRule = (index) => {
    const rules = [...(config.filtering?.rules || [])];
    rules.splice(index, 1);
    handleChange('filtering.rules', rules);
  };

  const updateFilterRule = (index, field, value) => {
    const rules = [...(config.filtering?.rules || [])];
    rules[index] = { ...rules[index], [field]: value };
    handleChange('filtering.rules', rules);
  };

  // Routing rule management
  const addRoutingRule = () => {
    const rules = [...(config.routing?.rules || [])];
    rules.push({
      name: '',
      match_type: 'vendor',
      match_value: '',
      endpoints: availableEndpoints.length > 0 ? [availableEndpoints[0]] : ['default'],
      priority: 100,
    });
    handleChange('routing.rules', rules);
  };

  const removeRoutingRule = (index) => {
    const rules = [...(config.routing?.rules || [])];
    rules.splice(index, 1);
    handleChange('routing.rules', rules);
  };

  const updateRoutingRule = (index, field, value) => {
    const rules = [...(config.routing?.rules || [])];
    rules[index] = { ...rules[index], [field]: value };
    handleChange('routing.rules', rules);
  };

  const toggleRoutingEndpoint = (ruleIndex, endpoint) => {
    const rules = [...(config.routing?.rules || [])];
    const currentEndpoints = rules[ruleIndex].endpoints || [];

    if (currentEndpoints.includes(endpoint)) {
      rules[ruleIndex].endpoints = currentEndpoints.filter(e => e !== endpoint);
    } else {
      rules[ruleIndex].endpoints = [...currentEndpoints, endpoint];
    }

    handleChange('routing.rules', rules);
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
      {/* MIB Upload Section */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">MIB Management</div>
            <div className="card-subtitle">Upload and manage SNMP MIB files</div>
          </div>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
          {uploadSuccess && (
            <div style={{
              padding: '12px',
              backgroundColor: 'rgba(34, 197, 94, 0.1)',
              border: '1px solid var(--success)',
              borderRadius: '8px',
              color: 'var(--success)',
              fontSize: '14px',
            }}>
              {uploadSuccess}
            </div>
          )}

          <div>
            <label className="btn btn-secondary" style={{ cursor: 'pointer' }}>
              <Upload size={18} />
              {uploadingMIB ? 'Uploading...' : 'Upload MIB File'}
              <input
                type="file"
                accept=".mib,.txt"
                onChange={handleMIBUpload}
                disabled={uploadingMIB}
                style={{ display: 'none' }}
              />
            </label>
            <div className="form-help" style={{ marginTop: '8px' }}>
              Upload SNMP MIB files (.mib or .txt format)
            </div>
          </div>

          <div style={{
            padding: '16px',
            backgroundColor: 'var(--bg-tertiary)',
            borderRadius: '8px',
            fontSize: '14px',
          }}>
            <div style={{ fontWeight: 600, marginBottom: '8px' }}>MIB Upload Tips:</div>
            <ul style={{ marginLeft: '20px', lineHeight: '1.8' }}>
              <li>Upload vendor-specific MIB files to enable trap processing</li>
              <li>MIBs are automatically compiled and made available to the agent</li>
              <li>Use the event config generator to create mappings from MIB traps</li>
            </ul>
          </div>
        </div>
      </div>

      {/* Listener Configuration */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">SNMP Listener Configuration</div>
            <div className="card-subtitle">Configure SNMP trap receiver settings</div>
          </div>
          {modified && (
            <span className="badge badge-warning">Unsaved Changes</span>
          )}
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          <div className="grid grid-2">
            <div>
              <label className="form-label">Listen Address</label>
              <input
                type="text"
                className="form-input"
                value={config.listen_address || '0.0.0.0:162'}
                onChange={(e) => handleChange('listen_address', e.target.value)}
                placeholder="0.0.0.0:162"
                required
              />
              <div className="form-help">IP and port to listen for SNMP traps</div>
            </div>
            <div>
              <label className="form-label">SNMP Version</label>
              <select
                className="form-input"
                value={config.snmp_version || '2c'}
                onChange={(e) => handleChange('snmp_version', e.target.value)}
              >
                <option value="1">v1</option>
                <option value="2c">v2c</option>
                <option value="3">v3</option>
              </select>
            </div>
          </div>

          {(config.snmp_version === '1' || config.snmp_version === '2c') && (
            <div>
              <label className="form-label">Community String</label>
              <input
                type="text"
                className="form-input"
                value={config.community || 'public'}
                onChange={(e) => handleChange('community', e.target.value)}
                placeholder="public"
              />
            </div>
          )}

          {config.snmp_version === '3' && (
            <>
              <div className="grid grid-2">
                <div>
                  <label className="form-label">Security Level</label>
                  <select
                    className="form-input"
                    value={config.v3?.security_level || 'authPriv'}
                    onChange={(e) => handleChange('v3.security_level', e.target.value)}
                  >
                    <option value="noAuthNoPriv">No Auth, No Priv</option>
                    <option value="authNoPriv">Auth, No Priv</option>
                    <option value="authPriv">Auth and Priv</option>
                  </select>
                </div>
                <div>
                  <label className="form-label">Security Name</label>
                  <input
                    type="text"
                    className="form-input"
                    value={config.v3?.security_name || ''}
                    onChange={(e) => handleChange('v3.security_name', e.target.value)}
                    placeholder="bigpanda"
                  />
                </div>
              </div>

              {(config.v3?.security_level === 'authNoPriv' || config.v3?.security_level === 'authPriv') && (
                <div className="grid grid-2">
                  <div>
                    <label className="form-label">Auth Protocol</label>
                    <select
                      className="form-input"
                      value={config.v3?.auth_protocol || 'SHA'}
                      onChange={(e) => handleChange('v3.auth_protocol', e.target.value)}
                    >
                      <option value="MD5">MD5</option>
                      <option value="SHA">SHA</option>
                      <option value="SHA256">SHA256</option>
                    </select>
                  </div>
                  <div>
                    <label className="form-label">Auth Password</label>
                    <input
                      type="password"
                      className="form-input"
                      value={config.v3?.auth_password || ''}
                      onChange={(e) => handleChange('v3.auth_password', e.target.value)}
                      placeholder="Authentication password"
                    />
                  </div>
                </div>
              )}

              {config.v3?.security_level === 'authPriv' && (
                <div className="grid grid-2">
                  <div>
                    <label className="form-label">Priv Protocol</label>
                    <select
                      className="form-input"
                      value={config.v3?.priv_protocol || 'AES'}
                      onChange={(e) => handleChange('v3.priv_protocol', e.target.value)}
                    >
                      <option value="DES">DES</option>
                      <option value="AES">AES</option>
                      <option value="AES256">AES256</option>
                    </select>
                  </div>
                  <div>
                    <label className="form-label">Priv Password</label>
                    <input
                      type="password"
                      className="form-input"
                      value={config.v3?.priv_password || ''}
                      onChange={(e) => handleChange('v3.priv_password', e.target.value)}
                      placeholder="Privacy password"
                    />
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      </div>

      {/* Filtering Rules */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">Filtering Rules</div>
            <div className="card-subtitle">Filter incoming SNMP traps</div>
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <input
              type="checkbox"
              id="filtering-enabled"
              checked={config.filtering?.enabled || false}
              onChange={(e) => handleChange('filtering.enabled', e.target.checked)}
            />
            <label htmlFor="filtering-enabled" style={{ cursor: 'pointer', fontSize: '14px' }}>
              Enable Filtering
            </label>
          </div>
        </div>

        {config.filtering?.enabled && (
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            {config.filtering?.rules?.map((rule, index) => (
              <div key={index} style={{
                padding: '16px',
                backgroundColor: 'var(--bg-tertiary)',
                borderRadius: '8px',
                display: 'flex',
                gap: '12px',
                alignItems: 'end',
              }}>
                <div style={{ flex: 1 }}>
                  <label className="form-label">Action</label>
                  <select
                    className="form-input"
                    value={rule.action}
                    onChange={(e) => updateFilterRule(index, 'action', e.target.value)}
                  >
                    <option value="drop">Drop</option>
                    <option value="accept">Accept</option>
                  </select>
                </div>
                <div style={{ flex: 1 }}>
                  <label className="form-label">Type</label>
                  <select
                    className="form-input"
                    value={rule.type}
                    onChange={(e) => updateFilterRule(index, 'type', e.target.value)}
                  >
                    <option value="oid">OID</option>
                    <option value="source">Source IP</option>
                    <option value="source_network">Source Network</option>
                  </select>
                </div>
                <div style={{ flex: 2 }}>
                  <label className="form-label">Pattern</label>
                  <input
                    type="text"
                    className="form-input"
                    value={rule.pattern}
                    onChange={(e) => updateFilterRule(index, 'pattern', e.target.value)}
                    placeholder={
                      rule.type === 'oid' ? '1.3.6.1.4.1.*' :
                      rule.type === 'source' ? '10.0.0.100' :
                      '10.0.0.0/8'
                    }
                  />
                </div>
                <button
                  type="button"
                  onClick={() => removeFilterRule(index)}
                  className="btn btn-secondary"
                  style={{ padding: '8px 12px' }}
                >
                  <Trash2 size={18} />
                </button>
              </div>
            ))}

            <button
              type="button"
              onClick={addFilterRule}
              className="btn btn-secondary"
              style={{ alignSelf: 'flex-start' }}
            >
              <Plus size={18} />
              Add Filter Rule
            </button>
          </div>
        )}
      </div>

      {/* Rate Limiting */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">Rate Limiting</div>
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
            <input
              type="checkbox"
              id="rate-limiting-enabled"
              checked={config.rate_limiting?.enabled || false}
              onChange={(e) => handleChange('rate_limiting.enabled', e.target.checked)}
            />
            <label htmlFor="rate-limiting-enabled" style={{ cursor: 'pointer', fontSize: '14px' }}>
              Enable Rate Limiting
            </label>
          </div>
        </div>

        {config.rate_limiting?.enabled && (
          <div className="grid grid-3">
            <div>
              <label className="form-label">Per Source (traps/sec)</label>
              <input
                type="number"
                className="form-input"
                value={config.rate_limiting?.per_source || 100}
                onChange={(e) => handleChange('rate_limiting.per_source', parseInt(e.target.value))}
                min="1"
              />
            </div>
            <div>
              <label className="form-label">Global (traps/sec)</label>
              <input
                type="number"
                className="form-input"
                value={config.rate_limiting?.global || 1000}
                onChange={(e) => handleChange('rate_limiting.global', parseInt(e.target.value))}
                min="1"
              />
            </div>
            <div>
              <label className="form-label">Burst</label>
              <input
                type="number"
                className="form-input"
                value={config.rate_limiting?.burst || 200}
                onChange={(e) => handleChange('rate_limiting.burst', parseInt(e.target.value))}
                min="1"
              />
            </div>
          </div>
        )}
      </div>

      {/* Performance Tuning */}
      <div className="card">
        <div className="card-header">
          <div className="card-title">Performance Tuning</div>
        </div>

        <div className="grid grid-3">
          <div>
            <label className="form-label">Workers</label>
            <input
              type="number"
              className="form-input"
              value={config.performance?.workers || 4}
              onChange={(e) => handleChange('performance.workers', parseInt(e.target.value))}
              min="1"
              max="32"
            />
            <div className="form-help">Processing workers</div>
          </div>
          <div>
            <label className="form-label">Buffer Size</label>
            <input
              type="number"
              className="form-input"
              value={config.performance?.buffer_size || 1000}
              onChange={(e) => handleChange('performance.buffer_size', parseInt(e.target.value))}
              min="100"
            />
            <div className="form-help">Channel buffer</div>
          </div>
          <div>
            <label className="form-label">Batch Size</label>
            <input
              type="number"
              className="form-input"
              value={config.performance?.batch_size || 50}
              onChange={(e) => handleChange('performance.batch_size', parseInt(e.target.value))}
              min="1"
            />
            <div className="form-help">Events per batch</div>
          </div>
        </div>
      </div>

      {/* BigPanda Endpoint Routing */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">BigPanda Endpoint Routing</div>
            <div className="card-subtitle">Route SNMP traps to specific BigPanda endpoints</div>
          </div>
        </div>

        {availableEndpoints.length === 0 ? (
          <div style={{
            padding: '20px',
            backgroundColor: 'rgba(245, 158, 11, 0.1)',
            borderRadius: '8px',
            border: '1px solid var(--warning)',
            fontSize: '14px',
          }}>
            <div style={{ fontWeight: 600, marginBottom: '8px' }}>No BigPanda Endpoints Configured</div>
            <div>Please configure at least one BigPanda endpoint in the BigPanda tab before setting up routing.</div>
          </div>
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
            {/* Default Endpoints */}
            <div>
              <label className="form-label">Default Endpoints</label>
              <div style={{ fontSize: '13px', color: 'var(--text-secondary)', marginBottom: '12px' }}>
                Traps that don't match any routing rules will be sent to these endpoints
              </div>
              <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
                {availableEndpoints.map(endpoint => {
                  const isSelected = (config.routing?.default_endpoints || []).includes(endpoint);
                  return (
                    <button
                      key={endpoint}
                      type="button"
                      onClick={() => {
                        const current = config.routing?.default_endpoints || [];
                        if (isSelected) {
                          // Don't allow removing if it's the last one
                          if (current.length > 1) {
                            handleChange('routing.default_endpoints', current.filter(e => e !== endpoint));
                          }
                        } else {
                          handleChange('routing.default_endpoints', [...current, endpoint]);
                        }
                      }}
                      style={{
                        padding: '8px 16px',
                        borderRadius: '6px',
                        border: isSelected ? '2px solid var(--accent-primary)' : '1px solid var(--border)',
                        backgroundColor: isSelected ? 'rgba(59, 130, 246, 0.1)' : 'var(--bg-tertiary)',
                        color: isSelected ? 'var(--accent-primary)' : 'var(--text-primary)',
                        cursor: 'pointer',
                        fontSize: '14px',
                        fontWeight: isSelected ? 600 : 400,
                      }}
                    >
                      {isSelected && '✓ '}{endpoint}
                    </button>
                  );
                })}
              </div>
            </div>

            {/* Routing Rules */}
            <div>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '12px' }}>
                <div>
                  <div style={{ fontSize: '15px', fontWeight: 600 }}>Conditional Routing Rules</div>
                  <div style={{ fontSize: '13px', color: 'var(--text-secondary)', marginTop: '4px' }}>
                    Route specific traps to different endpoints based on criteria
                  </div>
                </div>
                <button
                  type="button"
                  onClick={addRoutingRule}
                  className="btn btn-secondary"
                >
                  <Plus size={18} />
                  Add Rule
                </button>
              </div>

              {(config.routing?.rules || []).length > 0 ? (
                <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                  {config.routing.rules.map((rule, index) => (
                    <div key={index} style={{
                      padding: '16px',
                      backgroundColor: 'var(--bg-tertiary)',
                      borderRadius: '8px',
                      border: '1px solid var(--border)',
                    }}>
                      <div style={{ display: 'flex', gap: '12px', marginBottom: '12px', alignItems: 'end' }}>
                        <div style={{ flex: 2 }}>
                          <label className="form-label">Rule Name</label>
                          <input
                            type="text"
                            className="form-input"
                            value={rule.name || ''}
                            onChange={(e) => updateRoutingRule(index, 'name', e.target.value)}
                            placeholder="e.g., Cisco Devices, Network Datacenter"
                          />
                        </div>
                        <div style={{ flex: 1 }}>
                          <label className="form-label">Priority</label>
                          <input
                            type="number"
                            className="form-input"
                            value={rule.priority || 100}
                            onChange={(e) => updateRoutingRule(index, 'priority', parseInt(e.target.value))}
                            min="1"
                            max="1000"
                          />
                        </div>
                        <button
                          type="button"
                          onClick={() => removeRoutingRule(index)}
                          className="btn btn-secondary"
                          style={{ padding: '8px 12px' }}
                        >
                          <Trash2 size={18} />
                        </button>
                      </div>

                      <div style={{ display: 'flex', gap: '12px', marginBottom: '12px' }}>
                        <div style={{ flex: 1 }}>
                          <label className="form-label">Match Type</label>
                          <select
                            className="form-input"
                            value={rule.match_type || 'vendor'}
                            onChange={(e) => updateRoutingRule(index, 'match_type', e.target.value)}
                          >
                            <option value="vendor">Vendor</option>
                            <option value="oid">Exact OID</option>
                            <option value="oid_prefix">OID Prefix</option>
                            <option value="source">Source IP</option>
                            <option value="source_network">Source Network (CIDR)</option>
                          </select>
                        </div>
                        <div style={{ flex: 2 }}>
                          <label className="form-label">Match Value</label>
                          <input
                            type="text"
                            className="form-input"
                            value={rule.match_value || ''}
                            onChange={(e) => updateRoutingRule(index, 'match_value', e.target.value)}
                            placeholder={
                              rule.match_type === 'vendor' ? 'cisco, f5, netapp' :
                              rule.match_type === 'oid' ? '1.3.6.1.4.1.9.9.41.1.2.3.1' :
                              rule.match_type === 'oid_prefix' ? '1.3.6.1.4.1.9.9.41' :
                              rule.match_type === 'source' ? '10.0.0.1' :
                              '10.0.0.0/8'
                            }
                          />
                        </div>
                      </div>

                      <div>
                        <label className="form-label">Target Endpoints (select one or more)</label>
                        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px', marginTop: '8px' }}>
                          {availableEndpoints.map(endpoint => {
                            const isSelected = (rule.endpoints || []).includes(endpoint);
                            return (
                              <button
                                key={endpoint}
                                type="button"
                                onClick={() => toggleRoutingEndpoint(index, endpoint)}
                                style={{
                                  padding: '6px 12px',
                                  borderRadius: '6px',
                                  border: isSelected ? '2px solid var(--accent-primary)' : '1px solid var(--border)',
                                  backgroundColor: isSelected ? 'rgba(59, 130, 246, 0.1)' : 'transparent',
                                  color: isSelected ? 'var(--accent-primary)' : 'var(--text-primary)',
                                  cursor: 'pointer',
                                  fontSize: '13px',
                                  fontWeight: isSelected ? 600 : 400,
                                }}
                              >
                                {isSelected && '✓ '}{endpoint}
                              </button>
                            );
                          })}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div style={{
                  padding: '32px',
                  textAlign: 'center',
                  color: 'var(--text-tertiary)',
                  fontSize: '14px',
                }}>
                  No routing rules configured. All traps will be sent to the default endpoints.
                </div>
              )}
            </div>
          </div>
        )}
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

export default SNMPConfig;
