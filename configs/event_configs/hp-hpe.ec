[
  {
    "type": "NO_MIB",
    "trap": "cpqDa3AccelStatusChange",
    "trap-oid": "1.3.6.1.4.1.232.3.2.3.1.0.3",
    "trap-var-binds": {
      "cpqDa3AccelSerialNumber": 1,
      "cpqDa3AccelStatus": 2
    },
    "copy": {},
    "rename": {
      "cpqDa3AccelSerialNumber": "accel_serial",
      "cpqDa3AccelStatus": "accel_status"
    },
    "set": {
      "description": "HP/HPE array accelerator status changed",
      "mib_name": "CPQIDA-MIB",
      "vendor": "hp"
    },
    "map-status": {
      "accel_status": {
        "2": "ok",
        "3": "critical",
        "4": "critical"
      }
    },
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "accel_serial"
  },
  {
    "type": "NO_MIB",
    "trap": "cpqHeTemperatureFailed",
    "trap-oid": "1.3.6.1.4.1.232.0.6011",
    "trap-var-binds": {
      "cpqHeThermalSystemFanStatus": 1,
      "cpqHeThermalTempStatus": 2
    },
    "copy": {},
    "rename": {
      "cpqHeThermalSystemFanStatus": "fan_status",
      "cpqHeThermalTempStatus": "temp_status"
    },
    "set": {
      "status": "critical",
      "description": "HP/HPE temperature sensor failed",
      "mib_name": "CPQHOST-MIB",
      "vendor": "hp"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "cpqHeFltTolPowerSupplyFailed",
    "trap-oid": "1.3.6.1.4.1.232.0.6004",
    "trap-var-binds": {
      "cpqHeFltTolPowerSupplyStatus": 1,
      "cpqHeFltTolPowerSupplyBay": 2
    },
    "copy": {},
    "rename": {
      "cpqHeFltTolPowerSupplyStatus": "psu_status",
      "cpqHeFltTolPowerSupplyBay": "psu_bay"
    },
    "set": {
      "status": "critical",
      "description": "HP/HPE power supply failed",
      "mib_name": "CPQHOST-MIB",
      "vendor": "hp"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "psu_bay"
  }
]
