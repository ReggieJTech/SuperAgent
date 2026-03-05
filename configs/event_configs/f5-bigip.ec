[
  {
    "type": "NO_MIB",
    "trap": "bigipAgentStart",
    "trap-oid": "1.3.6.1.4.1.3375.2.4.0.1",
    "trap-var-binds": {},
    "copy": {},
    "rename": {},
    "set": {
      "status": "warning",
      "description": "F5 BIG-IP agent started",
      "mib_name": "F5-BIGIP-COMMON-MIB",
      "vendor": "f5"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "bigipAgentShutdown",
    "trap-oid": "1.3.6.1.4.1.3375.2.4.0.2",
    "trap-var-binds": {},
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "F5 BIG-IP agent shutdown",
      "mib_name": "F5-BIGIP-COMMON-MIB",
      "vendor": "f5"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "bigipNodeDown",
    "trap-oid": "1.3.6.1.4.1.3375.2.4.0.10",
    "trap-var-binds": {
      "nodeName": 1,
      "nodeAddr": 2
    },
    "copy": {},
    "rename": {
      "nodeName": "node_name",
      "nodeAddr": "node_address"
    },
    "set": {
      "status": "critical",
      "description": "F5 BIG-IP node is down",
      "mib_name": "F5-BIGIP-COMMON-MIB",
      "vendor": "f5"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "node_name"
  },
  {
    "type": "NO_MIB",
    "trap": "bigipNodeUp",
    "trap-oid": "1.3.6.1.4.1.3375.2.4.0.11",
    "trap-var-binds": {
      "nodeName": 1,
      "nodeAddr": 2
    },
    "copy": {},
    "rename": {
      "nodeName": "node_name",
      "nodeAddr": "node_address"
    },
    "set": {
      "status": "ok",
      "description": "F5 BIG-IP node is up",
      "mib_name": "F5-BIGIP-COMMON-MIB",
      "vendor": "f5"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "node_name"
  }
]
