package diff

import "github.com/driftctl-diff/internal/drift"

// CollapserOptions controls how results are collapsed.
type CollapserOptions struct {
	// MaxChangesPerResource caps the number of attribute changes kept per resource.
	// Zero means unlimited.
	MaxChangesPerResource int
	// OmitUnchanged drops resources that have no attribute changes.
	OmitUnchanged bool
}

// DefaultCollapserOptions returns sensible defaults (no limits, keep unchanged).
func DefaultCollapserOptions() CollapserOptions {
	return CollapserOptions{
		MaxChangesPerResource: 0,
		OmitUnchanged:         false,
	}
}

// Collapser reduces the verbosity of drift results by applying per-resource
// change caps and optionally dropping resources with no changes.
type Collapser struct {
	opts CollapserOptions
}

// NewCollapser creates a Collapser with the given options.
func NewCollapser(opts CollapserOptions) *Collapser {
	return &Collapser{opts: opts}
}

// Collapse applies the collapser options to results and returns a new slice.
func (c *Collapser) Collapse(results []drift.ResourceDiff) []drift.ResourceDiff {
	out := make([]drift.ResourceDiff, 0, len(results))
	for _, r := range results {
		if c.opts.OmitUnchanged && len(r.Changes) == 0 {
			continue
		}
		changes := r.Changes
		if c.opts.MaxChangesPerResource > 0 && len(changes) > c.opts.MaxChangesPerResource {
			changes = changes[:c.opts.MaxChangesPerResource]
		}
		out = append(out, drift.ResourceDiff{
			ResourceID:   r.ResourceID,
			ResourceType: r.ResourceType,
			Changes:      changes,
			Metadata:     r.Metadata,
		})
	}
	return out
}
