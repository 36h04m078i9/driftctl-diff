package diff

import (
	"strings"

	"github.com/snyk/driftctl-diff/internal/drift"
)

// SearchFilter holds criteria for filtering drift results.
type SearchFilter struct {
	ResourceType string
	ResourceID   string
	Attribute    string
	Kind         drift.ChangeKind
	KindSet      bool
}

// Searcher filters drift results based on a SearchFilter.
type Searcher struct {
	filter SearchFilter
}

// NewSearcher creates a Searcher with the given filter.
func NewSearcher(f SearchFilter) *Searcher {
	return &Searcher{filter: f}
}

// Search returns only the DriftResults whose changes match the filter.
func (s *Searcher) Search(results []drift.DriftResult) []drift.DriftResult {
	var out []drift.DriftResult
	for _, r := range results {
		if s.filter.ResourceType != "" && !strings.EqualFold(r.ResourceType, s.filter.ResourceType) {
			continue
		}
		if s.filter.ResourceID != "" && !strings.Contains(r.ResourceID, s.filter.ResourceID) {
			continue
		}
		matched := s.matchChanges(r.Changes)
		if len(matched) == 0 {
			continue
		}
		copy := r
		copy.Changes = matched
		out = append(out, copy)
	}
	return out
}

func (s *Searcher) matchChanges(changes []drift.AttributeChange) []drift.AttributeChange {
	var out []drift.AttributeChange
	for _, c := range changes {
		if s.filter.Attribute != "" && !strings.Contains(c.Attribute, s.filter.Attribute) {
			continue
		}
		if s.filter.KindSet && c.Kind != s.filter.Kind {
			continue
		}
		out = append(out, c)
	}
	return out
}
