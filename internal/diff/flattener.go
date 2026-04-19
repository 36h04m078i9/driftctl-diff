package diff

import "github.com/driftctl-diff/internal/drift"

// FlatChange is a single attribute-level drift record with resource context.
type FlatChange struct {
	ResourceType string
	ResourceID   string
	Attribute    string
	Kind         drift.ChangeKind
	WantValue    string
	GotValue     string
}

// Flattener expands a slice of DriftResult into individual FlatChange records,
// one per attribute change, making it easy to feed into tabular outputs.
type Flattener struct{}

// NewFlattener returns a new Flattener.
func NewFlattener() *Flattener {
	return &Flattener{}
}

// Flatten converts []drift.DriftResult into []FlatChange.
func (f *Flattener) Flatten(results []drift.DriftResult) []FlatChange {
	var out []FlatChange
	for _, r := range results {
		for _, c := range r.Changes {
			out = append(out, FlatChange{
				ResourceType: r.ResourceType,
				ResourceID:   r.ResourceID,
				Attribute:    c.Attribute,
				Kind:         c.Kind,
				WantValue:    c.WantValue,
				GotValue:     c.GotValue,
			})
		}
	}
	return out
}

// FlattenByKind returns only the FlatChange records matching the given kind.
func (f *Flattener) FlattenByKind(results []drift.DriftResult, kind drift.ChangeKind) []FlatChange {
	all := f.Flatten(results)
	var out []FlatChange
	for _, fc := range all {
		if fc.Kind == kind {
			out = append(out, fc)
		}
	}
	return out
}
