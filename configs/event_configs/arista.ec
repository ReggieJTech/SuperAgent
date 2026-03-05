[
  {
    "type": "NO_MIB",
    "trap": "aristaFanStatusChanged",
    "trap-oid": "1.3.6.1.4.1.30065.3.6.1.2.1",
    "trap-var-binds": {
      "aristaFanName": 1,
      "aristaFanState": 2
    },
    "copy": {},
    "rename": {
      "aristaFanName": "fan_name",
      "aristaFanState": "fan_state"
    },
    "set": {
      "description": "Arista fan status changed",
      "mib_name": "ARISTA-ENTITY-SENSOR-MIB",
      "vendor": "arista"
    },
    "map-status": {
      "fan_state": {
        "1": "ok",
        "2": "warning",
        "3": "critical"
      }
    },
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "fan_name"
  },
  {
    "type": "NO_MIB",
    "trap": "aristaPowerSupplyStatusChanged",
    "trap-oid": "1.3.6.1.4.1.30065.3.6.1.3.1",
    "trap-var-binds": {
      "aristaPowerSupplyName": 1,
      "aristaPowerSupplyState": 2
    },
    "copy": {},
    "rename": {
      "aristaPowerSupplyName": "psu_name",
      "aristaPowerSupplyState": "psu_state"
    },
    "set": {
      "description": "Arista power supply status changed",
      "mib_name": "ARISTA-ENTITY-SENSOR-MIB",
      "vendor": "arista"
    },
    "map-status": {
      "psu_state": {
        "1": "ok",
        "2": "warning",
        "3": "critical"
      }
    },
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "psu_name"
  }
]
