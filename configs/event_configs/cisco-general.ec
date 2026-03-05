[
  {
    "type": "NO_MIB",
    "trap": "linkDown",
    "trap-oid": "1.3.6.1.6.3.1.1.5.3",
    "trap-var-binds": {
      "ifIndex": 1,
      "ifAdminStatus": 2,
      "ifOperStatus": 3,
      "ifDescr": 4
    },
    "copy": {},
    "rename": {
      "ifIndex": "interface_index",
      "ifDescr": "interface_name",
      "ifAdminStatus": "admin_status",
      "ifOperStatus": "oper_status"
    },
    "set": {
      "status": "critical",
      "description": "Network interface is down",
      "mib_name": "IF-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "interface_name"
  },
  {
    "type": "NO_MIB",
    "trap": "linkUp",
    "trap-oid": "1.3.6.1.6.3.1.1.5.4",
    "trap-var-binds": {
      "ifIndex": 1,
      "ifAdminStatus": 2,
      "ifOperStatus": 3,
      "ifDescr": 4
    },
    "copy": {},
    "rename": {
      "ifIndex": "interface_index",
      "ifDescr": "interface_name",
      "ifAdminStatus": "admin_status",
      "ifOperStatus": "oper_status"
    },
    "set": {
      "status": "ok",
      "description": "Network interface is up",
      "mib_name": "IF-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "interface_name"
  },
  {
    "type": "NO_MIB",
    "trap": "coldStart",
    "trap-oid": "1.3.6.1.6.3.1.1.5.1",
    "trap-var-binds": {},
    "copy": {},
    "rename": {},
    "set": {
      "status": "warning",
      "description": "Device has been restarted (cold start)",
      "mib_name": "SNMPv2-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "warmStart",
    "trap-oid": "1.3.6.1.6.3.1.1.5.2",
    "trap-var-binds": {},
    "copy": {},
    "rename": {},
    "set": {
      "status": "warning",
      "description": "Device has been restarted (warm start)",
      "mib_name": "SNMPv2-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "authenticationFailure",
    "trap-oid": "1.3.6.1.6.3.1.1.5.5",
    "trap-var-binds": {},
    "copy": {},
    "rename": {},
    "set": {
      "status": "warning",
      "description": "SNMP authentication failure",
      "mib_name": "SNMPv2-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  }
]
