# BigPanda Super Agent - Web UI

Modern React-based web interface for the BigPanda Super Agent.

## Features

- **Dashboard**: Real-time monitoring with key metrics and charts
- **Plugins**: Manage and monitor receiver modules
- **Events**: Live event stream with WebSocket support
- **Queue**: Monitor event queue and processing status
- **SNMP**: View SNMP configurations and unknown traps
- **Webhooks**: Manage HTTP/HTTPS webhook endpoints
- **Dead Letter Queue**: Review and troubleshoot failed events

## Tech Stack

- **React 18**: Modern functional components with hooks
- **Vite**: Fast build tool and dev server
- **React Router**: Client-side routing
- **Recharts**: Beautiful data visualization
- **Lucide React**: Clean, consistent icons

## Development

### Prerequisites

- Node.js 18+ and npm

### Install Dependencies

```bash
cd web
npm install
```

### Development Server

```bash
npm run dev
```

The UI will be available at `http://localhost:3000` with proxy to the backend API at `http://localhost:8443`.

### Build for Production

```bash
npm run build
```

Build output will be in the `dist/` directory.

### Preview Production Build

```bash
npm run preview
```

## Configuration

The Vite development server is configured to proxy API requests to the backend:

- `/api/*` → `http://localhost:8443/api/*`
- `/health` → `http://localhost:8443/health`

For production, configure your web server (nginx, Apache, etc.) to proxy API requests to the agent backend.

## Project Structure

```
web/
├── src/
│   ├── pages/           # Page components
│   │   ├── Dashboard.jsx
│   │   ├── Plugins.jsx
│   │   ├── Events.jsx
│   │   ├── Queue.jsx
│   │   ├── SNMP.jsx
│   │   ├── Webhooks.jsx
│   │   └── DLQ.jsx
│   ├── components/      # Reusable components
│   ├── api.js          # API client
│   ├── App.jsx         # Main app component
│   ├── main.jsx        # Entry point
│   └── index.css       # Global styles
├── index.html
├── vite.config.js
└── package.json
```

## API Integration

The UI communicates with the backend REST API:

- **GET /health**: Health check
- **GET /api/v1/stats**: Agent statistics
- **GET /api/v1/plugins**: List plugins
- **GET /api/v1/events/recent**: Recent events
- **GET /api/v1/queue/stats**: Queue statistics
- **GET /api/v1/snmp/configs**: SNMP configurations
- **WS /api/v1/events/stream**: Live event stream

See the [API Reference](../docs/api-reference.md) for complete documentation.

## WebSocket Support

The Events page uses WebSocket for real-time event streaming. The connection is established at:

```
ws://localhost:8443/api/v1/events/stream
```

## Styling

The UI uses a custom CSS design system with:

- Dark theme optimized for monitoring
- Responsive grid layouts
- Consistent spacing and typography
- Status indicators and badges
- Clean, professional aesthetic

## Browser Support

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+

## License

Copyright © 2026 BigPanda Inc. All rights reserved.
