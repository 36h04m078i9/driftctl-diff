package diff

import (
	"sort"

	"github.com/acme/driftctl-diff/internal/drift"
)

// SortField defines the field to sort drift results by.
type SortField int

const (
	SortByResourceType SortField = iota
	SortByResourceID
	SortByChangeCount
)

// SortOrder defines ascending or descending order.
type SortOrder int

const (
	Ascending SortOrder = iota
	Descending
)

// Sorter sorts a slice of DriftResult values.
type Sorter struct {
	Field SortField
	Order SortOrder
}

// NewSorter returns a Sorter with the given field and order.
func NewSorter(field SortField, order SortOrder) *Sorter {
	return &Sorter{Field: field, Order: order}
}

// Sort returns a sorted copy of the provided results.
func (s *Sorter) Sort(results []drift.DriftResult) []drift.DriftResult {
	out := make([]drift.DriftResult, len(results))
	copy(out, results)

	sort.SliceStable(out, func(i, j int) bool {
		var less bool
		switch s.Field {
		case SortByResourceType:
			less = out[i].ResourceType < out[j].ResourceType
		case SortByResourceID:
			less = out[i].ResourceID < out[j].ResourceID
		case SortByChangeCount:
			less = len(out[i].Changes) < len(out[j].Changes)
		default:
			less = out[i].ResourceID < out[j].ResourceID
		}
		if s.Order == Descending {
			return !less
		}
		return less
	})

	return out
}
