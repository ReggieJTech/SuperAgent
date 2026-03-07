# Multi-Endpoint Routing Configuration

## Overview

The BigPanda SuperAgent now supports routing events from different sources (SNMP traps, webhooks) to multiple BigPanda integration endpoints. This enables:

1. **Multiple BigPanda Instances**: Send events to PROD and TEST BigPanda orgs simultaneously
2. **Per-Source Routing**: Network monitoring → Network app_key, Server monitoring → Server app_key
3. **Multi-Destination**: Fan-out events to multiple endpoints for redundancy or different teams
4. **Flexible Routing Rules**: Route based on OID, vendor, source IP, webhook source, etc.

## Configuration Structure

### Main Configuration

**File:** `config.yaml` (or `/etc/bigpanda-agent/config.yaml`)

```yaml
# Multiple BigPanda Endpoints
bigpanda_endpoints:
  # Production Network Monitoring
  - name: "prod-network"
    description: "Production BigPanda - Network Team"
    enabled: true
    api_url: "https://api.bigpanda.io/data/v2/alerts"
    token: "${BP_TOKEN_PROD}"
    app_key: "network_monitoring_key"
    batching:
      enabled: true
      max_size: 100
      max_wait: 5s
    retry:
      max_attempts: 5
      initial_backoff: 1s
      max_backoff: 60s
    rate_limit:
      events_per_second: 1000
      burst: 2000
    tags:
      environment: "production"
      team: "network"

  # Production Server Monitoring
  - name: "prod-servers"
    description: "Production BigPanda - Server Team"
    enabled: true
    api_url: "https://api.bigpanda.io/data/v2/alerts"
    token: "${BP_TOKEN_PROD}"
    app_key: "server_monitoring_key"
    batching:
      enabled: true
      max_size: 100
      max_wait: 5s
    tags:
      environment: "production"
      team: "servers"

  # Test/Staging Environment
  - name: "test-all"
    description: "Test BigPanda Instance"
    enabled: true
    api_url: "https://api-test.bigpanda.io/data/v2/alerts"  # Different instance
    token: "${BP_TOKEN_TEST}"
    app_key: "test_monitoring_key"
    batching:
      enabled: true
      max_size: 50
      max_wait: 10s
    tags:
      environment: "test"

  # EU Region Instance
  - name: "prod-eu"
    description: "EU Production BigPanda"
    enabled: true
    api_url: "https://api.eu.bigpanda.io/data/v2/alerts"  # EU endpoint
    token: "${BP_TOKEN_EU}"
    app_key: "eu_monitoring_key"
    tags:
      region: "eu"
      environment: "production"
```

### Backward Compatibility

The legacy single-endpoint configuration is still supported and automatically migrated:

```yaml
# Legacy format (still works)
bigpanda:
  api_url: "https://api.bigpanda.io/data/v2/alerts"
  token: "your-token"
  app_key: "your-app-key"
  # ... other settings
```

This is automatically converted to:
```yaml
bigpanda_endpoints:
  - name: "default"
    description: "Default BigPanda endpoint (migrated from legacy config)"
    enabled: true
    # ... legacy settings copied here
```

## SNMP Routing Configuration

**File:** `snmp.yaml` (or `/etc/bigpanda-agent/modules/snmp.yaml`)

```yaml
# SNMP Listener Settings
listen_address: "0.0.0.0:162"
snmp_version: "2c"
community: "public"

# Event Routing
routing:
  # Default endpoint(s) for all traps (if no rules match)
  default_endpoints:
    - "prod-network"
    - "test-all"  # Send all SNMP to both prod and test

  # Conditional routing rules (evaluated by priority, highest first)
  rules:
    # Cisco network devices → Network team + EU instance
    - name: "Cisco Devices"
      match_type: "vendor"
      match_value: "cisco"
      endpoints:
        - "prod-network"
        - "prod-eu"
      priority: 100

    # F5 load balancers → Server team
    - name: "F5 Load Balancers"
      match_type: "vendor"
      match_value: "f5"
      endpoints:
        - "prod-servers"
      priority: 90

    # Critical infrastructure (specific OID prefix) → Multiple teams
    - name: "Critical Infrastructure"
      match_type: "oid_prefix"
      match_value: "1.3.6.1.4.1.9.9.41"  # Cisco critical alerts
      endpoints:
        - "prod-network"
        - "prod-servers"  # Notify both teams
      priority: 200

    # Specific datacenter (by source network) → EU instance only
    - name: "EU Datacenter"
      match_type: "source_network"
      match_value: "10.100.0.0/16"
      endpoints:
        - "prod-eu"
      priority: 150

    # Specific device (by source IP)
    - name: "Core Router"
      match_type: "source"
      match_value: "10.0.0.1"
      endpoints:
        - "prod-network"
        - "test-all"  # Also send to test for analysis
      priority: 250

# Event configurations, filtering, rate limiting, etc.
event_configs_dir: "/etc/bigpanda-agent/snmp/event_configs"
# ... rest of SNMP config
```

### SNMP Routing Match Types

| Match Type | Description | Example |
|------------|-------------|---------|
| `vendor` | Match by vendor name from event config | `cisco`, `f5`, `netapp` |
| `oid` | Exact OID match | `1.3.6.1.4.1.9.9.41.1.2.3.1` |
| `oid_prefix` | OID prefix match | `1.3.6.1.4.1.9.9.41` |
| `source` | Exact source IP | `10.0.0.1` |
| `source_network` | Source IP in CIDR range | `10.0.0.0/8` |

### Routing Rule Priority

- Rules are evaluated in **descending priority** order (highest first)
- First matching rule determines the endpoint(s)
- If no rules match, `default_endpoints` are used
- Multiple endpoints can be specified for fan-out delivery

## Webhook Routing Configuration

**File:** `webhook.yaml` (or `/etc/bigpanda-agent/modules/webhook.yaml`)

```yaml
# Webhook Listener Settings
listen_address: "0.0.0.0:8080"

# Webhook Sources with Endpoint Routing
sources:
  # Prometheus AlertManager → Network Team
  - name: "prometheus-network"
    enabled: true
    path: "/webhook/prometheus/network"
    method: "POST"
    auth:
      type: "bearer"
      token: "${PROMETHEUS_TOKEN}"
    endpoints:
      - "prod-network"
      - "test-all"  # Also send to test
    transform:
      field_map:
        "alertname": "check"
        "instance": "host"
        "severity": "severity"
      status_map:
        "firing": "critical"
        "resolved": "ok"
      primary_key: "host"
      secondary_key: "check"

  # Datadog → Server Team
  - name: "datadog-servers"
    enabled: true
    path: "/webhook/datadog/servers"
    method: "POST"
    auth:
      type: "apikey"
      header: "X-API-Key"
      key: "${DATADOG_KEY}"
    endpoints:
      - "prod-servers"
    transform:
      field_map:
        "title": "check"
        "host": "host"
        "priority": "severity"

  # SolarWinds → Network Team (EU and US)
  - name: "solarwinds-network"
    enabled: true
    path: "/webhook/solarwinds"
    method: "POST"
    auth:
      type: "basic"
      username: "solarwinds"
      password: "${SOLARWINDS_PASSWORD}"
    endpoints:
      - "prod-network"
      - "prod-eu"  # Fan-out to both regions
    transform:
      field_map:
        "NodeName": "host"
        "AlertName": "check"
        "Severity": "severity"

  # Generic webhook → Default routing
  - name: "generic"
    enabled: true
    path: "/webhook/generic"
    method: "POST"
    auth:
      type: "none"
    endpoints:
      - "prod-network"  # Default if not specified elsewhere
    transform:
      field_map:
        "title": "check"
        "host": "host"
```

### Webhook Routing Strategy

Each webhook source explicitly specifies which BigPanda endpoint(s) to route to:

- **Single Endpoint**: `endpoints: ["prod-network"]`
- **Multi-Destination**: `endpoints: ["prod-network", "test-all"]`
- **Default Fallback**: If `endpoints` not specified, uses `["default"]`

## Use Cases

### 1. PROD + TEST Dual Delivery

Send all events to both production and test BigPanda instances:

```yaml
# SNMP
routing:
  default_endpoints:
    - "prod-network"
    - "test-all"

# Webhook
endpoints:
  - "prod-network"
  - "test-all"
```

### 2. Team-Based Routing

Route different event types to different teams:

```yaml
bigpanda_endpoints:
  - name: "team-network"
    app_key: "network_app_key"
  - name: "team-servers"
    app_key: "servers_app_key"
  - name: "team-security"
    app_key: "security_app_key"

# Route by vendor/source
routing:
  rules:
    - name: "Network Devices"
      match_type: "vendor"
      match_value: "cisco|juniper|arista"
      endpoints: ["team-network"]

    - name: "Servers"
      match_type: "vendor"
      match_value: "dell|hp|ibm"
      endpoints: ["team-servers"]
```

### 3. Regional Routing

Route events by geographic region:

```yaml
bigpanda_endpoints:
  - name: "us-prod"
    api_url: "https://api.bigpanda.io/data/v2/alerts"
  - name: "eu-prod"
    api_url: "https://api.eu.bigpanda.io/data/v2/alerts"

# Route by source network
routing:
  rules:
    - name: "US Datacenter"
      match_type: "source_network"
      match_value: "10.10.0.0/16"
      endpoints: ["us-prod"]

    - name: "EU Datacenter"
      match_type: "source_network"
      match_value: "10.100.0.0/16"
      endpoints: ["eu-prod"]
```

### 4. Critical Event Escalation

Send critical events to multiple teams:

```yaml
routing:
  rules:
    - name: "Critical Infrastructure"
      match_type: "oid_prefix"
      match_value: "1.3.6.1.4.1.CRITICAL"
      endpoints:
        - "team-network"
        - "team-servers"
        - "team-management"
      priority: 999
```

## Event Flow

```
┌─────────────────┐
│  SNMP Trap or   │
│  Webhook Event  │
└────────┬────────┘
         │
         ▼
┌─────────────────────────┐
│  Plugin (SNMP/Webhook)  │
│  - Receives event       │
│  - Applies routing      │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│  Routing Engine         │
│  - Evaluates rules      │
│  - Determines endpoint  │
└────────┬────────────────┘
         │
         ├────────┬────────┬────────┐
         ▼        ▼        ▼        ▼
     ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐
     │ Prod │ │ Test │ │  EU  │ │ Team │
     │  BP  │ │  BP  │ │  BP  │ │  BP  │
     └──────┘ └──────┘ └──────┘ └──────┘
```

## Configuration Validation

The agent validates:

1. **Unique endpoint names** - No duplicate `name` fields
2. **Required fields** - Each enabled endpoint must have `api_url`, `token`, `app_key`
3. **Valid routing references** - `endpoints` must reference existing endpoint names
4. **Default fallback** - At least one default endpoint exists

## Migration from Single Endpoint

Existing configurations with single `bigpanda:` section are automatically migrated:

**Before:**
```yaml
bigpanda:
  api_url: "https://api.bigpanda.io/data/v2/alerts"
  token: "token123"
  app_key: "app_key_123"
```

**After (automatic):**
```yaml
bigpanda_endpoints:
  - name: "default"
    enabled: true
    api_url: "https://api.bigpanda.io/data/v2/alerts"
    token: "token123"
    app_key: "app_key_123"
```

All existing SNMP and webhook configurations automatically route to the `default` endpoint.

## Benefits

1. **Flexibility**: Different event sources → different BigPanda integrations
2. **Redundancy**: Send events to multiple BigPanda instances
3. **Testing**: Easy PROD + TEST dual delivery
4. **Team Separation**: Different teams get different app_keys
5. **Regional Support**: Route to region-specific BigPanda instances
6. **Gradual Migration**: Test new configurations without disrupting production

## Next Steps

- [ ] Update forwarder to support multiple BigPanda clients
- [ ] Implement routing logic in SNMP and Webhook plugins
- [ ] Add UI for managing endpoints and routing rules
- [ ] Add connection testing for each endpoint
- [ ] Implement per-endpoint statistics and monitoring
