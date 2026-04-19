package diff

import (
	"sort"

	"github.com/acme/driftctl-diff/internal/drift"
)

// PivotEntry represents a single attribute pivoted across multiple resources.
type PivotEntry struct {
	Attribute string
	Resources []PivotResource
}

// PivotResource holds per-resource values for a pivoted attribute.
type PivotResource struct {
	ResourceID   string
	ResourceType string
	StateValue   string
	LiveValue    string
}

// Pivot reorganises drift results so that rows are attributes and columns are
// resources — the inverse of the default resource-centric view.
type Pivot struct{}

// NewPivot returns a new Pivot.
func NewPivot() *Pivot { return &Pivot{} }

// Compute builds a slice of PivotEntry from drift results, one entry per
// unique attribute name that has at least one change.
func (p *Pivot) Compute(results []drift.ResourceDiff) []PivotEntry {
	index := map[string]*PivotEntry{}
	var order []string

	for _, r := range results {
		for _, c := range r.Changes {
			if _, ok := index[c.Attribute]; !ok {
				index[c.Attribute] = &PivotEntry{Attribute: c.Attribute}
				order = append(order, c.Attribute)
			}
			index[c.Attribute].Resources = append(index[c.Attribute].Resources, PivotResource{
				ResourceID:   r.ResourceID,
				ResourceType: r.ResourceType,
				StateValue:   c.StateValue,
				LiveValue:    c.LiveValue,
			})
		}
	}

	sort.Strings(order)
	out := make([]PivotEntry, 0, len(order))
	for _, attr := range order {
		out = append(out, *index[attr])
	}
	return out
}
