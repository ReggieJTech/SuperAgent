[
  {
    "type": "NO_MIB",
    "trap": "diskFailed",
    "trap-oid": "1.3.6.1.4.1.789.0.3",
    "trap-var-binds": {
      "diskName": 1,
      "diskSerial": 2
    },
    "copy": {},
    "rename": {
      "diskName": "disk_name",
      "diskSerial": "disk_serial"
    },
    "set": {
      "status": "critical",
      "description": "NetApp disk has failed",
      "mib_name": "NETAPP-MIB",
      "vendor": "netapp"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "disk_name"
  },
  {
    "type": "NO_MIB",
    "trap": "volumeFull",
    "trap-oid": "1.3.6.1.4.1.789.0.5",
    "trap-var-binds": {
      "volumeName": 1,
      "volumeUsed": 2
    },
    "copy": {},
    "rename": {
      "volumeName": "volume_name",
      "volumeUsed": "volume_used_percent"
    },
    "set": {
      "status": "critical",
      "description": "NetApp volume is full",
      "mib_name": "NETAPP-MIB",
      "vendor": "netapp"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "volume_name"
  },
  {
    "type": "NO_MIB",
    "trap": "volumeNearlyFull",
    "trap-oid": "1.3.6.1.4.1.789.0.6",
    "trap-var-binds": {
      "volumeName": 1,
      "volumeUsed": 2
    },
    "copy": {},
    "rename": {
      "volumeName": "volume_name",
      "volumeUsed": "volume_used_percent"
    },
    "set": {
      "status": "warning",
      "description": "NetApp volume is nearly full",
      "mib_name": "NETAPP-MIB",
      "vendor": "netapp"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "volume_name"
  },
  {
    "type": "NO_MIB",
    "trap": "fanFailed",
    "trap-oid": "1.3.6.1.4.1.789.0.8",
    "trap-var-binds": {
      "fanNumber": 1
    },
    "copy": {},
    "rename": {
      "fanNumber": "fan_number"
    },
    "set": {
      "status": "critical",
      "description": "NetApp fan has failed",
      "mib_name": "NETAPP-MIB",
      "vendor": "netapp"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "fan_number"
  }
]
