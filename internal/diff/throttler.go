package diff

import (
	"fmt"
	"time"

	"github.com/acme/driftctl-diff/internal/drift"
)

// ThrottleOptions controls the rate-limiting behaviour of the Throttler.
type ThrottleOptions struct {
	// MaxPerSecond is the maximum number of DriftResults allowed per second.
	// Zero means no limit.
	MaxPerSecond int
	// BurstSize is the number of results that can be processed before the
	// per-second limit is enforced. Defaults to MaxPerSecond when zero.
	BurstSize int
}

// DefaultThrottleOptions returns a ThrottleOptions with no limits applied.
func DefaultThrottleOptions() ThrottleOptions {
	return ThrottleOptions{
		MaxPerSecond: 0,
		BurstSize:    0,
	}
}

// Throttler limits how many DriftResults are processed per second,
// introducing a sleep between batches when the rate is exceeded.
type Throttler struct {
	opts ThrottleOptions
	sleep func(time.Duration)
}

// NewThrottler creates a Throttler with the given options.
func NewThrottler(opts ThrottleOptions) *Throttler {
	return &Throttler{
		opts:  opts,
		sleep: time.Sleep,
	}
}

// Throttle returns a copy of results, pausing between batches when the
// configured rate limit is exceeded. If MaxPerSecond is zero, all results
// are returned immediately.
func (t *Throttler) Throttle(results []drift.DriftResult) ([]drift.DriftResult, error) {
	if len(results) == 0 {
		return results, nil
	}
	if t.opts.MaxPerSecond <= 0 {
		return results, nil
	}

	burst := t.opts.BurstSize
	if burst <= 0 {
		burst = t.opts.MaxPerSecond
	}
	if burst <= 0 {
		return nil, fmt.Errorf("throttler: burst size must be positive")
	}

	interval := time.Second / time.Duration(t.opts.MaxPerSecond)
	out := make([]drift.DriftResult, 0, len(results))

	for i, r := range results {
		out = append(out, r)
		if i > 0 && i%burst == 0 {
			t.sleep(interval)
		}
	}
	return out, nil
}
