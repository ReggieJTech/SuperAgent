import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link, useLocation } from 'react-router-dom';
import {
  LayoutDashboard,
  Activity,
  Layers,
  Database,
  Wifi,
  Webhook,
  AlertCircle,
  Settings
} from 'lucide-react';

import Dashboard from './pages/Dashboard';
import Plugins from './pages/Plugins';
import Events from './pages/Events';
import Queue from './pages/Queue';
import SNMP from './pages/SNMP';
import Webhooks from './pages/Webhooks';
import DLQ from './pages/DLQ';
import Configuration from './pages/Configuration';

function Navigation() {
  const location = useLocation();

  const navItems = [
    { path: '/', label: 'Dashboard', icon: LayoutDashboard },
    { path: '/plugins', label: 'Plugins', icon: Layers },
    { path: '/events', label: 'Events', icon: Activity },
    { path: '/queue', label: 'Queue', icon: Database },
    { path: '/snmp', label: 'SNMP', icon: Wifi },
    { path: '/webhooks', label: 'Webhooks', icon: Webhook },
    { path: '/dlq', label: 'Dead Letter Queue', icon: AlertCircle },
    { path: '/configuration', label: 'Configuration', icon: Settings },
  ];

  return (
    <div className="sidebar">
      <div className="sidebar-header">
        <h1>BigPanda Super Agent</h1>
        <p>Monitoring & Event Management</p>
      </div>
      <ul className="nav-menu">
        {navItems.map((item) => {
          const Icon = item.icon;
          const isActive = location.pathname === item.path;
          return (
            <li key={item.path}>
              <Link
                to={item.path}
                className={`nav-item ${isActive ? 'active' : ''}`}
              >
                <Icon size={20} />
                {item.label}
              </Link>
            </li>
          );
        })}
      </ul>
    </div>
  );
}

function App() {
  return (
    <Router>
      <div className="app">
        <Navigation />
        <div className="main-content">
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/plugins" element={<Plugins />} />
            <Route path="/events" element={<Events />} />
            <Route path="/queue" element={<Queue />} />
            <Route path="/snmp" element={<SNMP />} />
            <Route path="/webhooks" element={<Webhooks />} />
            <Route path="/dlq" element={<DLQ />} />
            <Route path="/configuration" element={<Configuration />} />
          </Routes>
        </div>
      </div>
    </Router>
  );
}

export default App;
