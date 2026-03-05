# Pre-bundled Event Configurations

This directory contains pre-bundled event configurations (.ec files) for common network vendors and devices.

## What are Event Configurations?

Event configurations define how SNMP traps are transformed into BigPanda events. Each .ec file contains a JSON array of trap definitions with transformation rules.

## Included Vendors

### Network Infrastructure

1. **cisco-general.ec** - Cisco Systems
   - Standard SNMP traps (linkDown, linkUp, coldStart, warmStart, authenticationFailure)
   - Works with: Cisco IOS, IOS-XE, IOS-XR, NX-OS, ASA

2. **juniper.ec** - Juniper Networks
   - Power supply failures
   - Fan failures
   - BGP peer state changes
   - Works with: Juniper MX, EX, QFX, SRX series

3. **arista.ec** - Arista Networks
   - Fan status changes
   - Power supply status changes
   - Works with: Arista 7000, 7500 series switches

### Security Appliances

4. **f5-bigip.ec** - F5 Networks BIG-IP
   - Agent start/shutdown
   - Node up/down events
   - Works with: F5 BIG-IP LTM, GTM, ASM

5. **palo-alto.ec** - Palo Alto Networks
   - System events
   - GlobalProtect gateway status
   - Works with: PA-Series firewalls

6. **checkpoint.ec** - Check Point Software
   - Threshold alerts
   - Firewall module state
   - Works with: Check Point Security Gateway

7. **fortinet.ec** - Fortinet FortiGate
   - CPU/Memory thresholds
   - VPN tunnel status
   - Works with: FortiGate firewalls

### Storage Systems

8. **netapp.ec** - NetApp
   - Disk failures
   - Volume full/nearly full
   - Fan failures
   - Works with: NetApp FAS, AFF systems

9. **dell-emc.ec** - Dell EMC
   - Fan failures
   - Power supply failures
   - Temperature probe alerts
   - Works with: Dell PowerEdge servers, EMC storage

### Servers & Hardware

10. **hp-hpe.ec** - HP/HPE
    - Array accelerator status
    - Temperature failures
    - Power supply failures
    - Works with: HP ProLiant, HPE servers

11. **Ibm2100-MIB.ec** - IBM 2100 Series
    - IBM-specific traps
    - Works with: IBM System Storage

### Specialized Equipment

12. **arbor.ec** - Arbor Networks
    - DDoS protection traps
    - Works with: Arbor Pravail, TMS

13. **I3IC-MIB.ec** - Interactive Intelligence (Genesys)
    - Call center infrastructure
    - System restarts
    - Event log entries
    - Switchover events
    - Works with: Interaction Center

14. **INFINERA-TRAP-MIB.ec** - Infinera
    - Optical networking equipment
    - Alarms and audits
    - Works with: Infinera DTN-X, XTM series

15. **PDU2-MIB.ec** - Power Distribution Units
    - PDU-specific traps
    - Power monitoring

16. **VERITAS-COMMAND-CENTRAL-MIB.ec** - Veritas Command Central
    - Backup and storage management
    - Works with: Veritas NetBackup, Backup Exec

17. **VERITAS-REG.ec** - Veritas Registry
18. **VERITAS-TC.ec** - Veritas Textual Conventions

## File Format

Each .ec file contains a JSON array of event configurations:

```json
[
  {
    "type": "NO_MIB",
    "trap": "linkDown",
    "trap-oid": "1.3.6.1.6.3.1.1.5.3",
    "trap-var-binds": {
      "ifIndex": 1,
      "ifDescr": 2
    },
    "rename": {
      "ifIndex": "interface_index",
      "ifDescr": "interface_name"
    },
    "set": {
      "status": "critical",
      "description": "Network interface is down"
    },
    "primary": "snmp_source_ip",
    "secondary": "interface_name"
  }
]
```

## Configuration Types

### MIB Type
Requires MIB compilation, uses MIB names directly:
```json
{
  "type": "MIB",
  "trap": "trapName",
  "mib": "VENDOR-MIB"
}
```

### NO_MIB Type (Recommended)
No MIB compilation needed, uses OIDs and position-based varbinds:
```json
{
  "type": "NO_MIB",
  "trap": "trapName",
  "trap-oid": "1.3.6.1.4.1.xxx",
  "trap-var-binds": {
    "field1": 1,
    "field2": 2
  }
}
```

## Transformation Operations

Event configs support these transformation operations:

1. **copy** - Duplicate field values
2. **rename** - Rename varbind fields
3. **set** - Set static field values (status, description, etc.)
4. **map-status** - Map varbind values to event status
5. **primary/secondary** - Define event identity keys

## Adding Custom Configurations

To add your own event configurations:

1. Create a new .ec file in this directory:
   ```bash
   /etc/bigpanda-agent/snmp/event_configs/custom.ec
   ```

2. Add trap definitions in JSON format

3. The agent will auto-reload configs (if enabled) or restart:
   ```bash
   sudo systemctl reload bigpanda-agent
   ```

## Testing

Test a configuration by sending a trap:

```bash
snmptrap -v 2c -c public localhost:162 '' \
  1.3.6.1.6.3.1.1.5.3 \
  1.3.6.1.2.1.2.2.1.1 i 1 \
  1.3.6.1.2.1.2.2.1.2 s "eth0"
```

Check if the trap was processed:

```bash
tail -f /var/log/bigpanda-agent/agent.log | grep snmp
```

## Vendor Coverage

This package provides coverage for:
- ✅ Network switches and routers (Cisco, Juniper, Arista)
- ✅ Firewalls (Palo Alto, Check Point, Fortinet)
- ✅ Load balancers (F5 BIG-IP)
- ✅ Storage systems (NetApp, Dell EMC)
- ✅ Servers (HP/HPE, Dell, IBM)
- ✅ Optical networking (Infinera)
- ✅ Call centers (Genesys)
- ✅ Security (Arbor DDoS)
- ✅ Backup systems (Veritas)
- ✅ Power distribution (PDU)

## Need More?

If you need event configurations for a vendor not listed here:

1. **Check BigPanda documentation** for additional configs
2. **Create your own** using the format above
3. **Contact BigPanda support** for assistance
4. **Use the MIB converter** (future feature) to generate configs from MIB files

## Status Determination

If status is not explicitly set in the config, it's automatically determined from the trap name:

- **Critical**: error, fail, down, critical, alarm, emergency, fatal, severe, offline, unreachable, dead
- **Warning**: warning, degraded, threshold, high, low, minor, caution
- **OK**: clear, ok, normal, up, online, available, recovered, restore, resolved

## Best Practices

1. **Use NO_MIB type** whenever possible (faster, no dependencies)
2. **Set explicit status** for known trap severities
3. **Use meaningful primary/secondary keys** for event correlation
4. **Test configs** before deploying to production
5. **Document custom changes** in comments (future feature)

## Updates

This configuration library is maintained by BigPanda. Check for updates:

```bash
# Check version
cat /opt/bigpanda-agent/VERSION

# Update agent
sudo apt update && sudo apt install bigpanda-agent
```

## Support

- Documentation: https://docs.bigpanda.io/super-agent
- Issues: https://github.com/ReggieJTech/SuperAgent/issues
- Email: support@bigpanda.io
