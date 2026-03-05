# API Reference

The BigPanda Super Agent provides a RESTful API for management and monitoring.

## Base URL

- HTTP: `http://localhost:8443/api/v1`
- HTTPS: `https://localhost:8443/api/v1`

## Authentication

Most endpoints require JWT authentication. Include the token in the `Authorization` header:

```
Authorization: Bearer <token>
```

### Login

**POST** `/api/v1/auth/login`

Authenticate and receive a JWT token.

**Request:**
```json
{
  "username": "admin",
  "password": "password"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2026-03-05T12:00:00Z",
    "user": {
      "username": "admin",
      "role": "admin"
    }
  },
  "time": "2026-03-04T12:00:00Z"
}
```

### Logout

**POST** `/api/v1/auth/logout`

Logout (client-side token invalidation).

### Refresh Token

**POST** `/api/v1/auth/refresh`

Refresh an expiring token.

## Health Endpoints

### Health Check

**GET** `/health`

Returns overall agent health status.

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "started": true,
    "uptime": "2h30m15s",
    "queue": {
      "status": "healthy",
      "size": 42
    },
    "forwarder": {
      "status": "healthy",
      "circuit_breaker": "closed"
    },
    "plugins": {
      "status": "healthy",
      "total_plugins": 1
    }
  }
}
```

### Liveness Probe

**GET** `/health/live`

Returns 200 if the server is alive.

### Readiness Probe

**GET** `/health/ready`

Returns 200 if the agent is ready to serve traffic.

## Statistics

### Get Stats

**GET** `/api/v1/stats`

Returns comprehensive statistics.

**Response:**
```json
{
  "success": true,
  "data": {
    "uptime": "2h30m15s",
    "start_time": "2026-03-04T09:30:00Z",
    "queue": {
      "enqueued": 1523,
      "dequeued": 1481,
      "dropped": 0,
      "size": 42
    },
    "forwarder": {
      "sent": 1481,
      "failed": 0,
      "retried": 5,
      "batches": 15
    },
    "plugins": {
      "totals": {
        "events_received": 1523,
        "events_sent": 1523,
        "events_dropped": 0,
        "errors": 0
      },
      "plugins": {
        "snmp": {
          "events_received": 1523,
          "events_sent": 1523,
          "traps_filtered": 12,
          "traps_unknown": 3,
          "started": true,
          "uptime": "2h30m10s"
        }
      }
    }
  }
}
```

## Agent Management

### Get Agent Info

**GET** `/api/v1/agent/info`

Returns agent information.

**Response:**
```json
{
  "success": true,
  "data": {
    "started": true,
    "uptime": "2h30m15s"
  }
}
```

### Get Agent Config

**GET** `/api/v1/agent/config`

Returns agent configuration (sensitive values redacted).

## Queue Management

### Get Queue Stats

**GET** `/api/v1/queue/stats`

Returns queue statistics.

**Response:**
```json
{
  "success": true,
  "data": {
    "enqueued": 1523,
    "dequeued": 1481,
    "dropped": 0,
    "size": 42,
    "maxSize": 100000
  }
}
```

### Get Queue Size

**GET** `/api/v1/queue/size`

Returns current queue size.

**Response:**
```json
{
  "success": true,
  "data": {
    "size": 42
  }
}
```

## Plugin Management

### List Plugins

**GET** `/api/v1/plugins`

Returns list of loaded plugins.

**Response:**
```json
{
  "success": true,
  "data": {
    "plugins": ["snmp"]
  }
}
```

### Get Plugin Info

**GET** `/api/v1/plugins/{name}`

Returns plugin information.

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "message": "Plugin is operating normally"
  }
}
```

### Get Plugin Stats

**GET** `/api/v1/plugins/{name}/stats`

Returns plugin statistics.

**Response:**
```json
{
  "success": true,
  "data": {
    "events_received": 1523,
    "events_sent": 1523,
    "events_dropped": 0,
    "errors": 0,
    "started": true,
    "uptime": "2h30m10s",
    "traps_filtered": 12,
    "traps_unknown": 3
  }
}
```

## Event Management

### Get Recent Events

**GET** `/api/v1/events/recent`

Returns recently processed events.

**Query Parameters:**
- `limit` (int): Number of events to return (default: 100)
- `offset` (int): Offset for pagination (default: 0)

**Response:**
```json
{
  "success": true,
  "data": {
    "events": [],
    "count": 0
  }
}
```

### Get DLQ Events

**GET** `/api/v1/events/dlq`

Returns dead letter queue events.

**Response:**
```json
{
  "success": true,
  "data": {
    "events": [],
    "count": 0
  }
}
```

### Event Stream (WebSocket)

**WS** `/api/v1/events/stream`

WebSocket endpoint for real-time event streaming.

**Message Format:**
```json
{
  "type": "event",
  "timestamp": "2026-03-04T12:00:00Z",
  "data": {
    "status": "critical",
    "primary_property": "192.168.1.1",
    "secondary_property": "linkDown"
  }
}
```

## SNMP Module

### List Event Configs

**GET** `/api/v1/snmp/configs`

Returns loaded SNMP event configurations.

**Response:**
```json
{
  "success": true,
  "data": {
    "configs": ["cisco-general", "f5-bigip", "netapp"],
    "count": 3
  }
}
```

### Get Unknown Traps

**GET** `/api/v1/snmp/unknown`

Returns recent unknown traps.

**Response:**
```json
{
  "success": true,
  "data": {
    "traps": [],
    "count": 0
  }
}
```

## Response Format

All API responses follow this format:

```json
{
  "success": true|false,
  "data": { ... },
  "error": "error message" (if success: false),
  "time": "2026-03-04T12:00:00Z"
}
```

### HTTP Status Codes

- `200 OK`: Successful request
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Authentication required or failed
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error
- `503 Service Unavailable`: Agent not ready

## Rate Limiting

API requests are rate limited to:
- 100 requests per minute per IP
- Burst of 20 requests

Rate limit headers:
- `X-RateLimit-Limit`: Request limit per minute
- `X-RateLimit-Remaining`: Remaining requests
- `X-RateLimit-Reset`: Time when limit resets

## CORS

Cross-Origin Resource Sharing (CORS) is enabled for all origins by default.

Configure in `config.yaml`:
```yaml
server:
  cors:
    enabled: true
    allowed_origins:
      - "https://dashboard.example.com"
```

## TLS/HTTPS

Enable HTTPS in `config.yaml`:

```yaml
server:
  listen_address: "0.0.0.0:8443"
  tls:
    enabled: true
    cert_file: "/etc/bigpanda-agent/certs/server.crt"
    key_file: "/etc/bigpanda-agent/certs/server.key"
    auto_generate: true
```

## Examples

### Using curl

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:8443/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' \
  | jq -r '.data.token')

# Get stats
curl -s http://localhost:8443/api/v1/stats \
  -H "Authorization: Bearer $TOKEN" \
  | jq .

# Get plugin info
curl -s http://localhost:8443/api/v1/plugins/snmp \
  -H "Authorization: Bearer $TOKEN" \
  | jq .
```

### Using JavaScript

```javascript
// Login
const response = await fetch('http://localhost:8443/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: 'admin', password: 'admin' })
});
const { data } = await response.json();
const token = data.token;

// Get stats
const statsResponse = await fetch('http://localhost:8443/api/v1/stats', {
  headers: { 'Authorization': `Bearer ${token}` }
});
const stats = await statsResponse.json();
```

### WebSocket Connection

```javascript
const ws = new WebSocket('ws://localhost:8443/api/v1/events/stream');

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received event:', message);
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};
```

## Error Handling

Errors are returned in this format:

```json
{
  "success": false,
  "error": "detailed error message",
  "time": "2026-03-04T12:00:00Z"
}
```

Common errors:
- `"authentication required"`: Missing or invalid token
- `"plugin not found"`: Invalid plugin name
- `"agent not started"`: Agent is not running
- `"invalid request"`: Malformed request body

## Support

- Documentation: https://docs.bigpanda.io/super-agent/api
- Issues: https://github.com/bigpanda/super-agent/issues
