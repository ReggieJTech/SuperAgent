[
  {
    "type": "MIB",
    "trap": "ccCritical",
    "mib": "VERITAS-COMMAND-CENTRAL-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A Critical alert trap from Command Central.",
      "mib_name": "VERITAS-COMMAND-CENTRAL-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "MIB",
    "trap": "ccError",
    "mib": "VERITAS-COMMAND-CENTRAL-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "An Error alert trap from Command Central.",
      "mib_name": "VERITAS-COMMAND-CENTRAL-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "MIB",
    "trap": "ccWarning",
    "mib": "VERITAS-COMMAND-CENTRAL-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "warning",
      "description": "A Warning alert trap from Command Central.",
      "mib_name": "VERITAS-COMMAND-CENTRAL-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "MIB",
    "trap": "ccInformational",
    "mib": "VERITAS-COMMAND-CENTRAL-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "An Informational alert trap from Command Central.",
      "mib_name": "VERITAS-COMMAND-CENTRAL-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  }
]