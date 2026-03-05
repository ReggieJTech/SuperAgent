import React, { useState, useEffect } from 'react';
import { Wifi, Search, Filter } from 'lucide-react';
import api from '../api';

function SNMP() {
  const [configs, setConfigs] = useState([]);
  const [unknownTraps, setUnknownTraps] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedVendor, setSelectedVendor] = useState('all');

  useEffect(() => {
    loadData();
    const interval = setInterval(() => loadUnknownTraps(), 10000); // Refresh unknown traps every 10s
    return () => clearInterval(interval);
  }, []);

  const loadData = async () => {
    try {
      const [configsData, unknownData] = await Promise.all([
        api.getSNMPConfigs(),
        api.getUnknownSNMPTraps(),
      ]);

      // Backend returns {configs: [], count: 0} and {traps: [], count: 0}
      const configList = configsData?.configs || configsData || [];
      const trapList = unknownData?.traps || unknownData || [];

      setConfigs(Array.isArray(configList) ? configList : []);
      setUnknownTraps(Array.isArray(trapList) ? trapList : []);
      setLoading(false);
      setError(null);
    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  };

  const loadUnknownTraps = async () => {
    try {
      const data = await api.getUnknownSNMPTraps();
      // Backend returns {traps: [], count: 0}
      const trapList = data?.traps || data || [];
      setUnknownTraps(Array.isArray(trapList) ? trapList : []);
    } catch (err) {
      console.error('Failed to load unknown traps:', err);
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
          <h2>SNMP</h2>
          <p>SNMP trap configurations and monitoring</p>
        </div>
        <div className="error-message">
          Failed to load SNMP data: {error}
        </div>
      </div>
    );
  }

  // Get unique vendors
  const vendors = ['all', ...new Set(configs.map(c => c.vendor).filter(Boolean))];

  // Filter configs
  const filteredConfigs = configs.filter(config => {
    const matchesSearch = !searchTerm ||
      config.oid?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      config.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      config.vendor?.toLowerCase().includes(searchTerm.toLowerCase());

    const matchesVendor = selectedVendor === 'all' || config.vendor === selectedVendor;

    return matchesSearch && matchesVendor;
  });

  // Group configs by vendor
  const configsByVendor = filteredConfigs.reduce((acc, config) => {
    const vendor = config.vendor || 'Unknown';
    if (!acc[vendor]) {
      acc[vendor] = [];
    }
    acc[vendor].push(config);
    return acc;
  }, {});

  return (
    <div>
      <div className="page-header">
        <h2>SNMP</h2>
        <p>SNMP trap configurations and monitoring</p>
      </div>

      {/* Summary Stats */}
      <div className="grid grid-4">
        <div className="stat-card">
          <div className="stat-label">Total Configurations</div>
          <div className="stat-value">{configs.length}</div>
          <div className="stat-change">Event mappings loaded</div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Vendors</div>
          <div className="stat-value">{vendors.length - 1}</div>
          <div className="stat-change">Supported vendors</div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Unknown Traps</div>
          <div className="stat-value" style={{
            color: unknownTraps.length > 0 ? 'var(--warning)' : 'var(--text-primary)'
          }}>
            {unknownTraps.length}
          </div>
          <div className="stat-change">Need configuration</div>
        </div>
        <div className="stat-card">
          <div className="stat-label">Active Filters</div>
          <div className="stat-value">{searchTerm || selectedVendor !== 'all' ? '1+' : '0'}</div>
          <div className="stat-change">Search criteria</div>
        </div>
      </div>

      {/* Unknown Traps Alert */}
      {unknownTraps.length > 0 && (
        <div className="card" style={{
          backgroundColor: 'rgba(245, 158, 11, 0.1)',
          borderColor: 'var(--warning)'
        }}>
          <div className="card-header">
            <div>
              <div className="card-title">Unknown SNMP Traps Detected</div>
              <div className="card-subtitle">
                {unknownTraps.length} trap(s) received without matching configuration
              </div>
            </div>
          </div>
          <div className="table-container">
            <table>
              <thead>
                <tr>
                  <th>OID</th>
                  <th>Source IP</th>
                  <th>First Seen</th>
                  <th>Count</th>
                </tr>
              </thead>
              <tbody>
                {unknownTraps.slice(0, 5).map((trap, index) => (
                  <tr key={index}>
                    <td><code>{trap.oid}</code></td>
                    <td>{trap.source_ip}</td>
                    <td>{new Date(trap.first_seen).toLocaleString()}</td>
                    <td>{trap.count}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {unknownTraps.length > 5 && (
            <div style={{ marginTop: '12px', fontSize: '14px', color: 'var(--text-secondary)' }}>
              ...and {unknownTraps.length - 5} more
            </div>
          )}
        </div>
      )}

      {/* Search and Filter */}
      <div className="card">
        <div className="card-header">
          <div>
            <div className="card-title">SNMP Event Configurations</div>
            <div className="card-subtitle">{filteredConfigs.length} configurations</div>
          </div>
        </div>

        <div style={{ display: 'flex', gap: '12px', marginBottom: '20px' }}>
          <div style={{ flex: 1, position: 'relative' }}>
            <Search
              size={18}
              style={{
                position: 'absolute',
                left: '12px',
                top: '50%',
                transform: 'translateY(-50%)',
                color: 'var(--text-tertiary)'
              }}
            />
            <input
              type="text"
              placeholder="Search by OID, name, or vendor..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              style={{
                width: '100%',
                padding: '10px 12px 10px 40px',
                backgroundColor: 'var(--bg-tertiary)',
                border: '1px solid var(--border)',
                borderRadius: '8px',
                color: 'var(--text-primary)',
                fontSize: '14px',
              }}
            />
          </div>
          <div style={{ position: 'relative' }}>
            <Filter
              size={18}
              style={{
                position: 'absolute',
                left: '12px',
                top: '50%',
                transform: 'translateY(-50%)',
                color: 'var(--text-tertiary)',
                pointerEvents: 'none'
              }}
            />
            <select
              value={selectedVendor}
              onChange={(e) => setSelectedVendor(e.target.value)}
              style={{
                padding: '10px 12px 10px 40px',
                backgroundColor: 'var(--bg-tertiary)',
                border: '1px solid var(--border)',
                borderRadius: '8px',
                color: 'var(--text-primary)',
                fontSize: '14px',
                cursor: 'pointer',
                minWidth: '200px',
              }}
            >
              {vendors.map(vendor => (
                <option key={vendor} value={vendor}>
                  {vendor === 'all' ? 'All Vendors' : vendor}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* Configs by Vendor */}
        {Object.entries(configsByVendor).map(([vendor, vendorConfigs]) => (
          <div key={vendor} style={{ marginBottom: '30px' }}>
            <h3 style={{
              fontSize: '18px',
              fontWeight: 600,
              marginBottom: '16px',
              display: 'flex',
              alignItems: 'center',
              gap: '8px'
            }}>
              <Wifi size={20} />
              {vendor}
              <span className="badge badge-info">{vendorConfigs.length}</span>
            </h3>

            <div className="table-container">
              <table>
                <thead>
                  <tr>
                    <th>OID</th>
                    <th>Name</th>
                    <th>Status Mapping</th>
                    <th>Priority</th>
                    <th>Tags</th>
                  </tr>
                </thead>
                <tbody>
                  {vendorConfigs.map((config, index) => (
                    <tr key={index}>
                      <td><code style={{ fontSize: '12px' }}>{config.oid}</code></td>
                      <td>{config.name || '-'}</td>
                      <td>
                        {config.status_map ? (
                          <span className="badge badge-success">✓ Configured</span>
                        ) : (
                          <span className="badge badge-info">Default</span>
                        )}
                      </td>
                      <td>
                        {config.priority && (
                          <span className={`badge ${
                            config.priority === 'critical' ? 'badge-error' :
                            config.priority === 'warning' ? 'badge-warning' :
                            'badge-info'
                          }`}>
                            {config.priority}
                          </span>
                        )}
                      </td>
                      <td style={{ fontSize: '12px' }}>
                        {config.tags && Object.keys(config.tags).length > 0 ? (
                          <div style={{ display: 'flex', flexWrap: 'wrap', gap: '4px' }}>
                            {Object.entries(config.tags).slice(0, 2).map(([key, value]) => (
                              <span key={key} className="badge badge-info" style={{ fontSize: '10px' }}>
                                {key}:{value}
                              </span>
                            ))}
                            {Object.keys(config.tags).length > 2 && (
                              <span className="badge badge-info" style={{ fontSize: '10px' }}>
                                +{Object.keys(config.tags).length - 2}
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
          </div>
        ))}

        {filteredConfigs.length === 0 && (
          <div className="empty-state">
            <div className="empty-state-icon">
              <Search size={48} />
            </div>
            <div>No configurations match your search criteria</div>
          </div>
        )}
      </div>
    </div>
  );
}

export default SNMP;
