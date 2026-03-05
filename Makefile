.PHONY: all build clean test lint install docker help

# Build variables
BINARY_NAME=bigpanda-agent
VERSION?=0.1.0
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildTime=$(BUILD_TIME) -w -s"

# Directories
BUILD_DIR=build
DIST_DIR=dist
SRC_DIR=cmd/agent

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Default target
all: clean build

## help: Show this help message
help:
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## build: Build the agent binary
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(SRC_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## build-all: Build binaries for all platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p $(DIST_DIR)
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./$(SRC_DIR)
	# Linux ARM64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 ./$(SRC_DIR)
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(SRC_DIR)
	# macOS ARM64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(SRC_DIR)
	# Windows AMD64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(SRC_DIR)
	@echo "Build complete for all platforms"

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...

## test-short: Run short tests
test-short:
	@echo "Running short tests..."
	$(GOTEST) -v -short ./...

## coverage: Generate test coverage report
coverage: test
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report: coverage.html"

## lint: Run linters
lint:
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	$(GOCLEAN)
	@rm -f coverage.txt coverage.html

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## vendor: Vendor dependencies
vendor:
	@echo "Vendoring dependencies..."
	$(GOMOD) vendor

## run: Build and run the agent
run: build
	@echo "Running agent..."
	./$(BUILD_DIR)/$(BINARY_NAME) -config configs/default.yaml

## install: Install the agent binary
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "Installed successfully"

## docker: Build Docker image
docker:
	@echo "Building Docker image..."
	docker build -t bigpanda/super-agent:$(VERSION) -t bigpanda/super-agent:latest .
	@echo "Docker image built: bigpanda/super-agent:$(VERSION)"

## docker-run: Run agent in Docker
docker-run:
	docker run --rm -it \
		-p 8443:8443 \
		-p 162:162/udp \
		-p 8080:8080 \
		-v $(PWD)/configs:/etc/bigpanda-agent \
		bigpanda/super-agent:latest

## package: Create distribution packages
package: build-all
	@echo "Creating distribution packages..."
	@mkdir -p $(DIST_DIR)/packages
	# Linux AMD64 tarball
	tar -czf $(DIST_DIR)/packages/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz \
		-C $(DIST_DIR) $(BINARY_NAME)-linux-amd64 \
		-C ../configs .
	# Linux ARM64 tarball
	tar -czf $(DIST_DIR)/packages/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz \
		-C $(DIST_DIR) $(BINARY_NAME)-linux-arm64 \
		-C ../configs .
	@echo "Packages created in $(DIST_DIR)/packages"

## validate: Validate configuration
validate: build
	./$(BUILD_DIR)/$(BINARY_NAME) -config configs/default.yaml -validate

## version: Show version
version: build
	./$(BUILD_DIR)/$(BINARY_NAME) -version
