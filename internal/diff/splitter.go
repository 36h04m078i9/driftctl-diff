package diff

import "github.com/nikoksr/driftctl-diff/internal/drift"

// SplitOptions controls how results are partitioned.
type SplitOptions struct {
	// ByKind splits into added, deleted, and changed buckets.
	ByKind bool
}

// SplitResult holds partitioned drift results.
type SplitResult struct {
	Added   []drift.ResourceDiff
	Deleted []drift.ResourceDiff
	Changed []drift.ResourceDiff
}

// Splitter partitions a slice of ResourceDiff into logical buckets.
type Splitter struct {
	opts SplitOptions
}

// NewSplitter returns a Splitter configured with opts.
func NewSplitter(opts SplitOptions) *Splitter {
	return &Splitter{opts: opts}
}

// Split partitions results into Added, Deleted, and Changed buckets.
// A resource is Added when every change has Kind == KindAdded,
// Deleted when every change has Kind == KindDeleted, and
// Changed otherwise (including mixed or empty-change resources).
func (s *Splitter) Split(results []drift.ResourceDiff) SplitResult {
	var out SplitResult
	for _, r := range results {
		if len(r.Changes) == 0 {
			out.Changed = append(out.Changed, r)
			continue
		}
		kind := r.Changes[0].Kind
		uniform := true
		for _, c := range r.Changes[1:] {
			if c.Kind != kind {
				uniform = false
				break
			}
		}
		if !uniform {
			out.Changed = append(out.Changed, r)
			continue
		}
		switch kind {
		case drift.KindAdded:
			out.Added = append(out.Added, r)
		case drift.KindDeleted:
			out.Deleted = append(out.Deleted, r)
		default:
			out.Changed = append(out.Changed, r)
		}
	}
	return out
}
