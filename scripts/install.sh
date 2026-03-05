#!/bin/bash
#
# BigPanda Super Agent Installer
# Copyright (c) 2026 BigPanda Inc.
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
AGENT_USER="bigpanda"
AGENT_GROUP="bigpanda"
INSTALL_DIR="/opt/bigpanda-agent"
CONFIG_DIR="/etc/bigpanda-agent"
DATA_DIR="/var/lib/bigpanda-agent"
LOG_DIR="/var/log/bigpanda-agent"
SYSTEMD_DIR="/etc/systemd/system"

# Functions
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_root() {
    if [ "$EUID" -ne 0 ]; then
        print_error "This script must be run as root"
        exit 1
    fi
}

check_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$ID
        VER=$VERSION_ID
        print_info "Detected OS: $OS $VER"
    else
        print_error "Cannot detect OS"
        exit 1
    fi
}

create_user() {
    if id "$AGENT_USER" &>/dev/null; then
        print_info "User $AGENT_USER already exists"
    else
        print_info "Creating user $AGENT_USER"
        useradd -r -s /bin/false -d /nonexistent -M "$AGENT_USER"
        print_success "User created"
    fi
}

create_directories() {
    print_info "Creating directories"

    mkdir -p "$INSTALL_DIR"
    mkdir -p "$CONFIG_DIR"/{modules,snmp/{event_configs,mibs},certs}
    mkdir -p "$DATA_DIR"/{queue,state}
    mkdir -p "$LOG_DIR"

    print_success "Directories created"
}

install_binary() {
    print_info "Installing binary"

    if [ -f "./bigpanda-agent" ]; then
        cp ./bigpanda-agent "$INSTALL_DIR/"
        chmod 755 "$INSTALL_DIR/bigpanda-agent"
        print_success "Binary installed to $INSTALL_DIR/bigpanda-agent"
    else
        print_error "Binary not found in current directory"
        exit 1
    fi
}

install_configs() {
    print_info "Installing configuration files"

    # Main config
    if [ ! -f "$CONFIG_DIR/config.yaml" ]; then
        if [ -f "./configs/default.yaml" ]; then
            cp ./configs/default.yaml "$CONFIG_DIR/config.yaml"
            print_success "Installed default config"
        fi
    else
        print_warning "Config already exists, skipping"
    fi

    # Module configs
    if [ -d "./configs/modules" ]; then
        cp -r ./configs/modules/* "$CONFIG_DIR/modules/" 2>/dev/null || true
        print_success "Installed module configs"
    fi

    # Event configs
    if [ -d "./configs/event_configs" ]; then
        cp -r ./configs/event_configs/* "$CONFIG_DIR/snmp/event_configs/" 2>/dev/null || true
        print_success "Installed event configs"
    fi
}

install_systemd() {
    print_info "Installing systemd service"

    if [ -f "./configs/systemd/bigpanda-agent.service" ]; then
        cp ./configs/systemd/bigpanda-agent.service "$SYSTEMD_DIR/"
        systemctl daemon-reload
        print_success "Systemd service installed"
    else
        print_warning "Systemd service file not found"
    fi
}

set_permissions() {
    print_info "Setting permissions"

    chown -R "$AGENT_USER:$AGENT_GROUP" "$CONFIG_DIR"
    chown -R "$AGENT_USER:$AGENT_GROUP" "$DATA_DIR"
    chown -R "$AGENT_USER:$AGENT_GROUP" "$LOG_DIR"

    # Config files should be readable by agent user
    chmod 750 "$CONFIG_DIR"
    chmod 640 "$CONFIG_DIR"/*.yaml 2>/dev/null || true

    # Data and log directories
    chmod 750 "$DATA_DIR"
    chmod 750 "$LOG_DIR"

    print_success "Permissions set"
}

configure_firewall() {
    print_info "Configuring firewall"

    if command -v ufw &> /dev/null; then
        print_info "Detected UFW firewall"
        read -p "Open SNMP trap port (162/udp)? [y/N] " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            ufw allow 162/udp
            print_success "Opened port 162/udp"
        fi

        read -p "Open Web UI port (8443/tcp)? [y/N] " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            ufw allow 8443/tcp
            print_success "Opened port 8443/tcp"
        fi
    elif command -v firewall-cmd &> /dev/null; then
        print_info "Detected firewalld"
        read -p "Open SNMP trap port (162/udp)? [y/N] " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            firewall-cmd --permanent --add-port=162/udp
            firewall-cmd --reload
            print_success "Opened port 162/udp"
        fi

        read -p "Open Web UI port (8443/tcp)? [y/N] " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            firewall-cmd --permanent --add-port=8443/tcp
            firewall-cmd --reload
            print_success "Opened port 8443/tcp"
        fi
    else
        print_warning "No supported firewall detected, please configure manually"
    fi
}

configure_agent() {
    print_info "Agent configuration"
    echo

    read -p "BigPanda API Token (or press Enter to skip): " BP_TOKEN
    read -p "BigPanda App Key (or press Enter to skip): " BP_APP_KEY

    if [ -n "$BP_TOKEN" ] && [ -n "$BP_APP_KEY" ]; then
        # Update config file
        sed -i "s/\${BP_TOKEN}/$BP_TOKEN/g" "$CONFIG_DIR/config.yaml"
        sed -i "s/\${BP_APP_KEY}/$BP_APP_KEY/g" "$CONFIG_DIR/config.yaml"
        print_success "Credentials configured"
    else
        print_warning "Skipped credential configuration"
        print_info "Edit $CONFIG_DIR/config.yaml manually to add credentials"
    fi
}

enable_service() {
    print_info "Service configuration"
    echo

    read -p "Enable BigPanda Agent to start on boot? [Y/n] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Nn]$ ]]; then
        systemctl enable bigpanda-agent
        print_success "Service enabled"
    fi

    read -p "Start BigPanda Agent now? [Y/n] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Nn]$ ]]; then
        systemctl start bigpanda-agent
        sleep 2
        if systemctl is-active --quiet bigpanda-agent; then
            print_success "Agent started successfully"
        else
            print_error "Failed to start agent"
            print_info "Check logs: journalctl -u bigpanda-agent -f"
        fi
    fi
}

print_summary() {
    echo
    echo "=================================="
    echo "Installation Complete!"
    echo "=================================="
    echo
    echo "Binary:        $INSTALL_DIR/bigpanda-agent"
    echo "Config:        $CONFIG_DIR/config.yaml"
    echo "Data:          $DATA_DIR"
    echo "Logs:          $LOG_DIR"
    echo
    echo "Useful commands:"
    echo "  Start:       sudo systemctl start bigpanda-agent"
    echo "  Stop:        sudo systemctl stop bigpanda-agent"
    echo "  Restart:     sudo systemctl restart bigpanda-agent"
    echo "  Status:      sudo systemctl status bigpanda-agent"
    echo "  Logs:        sudo journalctl -u bigpanda-agent -f"
    echo "  Health:      curl http://localhost:8443/health"
    echo
    echo "Next steps:"
    echo "  1. Edit config: sudo nano $CONFIG_DIR/config.yaml"
    echo "  2. Add credentials if not done during install"
    echo "  3. Configure modules: $CONFIG_DIR/modules/"
    echo "  4. Access Web UI: https://localhost:8443"
    echo
}

# Main installation flow
main() {
    echo "=================================="
    echo "BigPanda Super Agent Installer"
    echo "=================================="
    echo

    check_root
    check_os

    print_info "Starting installation..."
    echo

    create_user
    create_directories
    install_binary
    install_configs
    install_systemd
    set_permissions
    configure_firewall
    configure_agent
    enable_service

    print_summary
}

# Run installation
main
