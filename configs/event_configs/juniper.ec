[
  {
    "type": "NO_MIB",
    "trap": "jnxPowerSupplyFailure",
    "trap-oid": "1.3.6.1.4.1.2636.4.1.2",
    "trap-var-binds": {
      "jnxContentsDescr": 1,
      "jnxContentsSerialNo": 2
    },
    "copy": {},
    "rename": {
      "jnxContentsDescr": "psu_description",
      "jnxContentsSerialNo": "psu_serial"
    },
    "set": {
      "status": "critical",
      "description": "Juniper power supply failure",
      "mib_name": "JUNIPER-MIB",
      "vendor": "juniper"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "psu_description"
  },
  {
    "type": "NO_MIB",
    "trap": "jnxFanFailure",
    "trap-oid": "1.3.6.1.4.1.2636.4.1.3",
    "trap-var-binds": {
      "jnxContentsDescr": 1
    },
    "copy": {},
    "rename": {
      "jnxContentsDescr": "fan_description"
    },
    "set": {
      "status": "critical",
      "description": "Juniper fan failure",
      "mib_name": "JUNIPER-MIB",
      "vendor": "juniper"
    },
    "map-status": {},
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "fan_description"
  },
  {
    "type": "NO_MIB",
    "trap": "jnxBgpPeerStateChange",
    "trap-oid": "1.3.6.1.4.1.2636.4.5.1",
    "trap-var-binds": {
      "jnxBgpPeerRemoteAddr": 1,
      "jnxBgpPeerState": 2
    },
    "copy": {},
    "rename": {
      "jnxBgpPeerRemoteAddr": "peer_address",
      "jnxBgpPeerState": "peer_state"
    },
    "set": {
      "status": "warning",
      "description": "Juniper BGP peer state changed",
      "mib_name": "JUNIPER-BGP4-MIB",
      "vendor": "juniper"
    },
    "map-status": {
      "peer_state": {
        "1": "critical",
        "2": "critical",
        "6": "ok"
      }
    },
    "conditions": [],
    "custom-actions": {},
    "primary": "snmp_source_ip",
    "secondary": "peer_address"
  }
]
