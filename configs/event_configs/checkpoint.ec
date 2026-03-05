[
  {
    "type": "NO_MIB",
    "trap": "cpvHighThreshold",
    "trap-oid": "1.3.6.1.4.1.2620.1.6.7.1001.0.1",
    "trap-var-binds": {
      "cpvHighThresholdName": 1,
      "cpvHighThresholdValue": 2
    },
    "copy": {},
    "rename": {
      "cpvHighThresholdName": "threshold_name",
      "cpvHighThresholdValue": "threshold_value"
    },
    "set": {
      "status": "warning",
      "description": "Check Point high threshold exceeded",
      "mib_name": "CHECKPOINT-MIB",
      "vendor": "checkpoint"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "threshold_name"
  },
  {
    "type": "NO_MIB",
    "trap": "cpvFwModuleState",
    "trap-oid": "1.3.6.1.4.1.2620.1.1.1.0.102",
    "trap-var-binds": {
      "fwModuleState": 1
    },
    "copy": {},
    "rename": {
      "fwModuleState": "module_state"
    },
    "set": {
      "description": "Check Point firewall module state changed",
      "mib_name": "CHECKPOINT-MIB",
      "vendor": "checkpoint"
    },
    "map-status": {
      "module_state": {
        "1": "ok",
        "0": "critical"
      }
    },
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  }
]
