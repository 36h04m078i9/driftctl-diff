package diff

import (
	"sort"

	"github.com/acme/driftctl-diff/internal/drift"
)

// ReduceOptions controls which changes are kept after reduction.
type ReduceOptions struct {
	// KeepTopN retains only the N resources with the most changes.
	// Zero means keep all.
	KeepTopN int

	// OnlyKinds, when non-empty, keeps only results whose changes contain at
	// least one entry matching one of the specified kinds.
	OnlyKinds []drift.ChangeKind
}

// Reducer trims a drift result set down to a focused subset according to
// configurable criteria, making large runs easier to triage.
type Reducer struct {
	opts ReduceOptions
}

// NewReducer returns a Reducer configured with opts.
func NewReducer(opts ReduceOptions) *Reducer {
	return &Reducer{opts: opts}
}

// Reduce applies all configured reduction rules and returns the filtered,
// possibly re-ordered slice of DriftResult.
func (r *Reducer) Reduce(results []drift.DriftResult) []drift.DriftResult {
	out := make([]drift.DriftResult, 0, len(results))

	for _, res := range results {
		if !r.matchesKinds(res) {
			continue
		}
		out = append(out, res)
	}

	// Sort descending by number of changes so KeepTopN is deterministic.
	sort.SliceStable(out, func(i, j int) bool {
		return len(out[i].Changes) > len(out[j].Changes)
	})

	if r.opts.KeepTopN > 0 && len(out) > r.opts.KeepTopN {
		out = out[:r.opts.KeepTopN]
	}

	return out
}

// matchesKinds returns true when OnlyKinds is empty or the result contains at
// least one change whose Kind is in the allow-list.
func (r *Reducer) matchesKinds(res drift.DriftResult) bool {
	if len(r.opts.OnlyKinds) == 0 {
		return true
	}
	allowed := make(map[drift.ChangeKind]struct{}, len(r.opts.OnlyKinds))
	for _, k := range r.opts.OnlyKinds {
		allowed[k] = struct{}{}
	}
	for _, ch := range res.Changes {
		if _, ok := allowed[ch.Kind]; ok {
			return true
		}
	}
	return false
}
