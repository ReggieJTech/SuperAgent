# Plugin Development Guide

This guide explains how to develop plugins for the BigPanda Super Agent.

## Overview

Plugins are modular components that receive events from various sources (SNMP, webhooks, automation, etc.) and forward them to the agent's event queue for processing. The plugin system provides:

- **Standardized interface**: All plugins implement the same interface
- **Lifecycle management**: Automatic start, stop, and health monitoring
- **Error recovery**: Automatic restart on failures
- **Health monitoring**: Built-in health checks and metrics
- **Configuration**: YAML-based configuration per plugin

## Plugin Interface

Every plugin must implement the `Plugin` interface:

```go
type Plugin interface {
    // Name returns the plugin name
    Name() string

    // Version returns the plugin version
    Version() string

    // Description returns a short description
    Description() string

    // Init initializes the plugin with configuration
    Init(ctx context.Context, config PluginConfig) error

    // Start starts the plugin
    Start(ctx context.Context) error

    // Stop stops the plugin gracefully
    Stop(ctx context.Context) error

    // Health returns the plugin health status
    Health() HealthStatus

    // Stats returns plugin statistics
    Stats() map[string]interface{}
}
```

## Creating a Plugin

### 1. Use BasePlugin

The easiest way to create a plugin is to embed `BasePlugin`, which provides common functionality:

```go
package myplugin

import (
    "context"
    "github.com/ReggieJTech/SuperAgent/internal/plugin"
    "github.com/ReggieJTech/SuperAgent/internal/queue"
)

type MyPlugin struct {
    *plugin.BasePlugin
    config MyConfig
    stopChan chan struct{}
}

func NewMyPlugin() plugin.Plugin {
    base := plugin.NewBasePlugin("myplugin", "1.0.0", "My custom plugin")
    return &MyPlugin{
        BasePlugin: base,
        stopChan: make(chan struct{}),
    }
}
```

### 2. Implement Init

Parse your plugin configuration:

```go
type MyConfig struct {
    ListenAddress string `yaml:"listen_address"`
    // ... other config fields
}

func (p *MyPlugin) Init(ctx context.Context, config plugin.PluginConfig) error {
    // Call base init
    if err := p.BasePlugin.Init(ctx, config); err != nil {
        return err
    }

    // Parse plugin-specific config
    if config.Config != nil {
        // Convert map to struct
        // ... parse config
    }

    return nil
}
```

### 3. Implement Start

Start your plugin's main logic:

```go
func (p *MyPlugin) Start(ctx context.Context) error {
    if p.IsStarted() {
        return nil
    }

    log.Info().Str("plugin", p.Name()).Msg("Starting plugin")

    p.MarkStarted()
    p.SetHealthy()

    // Start your main goroutine
    go p.run(ctx)

    return nil
}

func (p *MyPlugin) run(ctx context.Context) {
    defer p.RecoverFromPanic()  // Important!

    for {
        select {
        case <-ctx.Done():
            return
        case <-p.stopChan:
            return
        default:
            // Your plugin logic here
            // Receive events, process them, and send to queue
            event := p.receiveEvent()
            if event != nil {
                p.sendToQueue(event)
            }
        }
    }
}
```

### 4. Implement Stop

Stop your plugin gracefully:

```go
func (p *MyPlugin) Stop(ctx context.Context) error {
    if !p.IsStarted() {
        return nil
    }

    log.Info().Str("plugin", p.Name()).Msg("Stopping plugin")

    close(p.stopChan)
    p.MarkStopped()

    return nil
}
```

### 5. Send Events to Queue

Get the queue from the config and enqueue events:

```go
func (p *MyPlugin) sendToQueue(event *queue.Event) error {
    q, ok := p.Queue().(queue.Queue)
    if !ok {
        return fmt.Errorf("invalid queue type")
    }

    // Increment metrics
    p.IncrementReceived()

    // Enqueue event
    if err := q.Enqueue(context.Background(), event); err != nil {
        p.IncrementDropped()
        p.IncrementErrors()
        return err
    }

    p.IncrementSent()
    return nil
}
```

## Health Monitoring

Use the built-in health methods:

```go
// Mark plugin as healthy
p.SetHealthy()

// Mark plugin as degraded (partial functionality)
p.SetDegraded("Some connections failing", map[string]interface{}{
    "failed_connections": 3,
})

// Mark plugin as unhealthy (not functioning)
p.SetUnhealthy("Cannot receive events", map[string]interface{}{
    "error": err.Error(),
})
```

The plugin loader will automatically monitor health and restart unhealthy plugins.

## Metrics

BasePlugin provides automatic metric tracking:

- `IncrementReceived()` - Count events received
- `IncrementSent()` - Count events sent to queue
- `IncrementDropped()` - Count events dropped
- `IncrementErrors()` - Count errors

Access metrics via `Stats()` method.

## Registration

Register your plugin in the agent initialization:

```go
// In internal/agent/agent.go
func (a *Agent) initializeComponents(ctx context.Context) error {
    // ...

    // Register plugins
    a.pluginLoader.Registry().Register("myplugin", myplugin.NewMyPlugin)

    // ...
}
```

## Configuration

### Main Config (config.yaml)

Enable your plugin:

```yaml
modules:
  - name: myplugin
    enabled: true
    config_file: "/etc/bigpanda-agent/modules/myplugin.yaml"
```

### Plugin Config (modules/myplugin.yaml)

```yaml
listen_address: "0.0.0.0:8080"
# ... your plugin-specific settings
```

## Best Practices

1. **Always use `RecoverFromPanic()`** in goroutines
2. **Check `IsStarted()`** before operations
3. **Use `SetHealth*()`** to report status
4. **Update metrics** for observability
5. **Handle context cancellation** properly
6. **Validate configuration** in `Init()`
7. **Close resources** in `Stop()`
8. **Log important events** with structured logging

## Example: Complete Plugin

See `internal/plugin/mock.go` for a complete working example.

## Built-in Plugins

- **SNMP**: SNMP trap receiver with MIB parsing
- **Webhook**: HTTP/HTTPS webhook receiver with transformation
- **Automation** (Future): Bidirectional automation task execution

## Testing Plugins

Create unit tests for your plugin:

```go
func TestMyPlugin(t *testing.T) {
    // Create plugin
    p := NewMyPlugin()

    // Mock config
    config := plugin.PluginConfig{
        Name: "myplugin",
        Enabled: true,
        Queue: queue.NewMemoryQueue(queue.Config{MaxSize: 100}),
    }

    // Test Init
    err := p.Init(context.Background(), config)
    assert.NoError(t, err)

    // Test Start
    err = p.Start(context.Background())
    assert.NoError(t, err)

    // Test functionality
    // ...

    // Test Stop
    err = p.Stop(context.Background())
    assert.NoError(t, err)
}
```

## Debugging

Enable debug logging for your plugin:

```yaml
logging:
  level: "debug"
```

Check plugin health:

```bash
curl http://localhost:8443/health
```

View plugin statistics:

```bash
curl http://localhost:8443/stats
```

## Support

- Issues: https://github.com/ReggieJTech/SuperAgent/issues
- Docs: https://docs.bigpanda.io/super-agent/plugins
