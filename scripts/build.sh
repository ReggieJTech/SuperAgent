#!/bin/bash
#
# Build script for BigPanda Super Agent
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# Build info
VERSION=${VERSION:-"0.1.0"}
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

# Directories
BUILD_DIR="build"
DIST_DIR="dist"

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

build_binary() {
    local os=$1
    local arch=$2
    local output=$3

    print_info "Building for $os/$arch..."

    GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build \
        -ldflags "-X main.Version=$VERSION -X main.GitCommit=$GIT_COMMIT -X main.BuildTime=$BUILD_TIME -w -s" \
        -o "$output" \
        ./cmd/agent

    print_success "Built $output"
}

create_tarball() {
    local platform=$1
    local binary=$2

    print_info "Creating tarball for $platform..."

    # Create temporary directory
    local temp_dir="$DIST_DIR/tmp/$platform"
    mkdir -p "$temp_dir"

    # Copy binary
    cp "$binary" "$temp_dir/bigpanda-agent"
    chmod +x "$temp_dir/bigpanda-agent"

    # Copy configs
    mkdir -p "$temp_dir/configs"
    cp -r configs/* "$temp_dir/configs/"

    # Copy scripts
    mkdir -p "$temp_dir/scripts"
    cp scripts/install.sh "$temp_dir/scripts/"
    chmod +x "$temp_dir/scripts/install.sh"

    # Copy docs
    cp README.md "$temp_dir/" 2>/dev/null || true

    # Create tarball
    tar -czf "$DIST_DIR/bigpanda-agent-$VERSION-$platform.tar.gz" \
        -C "$DIST_DIR/tmp" "$platform"

    # Cleanup
    rm -rf "$temp_dir"

    print_success "Created tarball: bigpanda-agent-$VERSION-$platform.tar.gz"
}

build_docker() {
    print_info "Building Docker image..."

    docker build \
        --build-arg VERSION="$VERSION" \
        --build-arg GIT_COMMIT="$GIT_COMMIT" \
        --build-arg BUILD_TIME="$BUILD_TIME" \
        -t "bigpanda/super-agent:$VERSION" \
        -t "bigpanda/super-agent:latest" \
        .

    print_success "Docker image built: bigpanda/super-agent:$VERSION"
}

main() {
    echo "=================================="
    echo "BigPanda Super Agent Build Script"
    echo "=================================="
    echo "Version: $VERSION"
    echo "Commit:  $GIT_COMMIT"
    echo "Time:    $BUILD_TIME"
    echo "=================================="
    echo

    # Clean previous builds
    print_info "Cleaning previous builds..."
    rm -rf "$BUILD_DIR" "$DIST_DIR"
    mkdir -p "$BUILD_DIR" "$DIST_DIR/tmp"

    # Build for all platforms
    print_info "Building binaries..."
    echo

    build_binary "linux" "amd64" "$DIST_DIR/bigpanda-agent-linux-amd64"
    build_binary "linux" "arm64" "$DIST_DIR/bigpanda-agent-linux-arm64"
    build_binary "darwin" "amd64" "$DIST_DIR/bigpanda-agent-darwin-amd64"
    build_binary "darwin" "arm64" "$DIST_DIR/bigpanda-agent-darwin-arm64"
    build_binary "windows" "amd64" "$DIST_DIR/bigpanda-agent-windows-amd64.exe"

    echo
    print_info "Creating distribution packages..."
    echo

    # Create tarballs
    create_tarball "linux-amd64" "$DIST_DIR/bigpanda-agent-linux-amd64"
    create_tarball "linux-arm64" "$DIST_DIR/bigpanda-agent-linux-arm64"
    create_tarball "darwin-amd64" "$DIST_DIR/bigpanda-agent-darwin-amd64"
    create_tarball "darwin-arm64" "$DIST_DIR/bigpanda-agent-darwin-arm64"

    # Build Docker image
    echo
    if command -v docker &> /dev/null; then
        build_docker
    else
        print_info "Docker not found, skipping Docker build"
    fi

    # Cleanup
    rm -rf "$DIST_DIR/tmp"

    echo
    print_success "Build complete!"
    echo
    echo "Artifacts in $DIST_DIR/:"
    ls -lh "$DIST_DIR/"
    echo
}

# Run build
main "$@"
