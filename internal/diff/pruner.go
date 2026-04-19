package diff

import "github.com/driftctl-diff/internal/drift"

// PruneOptions controls which changes are removed by the Pruner.
type PruneOptions struct {
	// RemoveUnchanged drops results with no attribute changes.
	RemoveUnchanged bool
	// MinChanges drops results with fewer than this many attribute changes.
	MinChanges int
	// ExcludeTypes is a set of resource types to drop entirely.
	ExcludeTypes map[string]bool
}

// DefaultPruneOptions returns a PruneOptions with safe defaults.
func DefaultPruneOptions() PruneOptions {
	return PruneOptions{
		RemoveUnchanged: true,
		MinChanges:      0,
		ExcludeTypes:    map[string]bool{},
	}
}

// Pruner removes noise from drift results based on configurable options.
type Pruner struct {
	opts PruneOptions
}

// NewPruner creates a Pruner with the given options.
func NewPruner(opts PruneOptions) *Pruner {
	return &Pruner{opts: opts}
}

// Prune filters out results that do not meet the configured thresholds.
func (p *Pruner) Prune(results []drift.ResourceDiff) []drift.ResourceDiff {
	out := make([]drift.ResourceDiff, 0, len(results))
	for _, r := range results {
		if p.opts.ExcludeTypes[r.ResourceType] {
			continue
		}
		if p.opts.RemoveUnchanged && len(r.Changes) == 0 {
			continue
		}
		if len(r.Changes) < p.opts.MinChanges {
			continue
		}
		out = append(out, r)
	}
	return out
}
