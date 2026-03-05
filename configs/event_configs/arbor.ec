[
  {
    "type": "NO_MIB",
    "trap": "subHostDown",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.7",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapSubhostName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Generated when a subhost transitions to inactive",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedProtectionGroupError",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.12",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedProtectionGroupId": 4,
      "aedProtectionGroupName": 5
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A protection group has an error",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedConfigError",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.16",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Generated when an internal configuration error is detected.",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedHwDeviceDown",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.18",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A hardware device has failed.",
      "mib_name": "AED-MIB",
      "base_trap_name": "aedHwDeviceDown"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedHwDeviceDownDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.19",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "A hardware device failure is no longer detected.",
      "mib_name": "AED-MIB",
      "base_trap_name": "aedHwDeviceDown"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedHwSensorCritical",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.20",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A hardware sensor is reading an alarm condition.",
      "mib_name": "AED-MIB",
      "base_trap_name": "aedHwSensorCritical"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedHwSensorCriticalDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.21",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "A hardware sensor is no longer reading an alarm condition.",
      "mib_name": "AED-MIB",
      "base_trap_name": "aedHwSensorCritical"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedSwComponentDown",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.22",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapSubhostName": 4,
      "aedTrapComponentName": 5
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A software program has failed.",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedSystemStatusCritical",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.24",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The AED system is experiencing a critical failure.",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedSystemStatusDegraded",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.25",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The AED system is experiencing degraded performance.",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedFilesystemCritical",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.27",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A filesystem is near capacity.",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedGRETunnelFailure",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.28",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A GRE tunnel failed.",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedNextHopUnreachable",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.29",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The next hop system is unreachable from the AED system.",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedPerformance",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.30",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The AED system is dropping traffic because of high traffic rates.",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedSystemStatusError",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.32",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The AED system is experiencing an error.",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedCloudSignalTimeout",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.46",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "There is a timeout communicating with cloud services",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedCloudSignalThreshold",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.47",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "There is a threshold error with cloud services",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedDeploymentModeChange",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.48",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedDeploymentMode": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The AED deployment mode changed",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedProtectionLevelChange",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.49",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedPreviousProtectionLevel": 4,
      "aedProtectionLevel": 5
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The AED protection level changed",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedTrapBlockHostDetail",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.50",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedBlockSrcAddr": 4,
      "aedBlockProtocol": 5,
      "aedBlockDstAddr": 6,
      "aedBlockSrcPort": 7,
      "aedBlockDstPort": 8,
      "aedTrapComponentName": 9,
      "aedInetAddressType": 10,
      "aedBlockSrcInetAddr": 11,
      "aedBlockDstInetAddr": 12
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A host was blocked by AED",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedTrapBlockHostSummary",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.51",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Number of additional hosts blocked by AED",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedTrapTraffic",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.52",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrafficLevel": 4,
      "aedTrafficUnits": 5,
      "aedProtectionGroupName": 6,
      "aedURL": 7
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Traffic exceeded the threshold for a AED protection group",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedTrapBotnetAttack",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.53",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrafficLevel": 4,
      "aedTrafficUnits": 5,
      "aedProtectionGroupName": 6,
      "aedURL": 7
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Botnet traffic which was detected but not blocked exceed the threshold for a AED protection group",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedTrapLicenseLimit",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.54",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedURL": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The total bandwidth of traffic in the AED system is approaching the license limit",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedTrapBlockedTraffic",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.55",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrafficLevel": 4,
      "aedTrafficUnits": 5,
      "aedProtectionGroupName": 6,
      "aedURL": 7
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Blocked traffic exceeded the threshold for a AED protection group",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedProtectionGroupLevelChange",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.56",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedPreviousProtectionLevel": 4,
      "aedProtectionLevel": 5,
      "aedProtectionGroupName": 6
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The AED protection level changed for one protection group.",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedProtectionGroupModeChange",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.57",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedDeploymentMode": 4,
      "aedProtectionGroupName": 5
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The AED protection mode changed for a protection group",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedChangeLogEntry",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.61",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedUsername": 4,
      "aedSubsystem": 5,
      "aedSettingType": 6
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A new entry has been added to the AED change log",
      "mib_name": "AED-MIB"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedAutomationActivation",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.62",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A Protection Category Automation is now active",
      "mib_name": "AED-MIB",
      "base_trap_name": "aedAutomationActivation"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedAutomationActivationDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.63",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "A Protection Category Automation is now inactive",
      "mib_name": "AED-MIB",
      "base_trap_name": "aedAutomationActivation"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedHwDeviceWarning",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.64",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "warning",
      "description": "A hardware device has triggered a warning condition",
      "mib_name": "AED-MIB",
      "base_trap_name": "aedHwDeviceWarning"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedHwDeviceWarningDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.65",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "warning",
      "description": "A hardware device warning is no longer detected.",
      "mib_name": "AED-MIB",
      "base_trap_name": "aedHwDeviceWarning"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedHwDeviceAlert",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.66",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A hardware device has triggered an alert",
      "mib_name": "AED-MIB",
      "base_trap_name": "aedHwDeviceAlert"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "aedHwDeviceAlertDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.13.3.0.67",
    "trap-var-binds": {
      "sysName": 1,
      "aedTrapString": 2,
      "aedTrapDetail": 3,
      "aedTrapComponentName": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "A hardware device alert is no longer detected.",
      "mib_name": "AED-MIB",
      "base_trap_name": "aedHwDeviceAlert"
    },
    "map-status": {},
    "primary": "sysName",
    "secondary": "aedTrapString"
  },
  {
    "type": "NO_MIB",
    "trap": "bandwidthAnomaly",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.1",
    "trap-var-binds": {
      "pdosAnomalyId": 1,
      "pdosAnomalyDirection": 2,
      "pdosAnomalyResource": 3,
      "pdosAnomalyLinkPercent": 4,
      "pdosAnomalyClassification": 5,
      "pdosAnomalyStart": 6,
      "pdosAnomalyDuration": 7,
      "pdosAnomalyRouterInterfaces": 8,
      "pdosUrl": 9
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Bandwidth anomaly detected by Peakflow",
      "mib_name": "PEAKFLOW-DOS-MIB"
    },
    "map-status": {},
    "primary": "pdosAnomalyId",
    "secondary": "pdosAnomalyResource"
  },
  {
    "type": "NO_MIB",
    "trap": "tcpflagsAnomaly",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.2",
    "trap-var-binds": {
      "pdosAnomalyId": 1,
      "pdosAnomalyDirection": 2,
      "pdosAnomalyResource": 3,
      "pdosAnomalyLinkPercent": 4,
      "pdosAnomalyClassification": 5,
      "pdosAnomalyStart": 6,
      "pdosAnomalyDuration": 7,
      "pdosAnomalyRouterInterfaces": 8,
      "pdosUrl": 9,
      "pdosAnomalyTcpFlags": 10
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "TCP flags anomaly detected by Peakflow",
      "mib_name": "PEAKFLOW-DOS-MIB"
    },
    "map-status": {},
    "primary": "pdosAnomalyId",
    "secondary": "pdosAnomalyResource"
  },
  {
    "type": "NO_MIB",
    "trap": "protocolAnomaly",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.3",
    "trap-var-binds": {
      "pdosAnomalyId": 1,
      "pdosAnomalyDirection": 2,
      "pdosAnomalyResource": 3,
      "pdosAnomalyLinkPercent": 4,
      "pdosAnomalyClassification": 5,
      "pdosAnomalyStart": 6,
      "pdosAnomalyDuration": 7,
      "pdosAnomalyRouterInterfaces": 8,
      "pdosUrl": 9,
      "pdosAnomalyProto": 10
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Protocol anomaly detected by Peakflow",
      "mib_name": "PEAKFLOW-DOS-MIB"
    },
    "map-status": {},
    "primary": "pdosAnomalyId",
    "secondary": "pdosAnomalyResource"
  },
  {
    "type": "NO_MIB",
    "trap": "heartbeatLoss",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.4",
    "trap-var-binds": {
      "pdosHeartbeatSource": 1
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Missing heartbeat from SP device to leader",
      "mib_name": "PEAKFLOW-DOS-MIB",
      "base_trap_name": "heartbeatLoss"
    },
    "map-status": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "internalError",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.5",
    "trap-var-binds": {
      "internalErrorLocation": 1,
      "internalErrorReason": 2
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Internal inconsistency or error",
      "mib_name": "PEAKFLOW-DOS-MIB"
    },
    "map-status": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "anomalyDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.6",
    "trap-var-binds": {
      "pdosAnomalyId": 1,
      "pdosAnomalyDirection": 2,
      "pdosAnomalyResource": 3,
      "pdosAnomalyLinkPercent": 4,
      "pdosAnomalyClassification": 5,
      "pdosAnomalyStart": 6,
      "pdosAnomalyDuration": 7,
      "pdosAnomalyRouterInterfaces": 8,
      "pdosUrl": 9
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "Some previously detected anomaly is no longer active",
      "mib_name": "PEAKFLOW-DOS-MIB",
      "base_trap_name": "anomaly"
    },
    "map-status": {},
    "primary": "pdosAnomalyId",
    "secondary": "pdosAnomalyResource"
  },
  {
    "type": "NO_MIB",
    "trap": "netflowMissing",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.7",
    "trap-var-binds": {
      "pdosAnomalyRouter": 1
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "NetFlow has not been received from a NetFlow transmitting router",
      "mib_name": "PEAKFLOW-DOS-MIB",
      "base_trap_name": "netflowMissing"
    },
    "map-status": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "netflowMissingDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.8",
    "trap-var-binds": {
      "pdosAnomalyRouter": 1
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "NetFlow has resumed from a router which previously was not forwarding NetFlow data",
      "mib_name": "PEAKFLOW-DOS-MIB",
      "base_trap_name": "netflowMissing"
    },
    "map-status": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "icmpMisuseAnomaly",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.9",
    "trap-var-binds": {
      "pdosAnomalyId": 1,
      "pdosAnomalyDirection": 2,
      "pdosAnomalyResource": 3,
      "pdosAnomalyLinkPercent": 4,
      "pdosAnomalyClassification": 5,
      "pdosAnomalyStart": 6,
      "pdosAnomalyDuration": 7,
      "pdosAnomalyRouterInterfaces": 8,
      "pdosUrl": 9,
      "pdosAnomalyProto": 10
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "ICMP misuse anomaly detected by Peakflow",
      "mib_name": "PEAKFLOW-DOS-MIB"
    },
    "map-status": {},
    "primary": "pdosAnomalyId",
    "secondary": "pdosAnomalyResource"
  },
  {
    "type": "NO_MIB",
    "trap": "tcpNullMisuseAnomaly",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.10",
    "trap-var-binds": {
      "pdosAnomalyId": 1,
      "pdosAnomalyDirection": 2,
      "pdosAnomalyResource": 3,
      "pdosAnomalyLinkPercent": 4,
      "pdosAnomalyClassification": 5,
      "pdosAnomalyStart": 6,
      "pdosAnomalyDuration": 7,
      "pdosAnomalyRouterInterfaces": 8,
      "pdosUrl": 9,
      "pdosAnomalyProto": 10,
      "pdosAnomalyTcpFlags": 11
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "TCP Null misuse anomaly detected by Peakflow",
      "mib_name": "PEAKFLOW-DOS-MIB"
    },
    "map-status": {},
    "primary": "pdosAnomalyId",
    "secondary": "pdosAnomalyResource"
  },
  {
    "type": "NO_MIB",
    "trap": "tcpSynMisuseAnomaly",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.11",
    "trap-var-binds": {
      "pdosAnomalyId": 1,
      "pdosAnomalyDirection": 2,
      "pdosAnomalyResource": 3,
      "pdosAnomalyLinkPercent": 4,
      "pdosAnomalyClassification": 5,
      "pdosAnomalyStart": 6,
      "pdosAnomalyDuration": 7,
      "pdosAnomalyRouterInterfaces": 8,
      "pdosUrl": 9,
      "pdosAnomalyProto": 10,
      "pdosAnomalyTcpFlags": 11
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "TCP SYN misuse anomaly detected by Peakflow",
      "mib_name": "PEAKFLOW-DOS-MIB"
    },
    "map-status": {},
    "primary": "pdosAnomalyId",
    "secondary": "pdosAnomalyResource"
  },
  {
    "type": "NO_MIB",
    "trap": "ipNullMisuseAnomaly",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.12",
    "trap-var-binds": {
      "pdosAnomalyId": 1,
      "pdosAnomalyDirection": 2,
      "pdosAnomalyResource": 3,
      "pdosAnomalyLinkPercent": 4,
      "pdosAnomalyClassification": 5,
      "pdosAnomalyStart": 6,
      "pdosAnomalyDuration": 7,
      "pdosAnomalyRouterInterfaces": 8,
      "pdosUrl": 9
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "IP Null misuse anomaly detected by Peakflow",
      "mib_name": "PEAKFLOW-DOS-MIB"
    },
    "map-status": {},
    "primary": "pdosAnomalyId",
    "secondary": "pdosAnomalyResource"
  },
  {
    "type": "NO_MIB",
    "trap": "ipFragMisuseAnomaly",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.13",
    "trap-var-binds": {
      "pdosAnomalyId": 1,
      "pdosAnomalyDirection": 2,
      "pdosAnomalyResource": 3,
      "pdosAnomalyLinkPercent": 4,
      "pdosAnomalyClassification": 5,
      "pdosAnomalyStart": 6,
      "pdosAnomalyDuration": 7,
      "pdosAnomalyRouterInterfaces": 8,
      "pdosUrl": 9,
      "pdosAnomalyProto": 10
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "IP Fragment misuse anomaly detected by Peakflow",
      "mib_name": "PEAKFLOW-DOS-MIB"
    },
    "map-status": {},
    "primary": "pdosAnomalyId",
    "secondary": "pdosAnomalyResource"
  },
  {
    "type": "NO_MIB",
    "trap": "ipPrivateMisuseAnomaly",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.14",
    "trap-var-binds": {
      "pdosAnomalyId": 1,
      "pdosAnomalyDirection": 2,
      "pdosAnomalyResource": 3,
      "pdosAnomalyLinkPercent": 4,
      "pdosAnomalyClassification": 5,
      "pdosAnomalyStart": 6,
      "pdosAnomalyDuration": 7,
      "pdosAnomalyRouterInterfaces": 8,
      "pdosUrl": 9
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "IP Private misuse anomaly detected by Peakflow",
      "mib_name": "PEAKFLOW-DOS-MIB"
    },
    "map-status": {},
    "primary": "pdosAnomalyId",
    "secondary": "pdosAnomalyResource"
  },
  {
    "type": "NO_MIB",
    "trap": "heartbeatLossDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.15",
    "trap-var-binds": {
      "pdosHeartbeatSource": 1
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "Heartbeat from SP device to leader now works",
      "mib_name": "PEAKFLOW-DOS-MIB",
      "base_trap_name": "heartbeatLoss"
    },
    "map-status": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "tcpRstMisuseAnomaly",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.16",
    "trap-var-binds": {
      "pdosAnomalyId": 1,
      "pdosAnomalyDirection": 2,
      "pdosAnomalyResource": 3,
      "pdosAnomalyLinkPercent": 4,
      "pdosAnomalyClassification": 5,
      "pdosAnomalyStart": 6,
      "pdosAnomalyDuration": 7,
      "pdosAnomalyRouterInterfaces": 8,
      "pdosUrl": 9,
      "pdosAnomalyProto": 10,
      "pdosAnomalyTcpFlags": 11
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "TCP RST misuse anomaly detected by Peakflow",
      "mib_name": "PEAKFLOW-DOS-MIB"
    },
    "map-status": {},
    "primary": "pdosAnomalyId",
    "secondary": "pdosAnomalyResource"
  },
  {
    "type": "NO_MIB",
    "trap": "totalTrafficMisuseAnomaly",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.17",
    "trap-var-binds": {
      "pdosAnomalyId": 1,
      "pdosAnomalyDirection": 2,
      "pdosAnomalyResource": 3,
      "pdosAnomalyLinkPercent": 4,
      "pdosAnomalyClassification": 5,
      "pdosAnomalyStart": 6,
      "pdosAnomalyDuration": 7,
      "pdosAnomalyRouterInterfaces": 8,
      "pdosUrl": 9
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Total Traffic misuse anomaly detected by Peakflow",
      "mib_name": "PEAKFLOW-DOS-MIB"
    },
    "map-status": {},
    "primary": "pdosAnomalyId",
    "secondary": "pdosAnomalyResource"
  },
  {
    "type": "NO_MIB",
    "trap": "fingerprintAnomaly",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.18",
    "trap-var-binds": {
      "pdosAnomalyId": 1,
      "pdosAnomalyDirection": 2,
      "pdosAnomalyResource": 3,
      "pdosAnomalyLinkPercent": 4,
      "pdosAnomalyClassification": 5,
      "pdosAnomalyStart": 6,
      "pdosAnomalyDuration": 7,
      "pdosAnomalyRouterInterfaces": 8,
      "pdosUrl": 9
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Fingerprint anomaly detected by Peakflow",
      "mib_name": "PEAKFLOW-DOS-MIB"
    },
    "map-status": {},
    "primary": "pdosAnomalyId",
    "secondary": "pdosAnomalyResource"
  },
  {
    "type": "NO_MIB",
    "trap": "dnsMisuseAnomaly",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.19",
    "trap-var-binds": {
      "pdosAnomalyId": 1,
      "pdosAnomalyDirection": 2,
      "pdosAnomalyResource": 3,
      "pdosAnomalyLinkPercent": 4,
      "pdosAnomalyClassification": 5,
      "pdosAnomalyStart": 6,
      "pdosAnomalyDuration": 7,
      "pdosAnomalyRouterInterfaces": 8,
      "pdosUrl": 9,
      "pdosAnomalyProto": 10
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "DNS misuse anomaly detected by Peakflow",
      "mib_name": "PEAKFLOW-DOS-MIB"
    },
    "map-status": {},
    "primary": "pdosAnomalyId",
    "secondary": "pdosAnomalyResource"
  },
  {
    "type": "NO_MIB",
    "trap": "udpMisuseAnomaly",
    "trap-oid": "1.3.6.1.4.1.9694.1.1.3.0.20",
    "trap-var-binds": {
      "pdosAnomalyId": 1,
      "pdosAnomalyDirection": 2,
      "pdosAnomalyResource": 3,
      "pdosAnomalyLinkPercent": 4,
      "pdosAnomalyClassification": 5,
      "pdosAnomalyStart": 6,
      "pdosAnomalyDuration": 7,
      "pdosAnomalyRouterInterfaces": 8,
      "pdosUrl": 9,
      "pdosAnomalyProto": 10
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "UDP misuse anomaly detected by Peakflow",
      "mib_name": "PEAKFLOW-DOS-MIB"
    },
    "map-status": {},
    "primary": "pdosAnomalyId",
    "secondary": "pdosAnomalyResource"
  },
  {
    "type": "NO_MIB",
    "trap": "flowDown",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.1",
    "trap-var-binds": {
      "spAlertID": 1,
      "spRouter": 2
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Flow data has not been received from a Flow transmitting router",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "flow"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "flowUp",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.2",
    "trap-var-binds": {
      "spAlertID": 1,
      "spRouter": 2
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "Flow data has resumed from a router which previously was not forwarding Flow data",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "flow"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "snmpDown",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.3",
    "trap-var-binds": {
      "spAlertID": 1,
      "spRouter": 2
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "SNMP requests are not being answered by the router",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "snmp"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "snmpUp",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.4",
    "trap-var-binds": {
      "spAlertID": 1,
      "spRouter": 2
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "SNMP requests are again being answered by the router",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "snmp"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "bgpDown",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.5",
    "trap-var-binds": {
      "spAlertID": 1,
      "spRouter": 2
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The BGP session with the router has transitioned down",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "bgp"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "bgpUp",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.6",
    "trap-var-binds": {
      "spAlertID": 1,
      "spRouter": 2
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "The BGP session with the router has transitioned up",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "bgp"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "collectorDown",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.7",
    "trap-var-binds": {
      "spAlertID": 1,
      "spCollector": 2,
      "spDetector": 3
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The SP device is down. Heartbeats are missing.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "collector"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "collectorUp",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.8",
    "trap-var-binds": {
      "spAlertID": 1,
      "spCollector": 2
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "The SP device is up. Heartbeats have resumed.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "collector"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "collectorStart",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.9",
    "trap-var-binds": {
      "spAlertID": 1,
      "spCollector": 2
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The SP device is started. Heartbeats have been received.",
      "mib_name": "PEAKFLOW-SP-MIB"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "bgpInstability",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.10",
    "trap-var-binds": {
      "spAlertID": 1,
      "spRouter": 2,
      "spBGPInstability": 3
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The BGP session with this router is exhibiting instability",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "bgpInstability"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "bgpInstabilityDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.11",
    "trap-var-binds": {
      "spAlertID": 1,
      "spRouter": 2
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "The BGP instability associated with this router has ended",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "bgpInstability"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "bgpTrap",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.12",
    "trap-var-binds": {
      "spAlertID": 1,
      "spBGPTrapName": 2,
      "spBGPTrapEvent": 3,
      "spBGPTrapPrefix": 4,
      "spBGPTrapOldAttributes": 5,
      "spBGPTrapNewAttributes": 6
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A BGP event matching this trap definition has occurred.",
      "mib_name": "PEAKFLOW-SP-MIB"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "interfaceUsage",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.13",
    "trap-var-binds": {
      "spAlertID": 1,
      "spRouter": 2,
      "spInterface": 3,
      "spInterfaceIndex": 4,
      "spInterfaceSpeed": 5,
      "spUsageType": 6,
      "spInterfaceUsage": 7
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The interface exceeded the configured traffic rate threshold.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "interfaceUsage"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "interfaceUsageDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.14",
    "trap-var-binds": {
      "spAlertID": 1,
      "spRouter": 2,
      "spInterface": 3,
      "spInterfaceIndex": 4,
      "spUsageType": 5
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "The interface is now within configured traffic thresholds",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "interfaceUsage"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "autoclassifyStart",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.15",
    "trap-var-binds": {
      "spAlertID": 1,
      "spCollector": 2,
      "spUsername": 3
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Autoclassification started on this Peakflow SP leader.",
      "mib_name": "PEAKFLOW-SP-MIB"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "configChange",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.16",
    "trap-var-binds": {
      "spAlertID": 1,
      "spCollector": 2,
      "spUsername": 3,
      "spVersion": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Configuration updated from this Peakflow SP leader.",
      "mib_name": "PEAKFLOW-SP-MIB"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "notificationLimit",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.17",
    "trap-var-binds": {
      "spAlertID": 1,
      "spCollector": 2
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Too many alerts have been generated by this Peakflow SP leader. Alerts will be temporarily suppressed. For more information about alerts that are being generated, please go to the Alerts page in the leader's UI.",
      "mib_name": "PEAKFLOW-SP-MIB"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "reportDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.18",
    "trap-var-binds": {
      "spReportName": 1,
      "spReportID": 2,
      "spReportURL": 3
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "The indicated report is finished and available for viewing at the listed URL.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "report"
    },
    "map-status": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "bgpHijack",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.19",
    "trap-var-binds": {
      "spAlertID": 1,
      "spRouter": 2,
      "spHijackRoute": 3,
      "spHijackAttr": 4,
      "spHijackLocal": 5
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A BGP announcement was seen for a prefix that is part of the configured local address space.",
      "mib_name": "PEAKFLOW-SP-MIB"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "managedObjectUsage",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.20",
    "trap-var-binds": {
      "spAlertID": 1,
      "spManagedObject": 2,
      "spManagedObjectFamily": 3,
      "spUsageType": 4,
      "spThreshold": 5,
      "spUsage": 6,
      "spUnit": 7
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The managed object exceeded the configured traffic rate threshold.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "managedObjectUsage"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "managedObjectUsageDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.21",
    "trap-var-binds": {
      "spAlertID": 1,
      "spManagedObject": 2,
      "spManagedObjectFamily": 3,
      "spUsageType": 4,
      "spUnit": 5
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "The manged object is no longer exceeding the traffic rate threshold.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "managedObjectUsage"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "hardwareFailure",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.22",
    "trap-var-binds": {
      "spAlertID": 1,
      "spCollector": 2,
      "spHardwareFailureDescription": 3
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A hardware failure has been detected on an SP device.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "hardwareFailure"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "hardwareFailureDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.23",
    "trap-var-binds": {
      "spAlertID": 1,
      "spCollector": 2,
      "spHardwareFailureDescription": 3
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "A hardware failure is no longer detected on an SP device.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "hardwareFailure"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "fingerprintFeedback",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.24",
    "trap-var-binds": {
      "spMitigationID": 1,
      "spFingerprintName": 2,
      "spFingerprintFeedback": 3
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Feedback received regarding a shared fingerprint.",
      "mib_name": "PEAKFLOW-SP-MIB"
    },
    "map-status": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "fingerprintReceive",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.25",
    "trap-var-binds": {
      "spMitigationID": 1,
      "spFingerprintName": 2,
      "spFingerprintSender": 3,
      "spFingerprintMessage": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A shared fingerprint was received.",
      "mib_name": "PEAKFLOW-SP-MIB"
    },
    "map-status": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "dnsBaseline",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.26",
    "trap-var-binds": {
      "spAlertID": 1,
      "spCollector": 2,
      "spDNSName": 3,
      "spDNSExpected": 4,
      "spDNSObserved": 5
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "An excessive number of queries for a domain name detected.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "dnsBaseline"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "dnsBaselineDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.27",
    "trap-var-binds": {
      "spAlertID": 1,
      "spCollector": 2,
      "spDNSName": 3,
      "spDNSExpected": 4,
      "spDNSObservedMean": 5,
      "spDNSObservedMax": 6
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "Query count for domain name has returned to normal levels.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "dnsBaseline"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "alertScript",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.28",
    "trap-var-binds": {
      "spMitigationID": 1,
      "spCollector": 2,
      "spMitigationName": 3,
      "spScriptCommand": 4,
      "spScriptHost": 5,
      "spScriptPort": 6,
      "spScriptStart": 7
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Alert script has been executed",
      "mib_name": "PEAKFLOW-SP-MIB"
    },
    "map-status": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "mitigationDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.29",
    "trap-var-binds": {
      "spMitigationID": 1,
      "spAlertID": 2,
      "spCollector": 3
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "Mitigation has completed",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "mitigation"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "mitigationTMSStart",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.30",
    "trap-var-binds": {
      "spMitigationID": 1,
      "spMitigationName": 2,
      "spAlertID": 3,
      "spManagedObject": 4,
      "spTMSPrefix": 5,
      "spTMSCommunity": 6,
      "spTMSTimeout": 7,
      "spMitigationStart": 8,
      "spTMSMultiPrefix": 9
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "TMS Mitigation has started",
      "mib_name": "PEAKFLOW-SP-MIB"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "mitigationThirdPartyStart",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.31",
    "trap-var-binds": {
      "spMitigationID": 1,
      "spMitigationName": 2,
      "spAlertID": 3,
      "spManagedObject": 4,
      "spThirdPartyZone": 5,
      "spThirdPartyAddr": 6,
      "spMitigationStart": 7
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Third Party Mitigation has started",
      "mib_name": "PEAKFLOW-SP-MIB"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "mitigationBlackholeStart",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.32",
    "trap-var-binds": {
      "spMitigationID": 1,
      "spMitigationName": 2,
      "spAlertID": 3,
      "spBlackholeCommunity": 4,
      "spBlackholeTimeout": 5,
      "spBlackholePrefix": 6,
      "spBlackholeNexthop": 7,
      "spMitigationStart": 8
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Blackhole Mitigation has started",
      "mib_name": "PEAKFLOW-SP-MIB"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "mitigationFlowspecStart",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.33",
    "trap-var-binds": {
      "spMitigationID": 1,
      "spMitigationName": 2,
      "spAlertID": 3,
      "spFlowspecCommunity": 4,
      "spFlowspecTimeout": 5,
      "spMitigationStart": 6
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Flowspec Mitigation has started",
      "mib_name": "PEAKFLOW-SP-MIB"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "spcommFailure",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.34",
    "trap-var-binds": {
      "spAlertID": 1,
      "spCollector": 2,
      "spCommFailureDestination": 3,
      "spCommFailureDescription": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "An SP internal communication failure has occurred.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "spcommFailure"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "spcommFailureDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.35",
    "trap-var-binds": {
      "spAlertID": 1,
      "spCollector": 2,
      "spCommFailureDestination": 3,
      "spCommFailureDescription": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "An SP internal communication failure has ended.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "spcommFailure"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "greDown",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.36",
    "trap-var-binds": {
      "spAlertID": 1,
      "spGreTunnelDestination": 2
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The GRE tunnel is down",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "greDown"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "greDownDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.37",
    "trap-var-binds": {
      "spAlertID": 1,
      "spGreTunnelDestination": 2
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "The GRE tunnel is back up",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "greDown"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "deviceSystemError",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.38",
    "trap-var-binds": {
      "spAlertID": 1,
      "spCollector": 2,
      "spSystemErrorType": 3,
      "spThreshold": 4,
      "spSystemErrorDescription": 5
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "An SP device system error alert has started.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "deviceSystemError"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "deviceSystemErrorDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.39",
    "trap-var-binds": {
      "spAlertID": 1,
      "spCollector": 2,
      "spSystemErrorType": 3,
      "spThreshold": 4,
      "spSystemErrorDescription": 5
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "An SP device system error alert has ended.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "deviceSystemError"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "fingerprintUsage",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.40",
    "trap-var-binds": {
      "spAlertID": 1,
      "spFingerprintName": 2,
      "spUsageType": 3,
      "spThreshold": 4,
      "spUsage": 5,
      "spUnit": 6
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The fingerprint exceeded the configured traffic rate threshold.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "fingerprintUsage"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "fingerprintUsageDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.41",
    "trap-var-binds": {
      "spAlertID": 1,
      "spFingerprintName": 2,
      "spUsageType": 3,
      "spUnit": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "The fingerprint is no longer exceeding the traffic rate threshold.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "fingerprintUsage"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "serviceUsage",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.42",
    "trap-var-binds": {
      "spAlertID": 1,
      "spServiceName": 2,
      "spUsageType": 3,
      "spApplicationName": 4,
      "spServiceElement": 5,
      "spThreshold": 6,
      "spUsage": 7
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The service exceeded the configured traffic rate threshold.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "serviceUsage"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "serviceUsageDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.43",
    "trap-var-binds": {
      "spAlertID": 1,
      "spServiceName": 2
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "The service is no longer exceeding the traffic rate threshold.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "serviceUsage"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "dosNetworkProfiledAlert",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.44",
    "trap-var-binds": {
      "spAlertID": 1,
      "pdosAnomalyClassification": 2,
      "pdosAnomalyDirection": 3,
      "pdosAnomalyStart": 4,
      "pdosAnomalyDuration": 5,
      "pdosUrl": 6,
      "spImpactBps": 7,
      "spImpactPps": 8,
      "spManagedObject": 9,
      "spDetectedCountries": 10
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "The managed object exceeded network and/or country baseline thresholds.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "dosNetworkProfiledAlert"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "dosNetworkProfiledAlertDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.45",
    "trap-var-binds": {
      "spAlertID": 1,
      "pdosAnomalyClassification": 2,
      "pdosAnomalyDirection": 3,
      "pdosAnomalyStart": 4,
      "pdosAnomalyDuration": 5,
      "pdosUrl": 6,
      "spImpactBps": 7,
      "spImpactPps": 8,
      "spManagedObject": 9,
      "spDetectedCountries": 10
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "The managed object is no longer exceeding the network and/or country baseline thresholds.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "dosNetworkProfiledAlert"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "dosHostDetectionAlert",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.46",
    "trap-var-binds": {
      "spAlertID": 1,
      "spAlertDetectionSignatures": 2,
      "pdosAnomalyDirection": 3,
      "pdosAnomalyStart": 4,
      "pdosAnomalyDuration": 5,
      "pdosUrl": 6,
      "spInetAddress": 7,
      "spInetAddressType": 8,
      "spImpactBps": 9,
      "spImpactPps": 10,
      "pdosAnomalyClassification": 11,
      "spManagedObjects": 12
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A Host alert was started after one or more signatures exceeded their thresholds.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "dosHostDetectionAlert"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "dosHostDetectionAlertDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.47",
    "trap-var-binds": {
      "spAlertID": 1,
      "spAlertDetectionSignatures": 2,
      "pdosAnomalyDirection": 3,
      "pdosAnomalyStart": 4,
      "pdosAnomalyDuration": 5,
      "pdosUrl": 6,
      "spInetAddress": 7,
      "spInetAddressType": 8,
      "spImpactBps": 9,
      "spImpactPps": 10,
      "pdosAnomalyClassification": 11,
      "spManagedObjects": 12
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "The Host alert ended and is no longer exceeding signature thresholds.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "dosHostDetectionAlert"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "routingFailover",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.48",
    "trap-var-binds": {
      "spAlertID": 1,
      "pdosUrl": 2,
      "spCollector": 3
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A routing failover event occurred on a collector.",
      "mib_name": "PEAKFLOW-SP-MIB"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "routingFailoverInterfaceDownAlert",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.49",
    "trap-var-binds": {
      "spAlertID": 1,
      "pdosUrl": 2,
      "spCollector": 3,
      "spRoutingFailoverInterfaces": 4
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "One or more interfaces involved in routing failover for a collector are down. The spRoutingFailoverInterfaces object documents the list of interfaces that are down.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "routingFailoverInterfaceDownAlert"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "routingFailoverInterfaceDownAlertDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.50",
    "trap-var-binds": {
      "spAlertID": 1,
      "pdosUrl": 2,
      "spCollector": 3
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "A previously started alert for down routing failover interfaces has finished.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "routingFailoverInterfaceDownAlert"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "trafficAutoMitigation",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.51",
    "trap-var-binds": {
      "spManagedObject": 1
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "Traffic has been seen by a TMS for a Managed Object which has been configured for traffic-based auto-mitigation.",
      "mib_name": "PEAKFLOW-SP-MIB"
    },
    "map-status": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "cloudSignalingMitigationRequest",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.52",
    "trap-var-binds": {
      "spManagedObject": 1,
      "spPravail": 2
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "A cloud signaling mitigation request has been seen and an alert created.",
      "mib_name": "PEAKFLOW-SP-MIB"
    },
    "map-status": {},
    "primary": "snmp_source_ip",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "licenseError",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.53",
    "trap-var-binds": {
      "spAlertID": 1,
      "spLicenseErrType": 2,
      "spLicenseErrCount": 3,
      "spThreshold": 4,
      "spLicenseErrDescription": 5
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "critical",
      "description": "An SP Deployment License error alert has started.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "licenseError"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  },
  {
    "type": "NO_MIB",
    "trap": "licenseErrorDone",
    "trap-oid": "1.3.6.1.4.1.9694.1.4.3.0.54",
    "trap-var-binds": {
      "spAlertID": 1,
      "spLicenseErrType": 2,
      "spThreshold": 3,
      "spLicenseErrLimit": 4,
      "spLicenseErrDescription": 5
    },
    "copy": {},
    "rename": {},
    "set": {
      "status": "ok",
      "description": "An SP Deployment License error alert has ended.",
      "mib_name": "PEAKFLOW-SP-MIB",
      "base_trap_name": "licenseError"
    },
    "map-status": {},
    "primary": "spAlertID",
    "secondary": "snmp_trap_name"
  }
]