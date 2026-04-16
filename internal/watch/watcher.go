// Package watch provides periodic drift re-detection for continuous monitoring.
package watch

import (
	"context"
	"time"

	"github.com/user/driftctl-diff/internal/drift"
	"github.com/user/driftctl-diff/internal/runner"
)

// Result holds the outcome of a single watch tick.
type Result struct {
	At      time.Time
	Changes []drift.Change
	Err     error
}

// Watcher repeatedly runs drift detection on a fixed interval.
type Watcher struct {
	runner   *runner.Runner
	interval time.Duration
	results  chan Result
}

// New creates a Watcher that emits Results on the returned channel.
func New(r *runner.Runner, interval time.Duration) *Watcher {
	return &Watcher{
		runner:   r,
		interval: interval,
		results:  make(chan Result, 1),
	}
}

// Results returns the read-only channel of drift results.
func (w *Watcher) Results() <-chan Result { return w.results }

// Start begins the watch loop and blocks until ctx is cancelled.
func (w *Watcher) Start(ctx context.Context) {
	defer close(w.results)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			changes, err := w.runner.Detect(ctx)
			select {
			case w.results <- Result{At: t, Changes: changes, Err: err}:
			case <-ctx.Done():
				return
			}
		}
	}
}
