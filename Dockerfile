# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_TIME=unknown
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME} -w -s" \
    -o bigpanda-agent \
    ./cmd/agent

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create user and directories
RUN addgroup -g 1000 bigpanda && \
    adduser -D -u 1000 -G bigpanda bigpanda && \
    mkdir -p /opt/bigpanda-agent \
             /etc/bigpanda-agent/modules \
             /etc/bigpanda-agent/snmp/event_configs \
             /etc/bigpanda-agent/snmp/mibs \
             /etc/bigpanda-agent/certs \
             /var/lib/bigpanda-agent/queue \
             /var/lib/bigpanda-agent/state \
             /var/log/bigpanda-agent && \
    chown -R bigpanda:bigpanda /opt/bigpanda-agent \
                                /etc/bigpanda-agent \
                                /var/lib/bigpanda-agent \
                                /var/log/bigpanda-agent

# Copy binary from builder
COPY --from=builder /build/bigpanda-agent /opt/bigpanda-agent/

# Copy default configs
COPY --chown=bigpanda:bigpanda configs/default.yaml /etc/bigpanda-agent/config.yaml
COPY --chown=bigpanda:bigpanda configs/modules/ /etc/bigpanda-agent/modules/
COPY --chown=bigpanda:bigpanda configs/event_configs/ /etc/bigpanda-agent/snmp/event_configs/

# Set working directory
WORKDIR /opt/bigpanda-agent

# Switch to non-root user
USER bigpanda

# Expose ports
EXPOSE 8443/tcp
EXPOSE 162/udp
EXPOSE 8080/tcp

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD ["/opt/bigpanda-agent/bigpanda-agent", "-config", "/etc/bigpanda-agent/config.yaml", "-validate"] || exit 1

# Environment variables
ENV BP_TOKEN="" \
    BP_APP_KEY=""

# Run agent
ENTRYPOINT ["/opt/bigpanda-agent/bigpanda-agent"]
CMD ["-config", "/etc/bigpanda-agent/config.yaml"]
