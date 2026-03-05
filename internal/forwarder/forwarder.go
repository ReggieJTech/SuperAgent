package forwarder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ReggieJTech/SuperAgent/internal/queue"
	"github.com/rs/zerolog/log"
)

// Forwarder sends events to BigPanda API
type Forwarder struct {
	config    Config
	queue     queue.Queue
	client    *http.Client
	batcher   *Batcher
	ratelimit *RateLimiter
	circuit   *CircuitBreaker
	mu        sync.RWMutex
	stopped   atomic.Bool

	// Metrics
	sent      atomic.Int64
	failed    atomic.Int64
	retried   atomic.Int64
	batches   atomic.Int64
	startTime time.Time
}

// Config contains forwarder configuration
type Config struct {
	APIURL       string
	Token        string
	AppKey       string
	Batching     BatchConfig
	Retry        RetryConfig
	RateLimit    RateLimitConfig
	Timeout      TimeoutConfig
}

// BatchConfig contains batching configuration
type BatchConfig struct {
	Enabled  bool
	MaxSize  int
	MaxWait  time.Duration
	MaxBytes int
}

// RetryConfig contains retry configuration
type RetryConfig struct {
	MaxAttempts       int
	InitialBackoff    time.Duration
	MaxBackoff        time.Duration
	BackoffMultiplier float64
}

// RateLimitConfig contains rate limiting configuration
type RateLimitConfig struct {
	EventsPerSecond int
	Burst           int
}

// TimeoutConfig contains timeout configuration
type TimeoutConfig struct {
	Connect time.Duration
	Request time.Duration
	Idle    time.Duration
}

// New creates a new forwarder
func New(config Config, q queue.Queue) *Forwarder {
	// Create HTTP client with timeouts
	client := &http.Client{
		Timeout: config.Timeout.Request,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     config.Timeout.Idle,
			DisableCompression:  false,
			DisableKeepAlives:   false,
		},
	}

	f := &Forwarder{
		config:    config,
		queue:     q,
		client:    client,
		startTime: time.Now(),
	}

	// Initialize batcher
	if config.Batching.Enabled {
		f.batcher = NewBatcher(config.Batching)
	}

	// Initialize rate limiter
	if config.RateLimit.EventsPerSecond > 0 {
		f.ratelimit = NewRateLimiter(
			config.RateLimit.EventsPerSecond,
			config.RateLimit.Burst,
		)
	}

	// Initialize circuit breaker
	f.circuit = NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:  5,
		ResetTimeout: 60 * time.Second,
	})

	return f
}

// Run starts the forwarder
func (f *Forwarder) Run(ctx context.Context, stopChan <-chan struct{}) {
	log.Info().Msg("Forwarder started")
	defer log.Info().Msg("Forwarder stopped")

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			// Flush any remaining batched events
			f.flush(ctx)
			return

		case <-ticker.C:
			// Process events from queue
			f.processQueue(ctx)
		}
	}
}

// processQueue processes events from the queue
func (f *Forwarder) processQueue(ctx context.Context) {
	// Check circuit breaker
	if !f.circuit.CanAttempt() {
		log.Warn().Msg("Circuit breaker is open, skipping send")
		return
	}

	// Dequeue event
	event, err := f.queue.Dequeue(ctx)
	if err != nil {
		if err != queue.ErrQueueEmpty {
			log.Error().Err(err).Msg("Failed to dequeue event")
		}
		return
	}

	// Check rate limit
	if f.ratelimit != nil {
		if !f.ratelimit.Allow() {
			// Put event back in queue
			if err := f.queue.Enqueue(ctx, event); err != nil {
				log.Error().Err(err).Msg("Failed to re-enqueue rate-limited event")
			}
			return
		}
	}

	// Send event (with batching if enabled)
	if f.config.Batching.Enabled {
		// Add to batch
		f.batcher.Add(event)

		// Check if batch is ready
		if f.batcher.ShouldFlush() {
			batch := f.batcher.Flush()
			f.sendBatch(ctx, batch)
		}
	} else {
		// Send immediately
		f.sendEvent(ctx, event)
	}
}

// sendEvent sends a single event
func (f *Forwarder) sendEvent(ctx context.Context, event *queue.Event) {
	events := []*queue.Event{event}
	f.sendBatch(ctx, events)
}

// sendBatch sends a batch of events
func (f *Forwarder) sendBatch(ctx context.Context, events []*queue.Event) {
	if len(events) == 0 {
		return
	}

	log.Debug().Int("count", len(events)).Msg("Sending batch to BigPanda")

	// Convert events to JSON
	payload, err := f.buildPayload(events)
	if err != nil {
		log.Error().Err(err).Msg("Failed to build payload")
		f.failed.Add(int64(len(events)))
		return
	}

	// Send with retry
	err = f.sendWithRetry(ctx, payload)
	if err != nil {
		log.Error().Err(err).Int("count", len(events)).Msg("Failed to send batch")
		f.failed.Add(int64(len(events)))
		f.circuit.RecordFailure()

		// Send failed events to DLQ
		for _, event := range events {
			if q, ok := f.queue.(*queue.BadgerQueue); ok {
				if dlqErr := q.SendToDLQ(event, err.Error()); dlqErr != nil {
					log.Error().Err(dlqErr).Msg("Failed to send to DLQ")
				}
			}
		}
	} else {
		f.sent.Add(int64(len(events)))
		f.batches.Add(1)
		f.circuit.RecordSuccess()
		log.Info().Int("count", len(events)).Msg("Batch sent successfully")
	}
}

// sendWithRetry sends a payload with exponential backoff retry
func (f *Forwarder) sendWithRetry(ctx context.Context, payload []byte) error {
	var lastErr error
	backoff := f.config.Retry.InitialBackoff

	for attempt := 1; attempt <= f.config.Retry.MaxAttempts; attempt++ {
		// Create request
		req, err := http.NewRequestWithContext(ctx, "POST", f.config.APIURL, bytes.NewReader(payload))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+f.config.Token)
		req.Header.Set("X-App-Key", f.config.AppKey)
		req.Header.Set("User-Agent", "BigPanda-Super-Agent/1.0")

		// Send request
		resp, err := f.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			log.Warn().
				Err(err).
				Int("attempt", attempt).
				Dur("backoff", backoff).
				Msg("Request failed, retrying")

			// Retry on network errors
			if attempt < f.config.Retry.MaxAttempts {
				f.retried.Add(1)
				time.Sleep(backoff)
				backoff = f.calculateBackoff(backoff)
				continue
			}
			return lastErr
		}

		// Read response body
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		// Check status code
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Success
			return nil
		}

		// Handle different error codes
		switch resp.StatusCode {
		case 400:
			// Bad request - don't retry
			return fmt.Errorf("bad request (400): %s", string(body))

		case 401, 403:
			// Authentication error - don't retry
			return fmt.Errorf("authentication error (%d): %s", resp.StatusCode, string(body))

		case 429:
			// Rate limited - retry with longer backoff
			lastErr = fmt.Errorf("rate limited (429)")
			log.Warn().Int("attempt", attempt).Msg("Rate limited, backing off")
			time.Sleep(backoff * 2)
			backoff = f.calculateBackoff(backoff * 2)

		case 500, 502, 503, 504:
			// Server error - retry
			lastErr = fmt.Errorf("server error (%d): %s", resp.StatusCode, string(body))
			log.Warn().
				Int("status", resp.StatusCode).
				Int("attempt", attempt).
				Msg("Server error, retrying")

			if attempt < f.config.Retry.MaxAttempts {
				f.retried.Add(1)
				time.Sleep(backoff)
				backoff = f.calculateBackoff(backoff)
				continue
			}

		default:
			// Unknown error
			lastErr = fmt.Errorf("unexpected status (%d): %s", resp.StatusCode, string(body))
			log.Warn().
				Int("status", resp.StatusCode).
				Int("attempt", attempt).
				Msg("Unexpected status, retrying")

			if attempt < f.config.Retry.MaxAttempts {
				f.retried.Add(1)
				time.Sleep(backoff)
				backoff = f.calculateBackoff(backoff)
				continue
			}
		}

		return lastErr
	}

	return lastErr
}

// buildPayload builds a JSON payload from events
func (f *Forwarder) buildPayload(events []*queue.Event) ([]byte, error) {
	// Convert events to JSON objects
	jsonEvents := make([]map[string]interface{}, 0, len(events))

	for _, event := range events {
		data := make(map[string]interface{})

		// Core fields
		data["status"] = event.Status
		data["primary_property"] = event.PrimaryKey
		data["secondary_property"] = event.SecondaryKey
		data["timestamp"] = event.Timestamp.Unix()

		// Optional fields
		if event.Check != "" {
			data["check"] = event.Check
		}
		if event.Description != "" {
			data["description"] = event.Description
		}
		if event.SourceSystem != "" {
			data["source_system"] = event.SourceSystem
		}

		// Tags
		for k, v := range event.Tags {
			data[k] = v
		}

		// Custom fields
		for k, v := range event.CustomFields {
			data[k] = v
		}

		jsonEvents = append(jsonEvents, data)
	}

	return json.Marshal(jsonEvents)
}

// calculateBackoff calculates the next backoff duration
func (f *Forwarder) calculateBackoff(current time.Duration) time.Duration {
	next := time.Duration(float64(current) * f.config.Retry.BackoffMultiplier)
	if next > f.config.Retry.MaxBackoff {
		return f.config.Retry.MaxBackoff
	}
	return next
}

// flush flushes any remaining batched events
func (f *Forwarder) flush(ctx context.Context) {
	if f.batcher == nil {
		return
	}

	batch := f.batcher.Flush()
	if len(batch) > 0 {
		log.Info().Int("count", len(batch)).Msg("Flushing remaining events")
		f.sendBatch(ctx, batch)
	}
}

// Stop stops the forwarder
func (f *Forwarder) Stop(ctx context.Context) {
	if f.stopped.Swap(true) {
		return
	}

	log.Info().Msg("Stopping forwarder")
	f.flush(ctx)
}

// Health returns health information
func (f *Forwarder) Health() map[string]interface{} {
	return map[string]interface{}{
		"status":         "healthy",
		"circuit_breaker": f.circuit.State(),
		"stopped":        f.stopped.Load(),
	}
}

// Stats returns forwarder statistics
func (f *Forwarder) Stats() map[string]interface{} {
	stats := map[string]interface{}{
		"sent":     f.sent.Load(),
		"failed":   f.failed.Load(),
		"retried":  f.retried.Load(),
		"batches":  f.batches.Load(),
		"uptime":   time.Since(f.startTime).String(),
	}

	if f.batcher != nil {
		stats["batch_size"] = f.batcher.Size()
	}

	return stats
}
