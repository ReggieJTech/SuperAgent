# BigPanda Super Agent

A self-contained, modular, production-ready monitoring agent for BigPanda that handles multiple types of event ingestion.

## Features

- **Self-contained deployment**: Single binary with no external dependencies
- **Modular architecture**: Plugin-based receiver modules (SNMP, Webhook, Automation)
- **Web UI**: Configuration, monitoring, and management interface
- **Queue-based processing**: Reliable event handling with retry logic
- **Security**: Production-ready for customer networks with TLS, authentication, and encryption
- **Scale**: Handle 100-1000+ events/second
- **Multi-auth**: Support local users, LDAP, and SSO

## Architecture

The BigPanda Super Agent consists of:

1. **Core Agent**: Central service managing configuration, plugins, and lifecycle
2. **Event Queue**: Persistent BadgerDB-based queue with retry and DLQ support
3. **BigPanda Forwarder**: Batching, rate limiting, and retry logic for API delivery
4. **Receiver Modules**:
   - **SNMP**: Trap reception with MIB parsing and event configuration
   - **Webhook**: HTTP/HTTPS endpoints with authentication and transformation
   - **Automation** (Future): Bidirectional automation task execution
5. **Web UI**: React-based management interface

## Quick Start

### Prerequisites

- Linux system (Ubuntu, RHEL, CentOS, Debian)
- Port 162/UDP for SNMP traps
- Port 8443/TCP for Web UI
- Port 8080/TCP for Webhooks (optional)

### Installation

```bash
# Download latest release
wget https://github.com/bigpanda/super-agent/releases/latest/download/bigpanda-agent-linux-amd64.tar.gz

# Extract and install
tar -xzf bigpanda-agent-linux-amd64.tar.gz
cd bigpanda-agent
sudo ./install.sh

# Configure BigPanda credentials
sudo nano /etc/bigpanda-agent/config.yaml

# Start the agent
sudo systemctl start bigpanda-agent
sudo systemctl enable bigpanda-agent

# Check status
sudo systemctl status bigpanda-agent

# Access Web UI
https://localhost:8443
```

### Docker Deployment

```bash
docker run -d \
  --name bigpanda-agent \
  -p 8443:8443 \
  -p 162:162/udp \
  -p 8080:8080 \
  -v /opt/bigpanda-config:/etc/bigpanda-agent \
  -v /opt/bigpanda-data:/var/lib/bigpanda-agent \
  -e BP_TOKEN=your_token \
  -e BP_APP_KEY=your_app_key \
  bigpanda/super-agent:latest
```

## Configuration

Main configuration file: `/etc/bigpanda-agent/config.yaml`

```yaml
server:
  listen_address: "0.0.0.0:8443"
  tls:
    enabled: true

bigpanda:
  api_url: "https://api.bigpanda.io/data/v2/alerts"
  token: "${BP_TOKEN}"
  app_key: "${BP_APP_KEY}"

modules:
  - name: snmp
    enabled: true
    config_file: "/etc/bigpanda-agent/modules/snmp.yaml"

  - name: webhook
    enabled: true
    config_file: "/etc/bigpanda-agent/modules/webhook.yaml"
```

See [Configuration Guide](docs/configuration.md) for full details.

## Development

### Building from Source

```bash
# Clone repository
git clone https://github.com/bigpanda/super-agent.git
cd super-agent

# Build
make build

# Run tests
make test

# Build Docker image
make docker
```

### Project Structure

```
bigpanda-super-agent/
├── cmd/agent/          # Main entry point
├── internal/           # Internal packages
│   ├── agent/          # Core agent service
│   ├── queue/          # Event queue
│   ├── forwarder/      # BigPanda API client
│   ├── plugin/         # Plugin system
│   ├── modules/        # Receiver modules
│   ├── webui/          # Web UI backend
│   ├── monitoring/     # Health and metrics
│   └── security/       # Auth and encryption
├── web/                # React frontend
├── configs/            # Default configurations
├── scripts/            # Build and install scripts
└── docs/               # Documentation
```

## Documentation

- [Architecture Guide](docs/architecture.md)
- [Configuration Reference](docs/configuration.md)
- [SNMP Module Guide](docs/snmp.md)
- [Webhook Module Guide](docs/webhook.md)
- [MIB Management](docs/mib-management.md)
- [API Reference](docs/api.md)
- [Troubleshooting](docs/troubleshooting.md)

## Supported Vendors (SNMP)

Pre-bundled event configurations for 60+ vendors including:
- Cisco (IOS, NXOS, ASA)
- F5 BIG-IP
- NetApp
- Dell EMC
- HP/HPE
- Juniper
- Palo Alto Networks
- And many more...

## Security

- TLS encryption for all web interfaces
- Multi-factor authentication support
- Credential encryption at rest
- IP whitelisting for webhooks
- API key rotation
- Audit logging

## Support

- Issues: 
- Documentation: 
- Email: 

## License

Copyright © 2026 BigPanda Inc. All rights reserved.
