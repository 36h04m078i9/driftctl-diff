package diff

import "github.com/driftctl-diff/internal/drift"

// TruncateOptions controls how results are truncated.
type TruncateOptions struct {
	MaxResults    int
	MaxChanges    int
	MaxValueLen   int
}

// DefaultTruncateOptions returns sensible defaults.
func DefaultTruncateOptions() TruncateOptions {
	return TruncateOptions{
		MaxResults:  100,
		MaxChanges:  20,
		MaxValueLen: 120,
	}
}

// Truncator limits the size of drift results for display.
type Truncator struct {
	opts TruncateOptions
}

// NewTruncator creates a Truncator with the given options.
func NewTruncator(opts TruncateOptions) *Truncator {
	return &Truncator{opts: opts}
}

// Truncate applies all limits and returns the (possibly shortened) slice.
// A boolean is also returned indicating whether any truncation occurred.
func (t *Truncator) Truncate(results []drift.ResourceDiff) ([]drift.ResourceDiff, bool) {
	truncated := false

	if t.opts.MaxResults > 0 && len(results) > t.opts.MaxResults {
		results = results[:t.opts.MaxResults]
		truncated = true
	}

	out := make([]drift.ResourceDiff, 0, len(results))
	for _, r := range results {
		changes := r.Changes
		if t.opts.MaxChanges > 0 && len(changes) > t.opts.MaxChanges {
			changes = changes[:t.opts.MaxChanges]
			truncated = true
		}
		truncated = t.truncateValues(changes) || truncated
		out = append(out, drift.ResourceDiff{
			ResourceID:   r.ResourceID,
			ResourceType: r.ResourceType,
			Changes:      changes,
		})
	}
	return out, truncated
}

func (t *Truncator) truncateValues(changes []drift.AttributeChange) bool {
	if t.opts.MaxValueLen <= 0 {
		return false
	}
	truncated := false
	for i := range changes {
		if len(changes[i].Got) > t.opts.MaxValueLen {
			changes[i].Got = changes[i].Got[:t.opts.MaxValueLen] + "..."
			truncated = true
		}
		if len(changes[i].Want) > t.opts.MaxValueLen {
			changes[i].Want = changes[i].Want[:t.opts.MaxValueLen] + "..."
			truncated = true
		}
	}
	return truncated
}
