[
  {
    "type": "MIB",
    "trap": "i3IcTrapRestart",
    "mib": "I3IC-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "This trap is generated when an Interaction Center subsystem has been restarted. See accompanying MIB variable.",
      "mib_name": "I3IC-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "MIB",
    "trap": "i3IcTrapInformationalEventLog",
    "mib": "I3IC-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "This trap is generated when Interaction Center has written an informational entry to the Event Log on the Interaction Center server See accompanying MIB variable.",
      "mib_name": "I3IC-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "MIB",
    "trap": "i3IcTrapWarningEventLog",
    "mib": "I3IC-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "warning",
      "description": "This trap is generated when Interaction Center has written a warning entry to the Event Log on the Interaction Center server See accompanying MIB variable.",
      "mib_name": "I3IC-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "MIB",
    "trap": "i3IcTrapErrorEventLog",
    "mib": "I3IC-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "This trap is generated when Interaction Center has written an error entry to the Event Log on the Interaction Center server See accompanying MIB variable.",
      "mib_name": "I3IC-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "MIB",
    "trap": "i3IcCommenceSwitchoverEvent",
    "mib": "I3IC-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "This trap is generated when the backup Interaction Center server is becoming active.",
      "mib_name": "I3IC-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "MIB",
    "trap": "i3IcTSInterfaceYellowAlarm",
    "mib": "I3IC-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "This trap is generated when a telephony interface has a yellow alarm raised on it.",
      "mib_name": "I3IC-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "MIB",
    "trap": "i3IcTSInterfaceRedAlarm",
    "mib": "I3IC-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "This trap is generated when a telephony interface has a red alarm raised on it.",
      "mib_name": "I3IC-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "MIB",
    "trap": "i3IcTSInterfaceAlarmCleared",
    "mib": "I3IC-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "This trap is generated when a telephony interface has an alarm cleared.",
      "mib_name": "I3IC-MIB",
      "base_trap_name": "i3IcTSInterfaceAlarmed"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "MIB",
    "trap": "i3IcTSInterfaceDChannelDown",
    "mib": "I3IC-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "This trap is generated when a telephony interface has a D channel go down.",
      "mib_name": "I3IC-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "MIB",
    "trap": "i3IcDiscontinueSwitchoverEvent",
    "mib": "I3IC-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "This trap is generated when the primary Interaction Center server is becoming inactive.",
      "mib_name": "I3IC-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "MIB",
    "trap": "i3IcHighLatencyState",
    "mib": "I3IC-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "This trap is generated when Interaction Center server is experiencing a high latency condition.",
      "mib_name": "I3IC-MIB",
      "base_trap_name": "i3IcHighLatencyState"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "MIB",
    "trap": "i3IcHighLatencyStateClear",
    "mib": "I3IC-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "This trap is generated when Interaction Center server is no longer experiencing a high latency condition.",
      "mib_name": "I3IC-MIB",
      "base_trap_name": "i3IcHighLatencyState"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "MIB",
    "trap": "i3IcTracingStoppedInsufficientDiskSpace",
    "mib": "I3IC-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "This trap is generated when tracing was stopped due to insufficient disk space.",
      "mib_name": "I3IC-MIB"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "MIB",
    "trap": "i3IcTracingStoppedInsufficientDiskSpaceCleared",
    "mib": "I3IC-MIB",
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "This trap is generated when tracing is resumed.",
      "mib_name": "I3IC-MIB",
      "base_trap_name": "i3IcTracingStoppedInsufficientDiskSpaceed"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  }
]