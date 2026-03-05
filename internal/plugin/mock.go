package plugin

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

// MockPlugin is a simple plugin for testing
type MockPlugin struct {
	*BasePlugin
	stopChan chan struct{}
}

// NewMockPlugin creates a new mock plugin
func NewMockPlugin() Plugin {
	base := NewBasePlugin("mock", "1.0.0", "Mock plugin for testing")
	return &MockPlugin{
		BasePlugin: base,
		stopChan:   make(chan struct{}),
	}
}

// Start starts the mock plugin
func (p *MockPlugin) Start(ctx context.Context) error {
	if p.IsStarted() {
		return nil
	}

	log.Info().Str("plugin", p.Name()).Msg("Starting mock plugin")

	p.MarkStarted()
	p.SetHealthy()

	// Simulate work
	go p.run(ctx)

	return nil
}

// Stop stops the mock plugin
func (p *MockPlugin) Stop(ctx context.Context) error {
	if !p.IsStarted() {
		return nil
	}

	log.Info().Str("plugin", p.Name()).Msg("Stopping mock plugin")

	close(p.stopChan)
	p.MarkStopped()

	return nil
}

// run simulates plugin work
func (p *MockPlugin) run(ctx context.Context) {
	defer p.RecoverFromPanic()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopChan:
			return
		case <-ticker.C:
			// Simulate receiving events
			p.IncrementReceived()
			p.IncrementSent()

			log.Debug().
				Str("plugin", p.Name()).
				Int64("received", p.eventsReceived.Load()).
				Msg("Mock plugin tick")
		}
	}
}
