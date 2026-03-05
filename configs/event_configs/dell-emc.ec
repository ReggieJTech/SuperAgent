[
  {
    "type": "NO_MIB",
    "trap": "alertFanFailure",
    "trap-oid": "1.3.6.1.4.1.674.10892.1.0.1052",
    "trap-var-binds": {
      "drsGlobalCurrStatus": 1,
      "drsChassisFQDD": 2
    },
    "copy": {},
    "rename": {
      "drsGlobalCurrStatus": "current_status",
      "drsChassisFQDD": "chassis_fqdd"
    },
    "set": {
      "status": "critical",
      "description": "Dell EMC fan failure detected",
      "mib_name": "DELL-RAC-MIB",
      "vendor": "dell_emc"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "chassis_fqdd"
  },
  {
    "type": "NO_MIB",
    "trap": "alertPowerSupplyFailure",
    "trap-oid": "1.3.6.1.4.1.674.10892.1.0.1053",
    "trap-var-binds": {
      "drsGlobalCurrStatus": 1,
      "drsPowerSupplyFQDD": 2
    },
    "copy": {},
    "rename": {
      "drsGlobalCurrStatus": "current_status",
      "drsPowerSupplyFQDD": "psu_fqdd"
    },
    "set": {
      "status": "critical",
      "description": "Dell EMC power supply failure",
      "mib_name": "DELL-RAC-MIB",
      "vendor": "dell_emc"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "psu_fqdd"
  },
  {
    "type": "NO_MIB",
    "trap": "alertTemperatureProbeNormal",
    "trap-oid": "1.3.6.1.4.1.674.10892.1.0.1154",
    "trap-var-binds": {
      "drsGlobalCurrStatus": 1,
      "drsTemperatureProbeFQDD": 2
    },
    "copy": {},
    "rename": {
      "drsTemperatureProbeFQDD": "probe_fqdd"
    },
    "set": {
      "status": "ok",
      "description": "Dell EMC temperature probe returned to normal",
      "mib_name": "DELL-RAC-MIB",
      "vendor": "dell_emc"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "probe_fqdd"
  },
  {
    "type": "NO_MIB",
    "trap": "alertTemperatureProbeFailure",
    "trap-oid": "1.3.6.1.4.1.674.10892.1.0.1052",
    "trap-var-binds": {
      "drsGlobalCurrStatus": 1,
      "drsTemperatureProbeFQDD": 2,
      "drsTemperatureProbeReading": 3
    },
    "copy": {},
    "rename": {
      "drsTemperatureProbeFQDD": "probe_fqdd",
      "drsTemperatureProbeReading": "temperature"
    },
    "set": {
      "status": "critical",
      "description": "Dell EMC temperature probe failure",
      "mib_name": "DELL-RAC-MIB",
      "vendor": "dell_emc"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "probe_fqdd"
  }
]
