[
  {
    "type": "NO_MIB",
    "trap": "fgTrapCpuHighThreshold",
    "trap-oid": "1.3.6.1.4.1.12356.101.6.0.101",
    "trap-var-binds": {
      "fgSysCpuUsage": 1
    },
    "copy": {},
    "rename": {
      "fgSysCpuUsage": "cpu_usage"
    },
    "set": {
      "status": "warning",
      "description": "Fortinet CPU usage high",
      "mib_name": "FORTINET-FORTIGATE-MIB",
      "vendor": "fortinet"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "fgTrapMemThreshold",
    "trap-oid": "1.3.6.1.4.1.12356.101.6.0.102",
    "trap-var-binds": {
      "fgSysMemUsage": 1
    },
    "copy": {},
    "rename": {
      "fgSysMemUsage": "memory_usage"
    },
    "set": {
      "status": "warning",
      "description": "Fortinet memory usage high",
      "mib_name": "FORTINET-FORTIGATE-MIB",
      "vendor": "fortinet"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "fgTrapVpnTunUp",
    "trap-oid": "1.3.6.1.4.1.12356.101.6.0.401",
    "trap-var-binds": {
      "fgVpnTunEntPhase1Name": 1,
      "fgVpnTunEntPhase2Name": 2
    },
    "copy": {},
    "rename": {
      "fgVpnTunEntPhase1Name": "vpn_phase1",
      "fgVpnTunEntPhase2Name": "vpn_phase2"
    },
    "set": {
      "status": "ok",
      "description": "Fortinet VPN tunnel is up",
      "mib_name": "FORTINET-FORTIGATE-MIB",
      "vendor": "fortinet"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "vpn_phase1"
  },
  {
    "type": "NO_MIB",
    "trap": "fgTrapVpnTunDown",
    "trap-oid": "1.3.6.1.4.1.12356.101.6.0.402",
    "trap-var-binds": {
      "fgVpnTunEntPhase1Name": 1,
      "fgVpnTunEntPhase2Name": 2
    },
    "copy": {},
    "rename": {
      "fgVpnTunEntPhase1Name": "vpn_phase1",
      "fgVpnTunEntPhase2Name": "vpn_phase2"
    },
    "set": {
      "status": "critical",
      "description": "Fortinet VPN tunnel is down",
      "mib_name": "FORTINET-FORTIGATE-MIB",
      "vendor": "fortinet"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "vpn_phase1"
  }
]
