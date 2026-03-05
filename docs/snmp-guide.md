# SNMP Module Guide

The SNMP module receives SNMP traps and forwards them as events to BigPanda.

## Features

- **SNMPv1/v2c/v3 support**: Receive traps from any SNMP version
- **Event configurations**: Transform traps using .ec files
- **MIB and NO_MIB modes**: Work with or without MIB compilation
- **Filtering**: Drop traps by OID, source IP, or network
- **Rate limiting**: Per-source and global rate limits
- **Unknown trap handling**: Log, drop, or forward unknown traps
- **Auto-reload**: Automatically reload event configs
- **Status determination**: Intelligent status mapping from trap names

## Configuration

### Main Config (config.yaml)

```yaml
modules:
  - name: snmp
    enabled: true
    config_file: "/etc/bigpanda-agent/modules/snmp.yaml"
```

### SNMP Module Config (modules/snmp.yaml)

```yaml
# SNMP trap listener
listen_address: "0.0.0.0:162"
snmp_version: "2c"          # v1, v2c, or v3
community: "public"

# SNMPv3 (if snmp_version: v3)
v3:
  security_level: "authPriv"
  auth_protocol: "SHA"
  auth_password: "${SNMP_AUTH_PASSWORD}"
  priv_protocol: "AES"
  priv_password: "${SNMP_PRIV_PASSWORD}"
  security_name: "bigpanda"

# Event configurations
event_configs_dir: "/etc/bigpanda-agent/snmp/event_configs"
mibs_dir: "/etc/bigpanda-agent/snmp/mibs"
auto_reload: true
reload_interval: 60s

# Filtering
filtering:
  enabled: true
  rules:
    - action: drop
      type: oid
      pattern: "1.3.6.1.4.1.9999.*"

    - action: drop
      type: source
      pattern: "10.0.0.100"

    - action: accept
      type: source_network
      pattern: "10.0.0.0/8"

# Unknown traps
unknown_traps:
  action: "log"              # log, drop, or forward
  log_details: true
  forward_as_critical: false

# Rate limiting
rate_limiting:
  enabled: true
  per_source: 100            # Max traps/sec per source
  global: 1000               # Max traps/sec total
  burst: 200

# Performance
performance:
  workers: 4
  buffer_size: 1000
  batch_size: 50

# Logging
logging:
  debug: false
  log_received_traps: true
  log_filtered_traps: true
  log_unknown_traps: true
```

## Event Configurations (.ec files)

Event configurations define how SNMP traps are transformed into BigPanda events.

### MIB Type

Uses MIB names directly (requires MIB compilation):

```json
{
  "type": "MIB",
  "trap": "linkDown",
  "mib": "IF-MIB",
  "copy": {},
  "rename": {
    "ifIndex": "interface_index",
    "ifDescr": "interface_name"
  },
  "set": {
    "status": "critical",
    "description": "Network interface is down"
  },
  "map-status": {},
  "primary": "snmp_source_ip",
  "secondary": "interface_name"
}
```

### NO_MIB Type

Uses OIDs and position-based varbind mapping (no MIB needed):

```json
{
  "type": "NO_MIB",
  "trap": "emsAlarm",
  "trap-oid": "1.3.6.1.4.1.21296.2.3.1.1",
  "trap-var-binds": {
    "severity": 13,
    "description": 17,
    "objectName": 9
  },
  "copy": {},
  "rename": {},
  "set": {
    "status": "critical",
    "mib_name": "INFINERA-TRAP-MIB"
  },
  "map-status": {
    "severity": {
      "1": "critical",
      "2": "warning",
      "3": "ok"
    }
  },
  "primary": "snmp_source_ip",
  "secondary": "objectName"
}
```

## Transformation Operations

### 1. Copy

Copy varbind values to new fields:

```json
"copy": {
  "ifIndex": "backup_index"
}
```

### 2. Rename

Rename varbind fields:

```json
"rename": {
  "ifIndex": "interface_index",
  "ifDescr": "interface_name"
}
```

### 3. Set

Set static field values:

```json
"set": {
  "status": "critical",
  "description": "Interface down",
  "source_system": "network"
}
```

### 4. Map Status

Map varbind values to status:

```json
"map-status": {
  "severity": {
    "critical": "critical",
    "major": "critical",
    "minor": "warning",
    "info": "ok"
  }
}
```

### 5. Primary/Secondary Keys

Define event identity:

```json
"primary": "snmp_source_ip",
"secondary": "interface_name"
```

Special values:
- `snmp_source_ip` - Source IP of trap
- `snmp_trap_name` - Trap name
- `snmp_trap_oid` - Trap OID
- Any varbind name

## Status Determination

If status is not explicitly set, it's determined from the trap name:

### Critical Patterns
`critical`, `error`, `fail`, `down`, `alarm`, `emergency`, `fatal`, `severe`, `offline`, `unreachable`, `dead`

### Warning Patterns
`warning`, `degraded`, `threshold`, `high`, `low`, `minor`, `caution`

### OK Patterns
`clear`, `ok`, `normal`, `up`, `online`, `available`, `recovered`, `restore`

## Filtering

### By OID Pattern

```yaml
- action: drop
  type: oid
  pattern: "1.3.6.1.4.1.9999.*"
```

### By Source IP

```yaml
- action: drop
  type: source
  pattern: "10.0.0.100"
```

### By Network

```yaml
- action: accept
  type: source_network
  pattern: "10.0.0.0/8"
```

## Testing

### Send Test Trap

```bash
# SNMPv2c
snmptrap -v 2c -c public localhost:162 '' \
  1.3.6.1.4.1.8072.2.3.0.1 \
  1.3.6.1.4.1.8072.2.3.2.1 i 123456

# SNMPv3
snmptrap -v 3 -l authPriv \
  -u bigpanda \
  -a SHA -A authpass \
  -x AES -X privpass \
  localhost:162 '' \
  1.3.6.1.4.1.8072.2.3.0.1 \
  1.3.6.1.4.1.8072.2.3.2.1 i 123456
```

### View Statistics

```bash
curl http://localhost:8443/stats | jq '.plugins.snmp'
```

### Check Logs

```bash
tail -f /var/log/bigpanda-agent/agent.log | grep snmp
```

## Pre-bundled Vendors

The agent ships with event configurations for 60+ vendors:

- Cisco (IOS, NXOS, ASA, UCS)
- F5 BIG-IP
- NetApp
- Dell EMC
- HP/HPE
- Juniper
- Palo Alto Networks
- Check Point
- Fortinet
- Arista
- And many more...

## Adding Custom Configs

1. Create .ec file in event configs directory:

```bash
/etc/bigpanda-agent/snmp/event_configs/custom.ec
```

2. Add trap definitions (JSON array)

3. Reload agent or wait for auto-reload:

```bash
sudo systemctl reload bigpanda-agent
```

## Troubleshooting

### No traps received

1. Check firewall allows UDP 162:
```bash
sudo ufw allow 162/udp
```

2. Check listener is running:
```bash
sudo netstat -ulnp | grep 162
```

3. Enable debug logging:
```yaml
logging:
  debug: true
  log_received_traps: true
```

### Unknown traps

1. Check event config exists:
```bash
grep -r "trap_oid_here" /etc/bigpanda-agent/snmp/event_configs/
```

2. View unknown traps in logs:
```bash
grep "Unknown trap" /var/log/bigpanda-agent/agent.log
```

3. Create event config for the trap

### Rate limiting

Check rate limit stats:
```bash
curl http://localhost:8443/stats | jq '.plugins.snmp.traps_dropped'
```

Adjust limits in config:
```yaml
rate_limiting:
  per_source: 200
  global: 2000
```

## Best Practices

1. **Use NO_MIB type** when possible (faster, no MIB dependency)
2. **Pre-bundle common configs** in your deployment
3. **Enable filtering** to drop noisy traps
4. **Set rate limits** to prevent overload
5. **Monitor unknown traps** and add configs as needed
6. **Use auto-reload** to pick up config changes
7. **Test configs** before deploying to production

## Performance

- Handles 1000+ traps/sec sustained
- Low latency (< 10ms processing)
- Concurrent trap processing
- Rate limiting prevents overload
- Persistent queue prevents data loss

## Security

- SNMPv3 encryption support
- Community string configuration
- IP-based filtering
- No shell execution
- Sandboxed processing
