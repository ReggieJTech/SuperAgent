#!/bin/bash
#
# BigPanda Super Agent Uninstaller
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Configuration
AGENT_USER="bigpanda"
INSTALL_DIR="/opt/bigpanda-agent"
CONFIG_DIR="/etc/bigpanda-agent"
DATA_DIR="/var/lib/bigpanda-agent"
LOG_DIR="/var/log/bigpanda-agent"
SYSTEMD_SERVICE="/etc/systemd/system/bigpanda-agent.service"

print_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
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

confirm_uninstall() {
    echo "=================================="
    echo "BigPanda Super Agent Uninstaller"
    echo "=================================="
    echo
    print_info "This will remove the BigPanda Super Agent from your system."
    echo
    read -p "Are you sure you want to continue? [y/N] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Uninstall cancelled."
        exit 0
    fi
}

stop_service() {
    if systemctl is-active --quiet bigpanda-agent; then
        print_info "Stopping service..."
        systemctl stop bigpanda-agent
        print_success "Service stopped"
    fi

    if systemctl is-enabled --quiet bigpanda-agent 2>/dev/null; then
        print_info "Disabling service..."
        systemctl disable bigpanda-agent
        print_success "Service disabled"
    fi
}

remove_service() {
    if [ -f "$SYSTEMD_SERVICE" ]; then
        print_info "Removing systemd service..."
        rm -f "$SYSTEMD_SERVICE"
        systemctl daemon-reload
        print_success "Service removed"
    fi
}

remove_files() {
    print_info "Removing files..."

    # Binary
    if [ -d "$INSTALL_DIR" ]; then
        rm -rf "$INSTALL_DIR"
        print_success "Removed $INSTALL_DIR"
    fi

    # Ask about config
    read -p "Remove configuration files? [y/N] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        if [ -d "$CONFIG_DIR" ]; then
            rm -rf "$CONFIG_DIR"
            print_success "Removed $CONFIG_DIR"
        fi
    else
        print_info "Kept configuration files in $CONFIG_DIR"
    fi

    # Ask about data
    read -p "Remove data files (queue, state)? [y/N] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        if [ -d "$DATA_DIR" ]; then
            rm -rf "$DATA_DIR"
            print_success "Removed $DATA_DIR"
        fi
    else
        print_info "Kept data files in $DATA_DIR"
    fi

    # Ask about logs
    read -p "Remove log files? [y/N] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        if [ -d "$LOG_DIR" ]; then
            rm -rf "$LOG_DIR"
            print_success "Removed $LOG_DIR"
        fi
    else
        print_info "Kept log files in $LOG_DIR"
    fi
}

remove_user() {
    if id "$AGENT_USER" &>/dev/null; then
        read -p "Remove user '$AGENT_USER'? [y/N] " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            userdel "$AGENT_USER" 2>/dev/null || true
            print_success "User removed"
        fi
    fi
}

main() {
    check_root
    confirm_uninstall

    echo
    stop_service
    remove_service
    remove_files
    remove_user

    echo
    print_success "Uninstall complete!"
    echo
}

main
