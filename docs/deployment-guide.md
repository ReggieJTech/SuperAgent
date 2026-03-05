# Deployment Guide

This guide covers all deployment methods for the BigPanda Super Agent.

## Table of Contents

- [Linux Installation](#linux-installation)
- [Docker Deployment](#docker-deployment)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Configuration](#configuration)
- [Post-Installation](#post-installation)

## Linux Installation

### Prerequisites

- Linux OS (Ubuntu, Debian, RHEL, CentOS, or similar)
- Root or sudo access
- 1GB RAM minimum
- 2GB disk space
- Open ports: 162/UDP (SNMP), 8443/TCP (Web UI)

### Installation Steps

1. **Download Release**

```bash
wget https://github.com/ReggieJTech/SuperAgent/releases/latest/download/bigpanda-agent-linux-amd64.tar.gz
tar -xzf bigpanda-agent-linux-amd64.tar.gz
cd bigpanda-agent
```

2. **Run Installer**

```bash
sudo ./scripts/install.sh
```

The installer will:
- Create `bigpanda` user
- Install binary to `/opt/bigpanda-agent`
- Create directories in `/etc`, `/var/lib`, `/var/log`
- Install systemd service
- Configure firewall (optional)
- Prompt for BigPanda credentials

3. **Verify Installation**

```bash
sudo systemctl status bigpanda-agent
curl -k https://localhost:8443/health
```

### Manual Installation

If you prefer manual installation:

```bash
# Create user
sudo useradd -r -s /bin/false bigpanda

# Create directories
sudo mkdir -p /opt/bigpanda-agent
sudo mkdir -p /etc/bigpanda-agent/{modules,snmp/{event_configs,mibs},certs}
sudo mkdir -p /var/lib/bigpanda-agent/{queue,state}
sudo mkdir -p /var/log/bigpanda-agent

# Install binary
sudo cp bigpanda-agent /opt/bigpanda-agent/
sudo chmod 755 /opt/bigpanda-agent/bigpanda-agent

# Install configs
sudo cp configs/default.yaml /etc/bigpanda-agent/config.yaml
sudo cp -r configs/modules/* /etc/bigpanda-agent/modules/
sudo cp -r configs/event_configs/* /etc/bigpanda-agent/snmp/event_configs/

# Install systemd service
sudo cp configs/systemd/bigpanda-agent.service /etc/systemd/system/
sudo systemctl daemon-reload

# Set permissions
sudo chown -R bigpanda:bigpanda /opt/bigpanda-agent
sudo chown -R bigpanda:bigpanda /etc/bigpanda-agent
sudo chown -R bigpanda:bigpanda /var/lib/bigpanda-agent
sudo chown -R bigpanda:bigpanda /var/log/bigpanda-agent

# Enable and start
sudo systemctl enable bigpanda-agent
sudo systemctl start bigpanda-agent
```

### Uninstallation

```bash
sudo ./scripts/uninstall.sh
```

Or manually:

```bash
sudo systemctl stop bigpanda-agent
sudo systemctl disable bigpanda-agent
sudo rm /etc/systemd/system/bigpanda-agent.service
sudo rm -rf /opt/bigpanda-agent
sudo rm -rf /etc/bigpanda-agent
sudo rm -rf /var/lib/bigpanda-agent
sudo rm -rf /var/log/bigpanda-agent
sudo userdel bigpanda
```

## Docker Deployment

### Prerequisites

- Docker 20.10+
- Docker Compose (optional)
- 1GB RAM minimum

### Quick Start

```bash
docker run -d \
  --name bigpanda-agent \
  -p 8443:8443 \
  -p 162:162/udp \
  -e BP_TOKEN=your_token \
  -e BP_APP_KEY=your_app_key \
  -v /opt/bigpanda-config:/etc/bigpanda-agent \
  -v /opt/bigpanda-data:/var/lib/bigpanda-agent \
  reggiejtech/super-agent:latest
```

### Docker Compose

```bash
# Clone repository or download docker-compose.yml
wget https://raw.githubusercontent.com/reggiejtech/super-agent/main/docker-compose.yml

# Create .env file
cat > .env <<EOF
BP_TOKEN=your_token
BP_APP_KEY=your_app_key
EOF

# Start
docker-compose up -d

# Check logs
docker-compose logs -f

# Stop
docker-compose down
```

### Build Custom Image

```bash
git clone https://github.com/ReggieJTech/SuperAgent.git
cd super-agent
docker build -t reggiejtech/super-agent:custom .
```

### Docker Volume Management

```bash
# Backup queue data
docker run --rm -v bigpanda-agent_agent-data:/data -v $(pwd):/backup alpine \
  tar czf /backup/bigpanda-data-backup.tar.gz -C /data .

# Restore queue data
docker run --rm -v bigpanda-agent_agent-data:/data -v $(pwd):/backup alpine \
  tar xzf /backup/bigpanda-data-backup.tar.gz -C /data
```

## Kubernetes Deployment

### Prerequisites

- Kubernetes 1.20+
- kubectl configured
- Helm 3+ (optional)
- Persistent storage provider

### Quick Start

```bash
# Create namespace
kubectl create namespace bigpanda

# Create secret with credentials
kubectl create secret generic bigpanda-credentials \
  --from-literal=token=YOUR_TOKEN \
  --from-literal=app_key=YOUR_APP_KEY \
  -n bigpanda

# Deploy
kubectl apply -k k8s/ -n bigpanda

# Verify
kubectl get pods -n bigpanda
kubectl logs -f deployment/bigpanda-agent -n bigpanda
```

### Helm Deployment (Future)

```bash
helm repo add bigpanda https://charts.bigpanda.io
helm install bigpanda-agent reggiejtech/super-agent \
  --set bigpanda.token=YOUR_TOKEN \
  --set bigpanda.appKey=YOUR_APP_KEY \
  --namespace bigpanda \
  --create-namespace
```

### Exposing SNMP Traps

For SNMP trap reception, expose UDP port 162:

**NodePort:**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: bigpanda-agent-snmp
spec:
  type: NodePort
  ports:
  - port: 162
    targetPort: 162
    nodePort: 30162
    protocol: UDP
```

**LoadBalancer:**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: bigpanda-agent-snmp
spec:
  type: LoadBalancer
  ports:
  - port: 162
    targetPort: 162
    protocol: UDP
```

## Configuration

### BigPanda Credentials

Edit `/etc/bigpanda-agent/config.yaml`:

```yaml
bigpanda:
  api_url: "https://api.bigpanda.io/data/v2/alerts"
  token: "YOUR_TOKEN"
  app_key: "YOUR_APP_KEY"
```

Or use environment variables:

```bash
export BP_TOKEN="your_token"
export BP_APP_KEY="your_app_key"
```

### Enable SNMP Module

Edit `/etc/bigpanda-agent/modules/snmp.yaml`:

```yaml
listen_address: "0.0.0.0:162"
snmp_version: "2c"
community: "public"
event_configs_dir: "/etc/bigpanda-agent/snmp/event_configs"
```

### TLS Certificates

**Auto-generated (default):**
```yaml
server:
  tls:
    enabled: true
    auto_generate: true
```

**Custom certificates:**
```yaml
server:
  tls:
    enabled: true
    cert_file: "/etc/bigpanda-agent/certs/server.crt"
    key_file: "/etc/bigpanda-agent/certs/server.key"
```

Generate self-signed:
```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /etc/bigpanda-agent/certs/server.key \
  -out /etc/bigpanda-agent/certs/server.crt
```

## Post-Installation

### Verify Agent Health

```bash
# Check health
curl -k https://localhost:8443/health

# Check statistics
curl -k https://localhost:8443/api/v1/stats | jq

# View logs
sudo journalctl -u bigpanda-agent -f
```

### Test SNMP Reception

```bash
# Send test trap
snmptrap -v 2c -c public localhost:162 '' \
  1.3.6.1.6.3.1.1.5.3 \
  1.3.6.1.2.1.2.2.1.1 i 1

# Check logs for received trap
sudo journalctl -u bigpanda-agent | grep "Trap received"
```

### Access Web UI

Open browser to: `https://localhost:8443`

Default credentials: (to be configured)

### Configure Firewall

**UFW (Ubuntu/Debian):**
```bash
sudo ufw allow 162/udp
sudo ufw allow 8443/tcp
```

**firewalld (RHEL/CentOS):**
```bash
sudo firewall-cmd --permanent --add-port=162/udp
sudo firewall-cmd --permanent --add-port=8443/tcp
sudo firewall-cmd --reload
```

### Monitor Performance

```bash
# CPU and memory usage
top -p $(pgrep bigpanda-agent)

# Queue size
curl -k https://localhost:8443/api/v1/queue/size

# Plugin statistics
curl -k https://localhost:8443/api/v1/plugins/snmp/stats
```

## Upgrading

### Linux

```bash
# Stop service
sudo systemctl stop bigpanda-agent

# Backup config
sudo cp -r /etc/bigpanda-agent /etc/bigpanda-agent.backup

# Install new version
sudo ./scripts/install.sh

# Start service
sudo systemctl start bigpanda-agent
```

### Docker

```bash
# Pull new image
docker pull reggiejtech/super-agent:latest

# Recreate container
docker-compose up -d
```

### Kubernetes

```bash
kubectl set image deployment/bigpanda-agent \
  agent=reggiejtech/super-agent:v1.1.0 \
  -n bigpanda
```

## Troubleshooting

### Agent Won't Start

```bash
# Check service status
sudo systemctl status bigpanda-agent

# Check logs
sudo journalctl -u bigpanda-agent -n 100

# Validate config
/opt/bigpanda-agent/bigpanda-agent -config /etc/bigpanda-agent/config.yaml -validate
```

### No Traps Received

```bash
# Check if port 162 is listening
sudo netstat -ulnp | grep 162

# Check firewall
sudo ufw status
sudo iptables -L -n | grep 162

# Test locally
snmptrap -v 2c -c public localhost:162 '' 1.3.6.1.6.3.1.1.5.1
```

### High Memory Usage

```bash
# Check queue size
curl -k https://localhost:8443/api/v1/queue/size

# Adjust queue limits in config
queue:
  max_size: 50000  # Reduce if needed
```

### Connection to BigPanda Fails

```bash
# Test connectivity
curl -v https://api.bigpanda.io/data/v2/alerts

# Check credentials
curl -H "Authorization: Bearer YOUR_TOKEN" \
     -H "X-App-Key: YOUR_APP_KEY" \
     https://api.bigpanda.io/data/v2/alerts
```

## Support

- Documentation: https://docs.bigpanda.io/super-agent
- Issues: https://github.com/ReggieJTech/SuperAgent/issues
- Email: support@bigpanda.io
