package diff

import "github.com/driftctl/driftctl-diff/internal/drift"

// Stats holds aggregate counts derived from a slice of DriftResults.
type Stats struct {
	TotalResources  int
	DriftedResources int
	AddedAttributes  int
	ChangedAttributes int
	DeletedAttributes int
}

// Compute calculates drift statistics from the provided results.
func Compute(results []drift.DriftResult) Stats {
	s := Stats{
		TotalResources: len(results),
	}
	for _, r := range results {
		if len(r.Changes) == 0 {
			continue
		}
		s.DriftedResources++
		for _, c := range r.Changes {
			switch c.Kind {
			case drift.Added:
				s.AddedAttributes++
			case drift.Changed:
				s.ChangedAttributes++
			case drift.Deleted:
				s.DeletedAttributes++
			}
		}
	}
	return s
}

// DriftPercent returns the percentage of resources that have drifted.
// Returns 0 if there are no resources.
func (s Stats) DriftPercent() float64 {
	if s.TotalResources == 0 {
		return 0
	}
	return float64(s.DriftedResources) / float64(s.TotalResources) * 100
}

// TotalChanges returns the total number of attribute changes across all resources.
func (s Stats) TotalChanges() int {
	return s.AddedAttributes + s.ChangedAttributes + s.DeletedAttributes
}
