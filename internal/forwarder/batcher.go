package forwarder

import (
	"sync"
	"time"

	"github.com/ReggieJTech/SuperAgent/internal/queue"
)

// Batcher batches events together
type Batcher struct {
	config    BatchConfig
	events    []*queue.Event
	totalSize int
	lastFlush time.Time
	mu        sync.Mutex
}

// NewBatcher creates a new batcher
func NewBatcher(config BatchConfig) *Batcher {
	return &Batcher{
		config:    config,
		events:    make([]*queue.Event, 0, config.MaxSize),
		lastFlush: time.Now(),
	}
}

// Add adds an event to the batch
func (b *Batcher) Add(event *queue.Event) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.events = append(b.events, event)

	// Estimate size (rough approximation)
	// Each event is roughly 500 bytes
	b.totalSize += 500
}

// ShouldFlush returns true if the batch should be flushed
func (b *Batcher) ShouldFlush() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Check if batch is full by count
	if len(b.events) >= b.config.MaxSize {
		return true
	}

	// Check if batch is full by size
	if b.totalSize >= b.config.MaxBytes {
		return true
	}

	// Check if max wait time has elapsed
	if time.Since(b.lastFlush) >= b.config.MaxWait {
		return true
	}

	return false
}

// Flush returns all events and resets the batch
func (b *Batcher) Flush() []*queue.Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	events := b.events
	b.events = make([]*queue.Event, 0, b.config.MaxSize)
	b.totalSize = 0
	b.lastFlush = time.Now()

	return events
}

// Size returns the current batch size
func (b *Batcher) Size() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.events)
}
