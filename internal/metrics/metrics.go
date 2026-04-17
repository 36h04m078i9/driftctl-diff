// Package metrics tracks runtime statistics for a drift scan run.
package metrics

import (
	"sync"
	"time"
)

// Counters holds scan-level statistics.
type Counters struct {
	mu              sync.Mutex
	ResourcesTotal  int
	ResourcesDrifted int
	AttributesChecked int
	FetchErrors     int
	Duration        time.Duration
}

// Collector gathers metrics during a scan.
type Collector struct {
	start    time.Time
	counters Counters
}

// New creates a new Collector and starts the clock.
func New() *Collector {
	return &Collector{start: time.Now()}
}

// IncResources increments the total resource count.
func (c *Collector) IncResources(n int) {
	c.counters.mu.Lock()
	defer c.counters.mu.Unlock()
	c.counters.ResourcesTotal += n
}

// IncDrifted increments the drifted resource count.
func (c *Collector) IncDrifted(n int) {
	c.counters.mu.Lock()
	defer c.counters.mu.Unlock()
	c.counters.ResourcesDrifted += n
}

// IncAttributes increments the attribute check count.
func (c *Collector) IncAttributes(n int) {
	c.counters.mu.Lock()
	defer c.counters.mu.Unlock()
	c.counters.AttributesChecked += n
}

// IncFetchErrors increments the fetch error count.
func (c *Collector) IncFetchErrors(n int) {
	c.counters.mu.Lock()
	defer c.counters.mu.Unlock()
	c.counters.FetchErrors += n
}

// Snapshot finalises timing and returns a copy of the counters.
func (c *Collector) Snapshot() Counters {
	c.counters.mu.Lock()
	defer c.counters.mu.Unlock()
	c.counters.Duration = time.Since(c.start)
	copy := c.counters
	return copy
}
