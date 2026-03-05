# Supported Vendors

The BigPanda Super Agent includes pre-bundled event configurations for the following vendors and devices.

## Summary

- **Total Event Configurations**: 18 files
- **Vendor Categories**: 6 (Network, Security, Storage, Servers, Specialized, Power)
- **Configuration Type**: Primarily NO_MIB (no MIB compilation required)

## Vendor List

### Network Infrastructure

| Vendor | File | Devices | Trap Count |
|--------|------|---------|------------|
| Cisco Systems | cisco-general.ec | IOS, IOS-XE, IOS-XR, NX-OS, ASA | 5+ |
| Juniper Networks | juniper.ec | MX, EX, QFX, SRX series | 3+ |
| Arista Networks | arista.ec | 7000, 7500 series switches | 2+ |

### Security Appliances

| Vendor | File | Devices | Trap Count |
|--------|------|---------|------------|
| F5 Networks | f5-bigip.ec | BIG-IP LTM, GTM, ASM | 4+ |
| Palo Alto Networks | palo-alto.ec | PA-Series firewalls | 3+ |
| Check Point | checkpoint.ec | Security Gateway | 2+ |
| Fortinet | fortinet.ec | FortiGate firewalls | 4+ |

### Storage Systems

| Vendor | File | Devices | Trap Count |
|--------|------|---------|------------|
| NetApp | netapp.ec | FAS, AFF systems | 4+ |
| Dell EMC | dell-emc.ec | PowerEdge servers, EMC storage | 4+ |

### Servers & Hardware

| Vendor | File | Devices | Trap Count |
|--------|------|---------|------------|
| HP/HPE | hp-hpe.ec | ProLiant, HPE servers | 3+ |
| IBM | Ibm2100-MIB.ec | System Storage 2100 | Multiple |

### Specialized Equipment

| Vendor | File | Use Case | Trap Count |
|--------|------|----------|------------|
| Arbor Networks | arbor.ec | DDoS Protection (Pravail, TMS) | Multiple |
| Interactive Intelligence | I3IC-MIB.ec | Call Center (Interaction Center) | 10+ |
| Infinera | INFINERA-TRAP-MIB.ec | Optical Networking (DTN-X, XTM) | Multiple |
| Veritas | VERITAS-*.ec | Backup & Storage Management | Multiple |

### Power & Infrastructure

| Type | File | Use Case | Trap Count |
|------|------|----------|------------|
| PDU | PDU2-MIB.ec | Power Distribution Units | Multiple |

## Trap Coverage by Category

### Common Traps (All Vendors)
- Link Up/Down
- Device Cold Start
- Device Warm Start
- Authentication Failures

### Hardware Monitoring
- Power Supply Failures
- Fan Failures
- Temperature Alerts
- Disk Failures

### Performance & Capacity
- CPU High Threshold
- Memory High Threshold
- Volume Full/Nearly Full
- Disk Space Alerts

### Network Services
- BGP Peer State Changes
- VPN Tunnel Up/Down
- Interface Status Changes
- Node Up/Down (Load Balancers)

### System Events
- Agent Start/Shutdown
- System Restarts
- Configuration Changes
- Module State Changes

## Configuration Format

All event configurations use JSON format with the following structure:

```json
{
  "type": "NO_MIB",           // or "MIB"
  "trap": "trapName",
  "trap-oid": "1.3.6.1.x.x.x",
  "trap-var-binds": {},       // Position-based varbind mapping
  "rename": {},               // Field renaming
  "set": {},                  // Static values
  "map-status": {},           // Status mapping
  "primary": "field",         // Primary key
  "secondary": "field"        // Secondary key
}
```

## Adding New Vendors

To add support for a new vendor:

1. **Obtain vendor MIB files** or trap documentation
2. **Create .ec file** in `/etc/bigpanda-agent/snmp/event_configs/`
3. **Define trap mappings** using NO_MIB format (preferred)
4. **Test with sample traps** from the vendor's device
5. **Document** in this file and submit PR (optional)

## Enterprise Support

For enterprise customers requiring:
- Custom vendor integrations
- MIB conversion services
- Trap mapping consultation
- 24/7 support

Contact: enterprise@bigpanda.io

## Updates

Check for vendor config updates:
```bash
# Check for new configs
git pull origin main configs/event_configs/

# Or update agent package
sudo apt update && sudo apt upgrade bigpanda-agent
```

## Community Contributions

Have event configurations for a vendor not listed here?

Submit a pull request:
https://github.com/ReggieJTech/SuperAgent/pulls

## Testing Configurations

Test a vendor config:

```bash
# Send test trap
snmptrap -v 2c -c public localhost:162 '' <trap-oid>

# View in logs
tail -f /var/log/bigpanda-agent/agent.log | grep snmp

# Check statistics
curl http://localhost:8443/stats | jq '.plugins.snmp'
```

## Compatibility

All configurations are tested and compatible with:
- BigPanda Super Agent v1.0.0+
- SNMPv1, SNMPv2c, SNMPv3
- BigPanda API v2

## License

These event configurations are provided under the BigPanda Super Agent license.
Commercial use is permitted for licensed BigPanda customers.
