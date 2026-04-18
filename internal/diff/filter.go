package diff

import (
	"strings"

	"github.com/acme/driftctl-diff/internal/drift"
)

// FilterOptions controls which drift results are included.
type FilterOptions struct {
	ResourceType string
	ResourceID   string
	Kinds        []drift.ChangeKind
	MinChanges   int
}

// Filter returns a subset of results matching the given options.
type Filter struct {
	opts FilterOptions
}

// NewFilter creates a Filter with the provided options.
func NewFilter(opts FilterOptions) *Filter {
	return &Filter{opts: opts}
}

// Apply filters the slice of DriftResult according to the options.
func (f *Filter) Apply(results []drift.DriftResult) []drift.DriftResult {
	out := make([]drift.DriftResult, 0, len(results))
	for _, r := range results {
		if f.opts.ResourceType != "" &&
			!strings.EqualFold(r.ResourceType, f.opts.ResourceType) {
			continue
		}
		if f.opts.ResourceID != "" &&
			!strings.Contains(r.ResourceID, f.opts.ResourceID) {
			continue
		}
		if len(f.opts.Kinds) > 0 && !f.anyKindMatch(r) {
			continue
		}
		if f.opts.MinChanges > 0 && len(r.Changes) < f.opts.MinChanges {
			continue
		}
		out = append(out, r)
	}
	return out
}

func (f *Filter) anyKindMatch(r drift.DriftResult) bool {
	for _, c := range r.Changes {
		for _, k := range f.opts.Kinds {
			if c.Kind == k {
				return true
			}
		}
	}
	return false
}
