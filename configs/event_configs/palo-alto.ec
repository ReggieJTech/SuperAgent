[
  {
    "type": "NO_MIB",
    "trap": "panCommonEventLog",
    "trap-oid": "1.3.6.1.4.1.25461.2.1.3.2.0.1",
    "trap-var-binds": {
      "panEventType": 1,
      "panEventSubType": 2,
      "panVsys": 3,
      "panSeqno": 4,
      "panActionflags": 5,
      "panSystemEventId": 6,
      "panSystemObject": 7,
      "panSystemModule": 8,
      "panSystemSeverity": 9,
      "panSystemDescription": 10
    },
    "copy": {},
    "rename": {
      "panEventType": "event_type",
      "panSystemSeverity": "severity",
      "panSystemDescription": "description",
      "panSystemModule": "module"
    },
    "set": {
      "mib_name": "PAN-COMMON-MIB",
      "vendor": "palo_alto"
    },
    "map-status": {
      "severity": {
        "critical": "critical",
        "high": "critical",
        "medium": "warning",
        "low": "warning",
        "informational": "ok"
      }
    },
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "module"
  },
  {
    "type": "NO_MIB",
    "trap": "panGlobalProtectGatewayUp",
    "trap-oid": "1.3.6.1.4.1.25461.2.1.3.2.0.13",
    "trap-var-binds": {
      "panGPGWName": 1
    },
    "copy": {},
    "rename": {
      "panGPGWName": "gateway_name"
    },
    "set": {
      "status": "ok",
      "description": "Palo Alto GlobalProtect gateway is up",
      "mib_name": "PAN-COMMON-MIB",
      "vendor": "palo_alto"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "gateway_name"
  },
  {
    "type": "NO_MIB",
    "trap": "panGlobalProtectGatewayDown",
    "trap-oid": "1.3.6.1.4.1.25461.2.1.3.2.0.14",
    "trap-var-binds": {
      "panGPGWName": 1
    },
    "copy": {},
    "rename": {
      "panGPGWName": "gateway_name"
    },
    "set": {
      "status": "critical",
      "description": "Palo Alto GlobalProtect gateway is down",
      "mib_name": "PAN-COMMON-MIB",
      "vendor": "palo_alto"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "gateway_name"
  }
]
