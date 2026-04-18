package diff

import (
	"sort"

	"github.com/acme/driftctl-diff/internal/drift"
)

// CompareOptions controls how two sets of drift results are compared.
type CompareOptions struct {
	IgnoreAdded   bool
	IgnoreRemoved bool
	IgnoreChanged bool
}

// Comparator compares two snapshots of drift results and surfaces what changed
// between them (e.g. new drift introduced, drift resolved).
type Comparator struct {
	opts CompareOptions
}

// NewComparator returns a Comparator with the given options.
func NewComparator(opts CompareOptions) *Comparator {
	return &Comparator{opts: opts}
}

// CompareResult holds the outcome of comparing two drift result sets.
type CompareResult struct {
	NewDrift      []drift.ResourceDiff // drift present in current but not baseline
	ResolvedDrift []drift.ResourceDiff // drift present in baseline but not current
	Persisted     []drift.ResourceDiff // drift present in both
}

// Compare compares baseline results against current results.
func (c *Comparator) Compare(baseline, current []drift.ResourceDiff) CompareResult {
	baseMap := indexByKey(baseline)
	currMap := indexByKey(current)

	var result CompareResult

	for k, curr := range currMap {
		if _, found := baseMap[k]; !found {
			if !c.opts.IgnoreAdded {
				result.NewDrift = append(result.NewDrift, curr)
			}
		} else {
			if !c.opts.IgnoreChanged {
				result.Persisted = append(result.Persisted, curr)
			}
		}
	}

	for k, base := range baseMap {
		if _, found := currMap[k]; !found {
			if !c.opts.IgnoreRemoved {
				result.ResolvedDrift = append(result.ResolvedDrift, base)
			}
		}
	}

	sortDiffs(result.NewDrift)
	sortDiffs(result.ResolvedDrift)
	sortDiffs(result.Persisted)

	return result
}

func indexByKey(diffs []drift.ResourceDiff) map[string]drift.ResourceDiff {
	m := make(map[string]drift.ResourceDiff, len(diffs))
	for _, d := range diffs {
		m[d.ResourceType+"."+d.ResourceID] = d
	}
	return m
}

func sortDiffs(diffs []drift.ResourceDiff) {
	sort.Slice(diffs, func(i, j int) bool {
		if diffs[i].ResourceType != diffs[j].ResourceType {
			return diffs[i].ResourceType < diffs[j].ResourceType
		}
		return diffs[i].ResourceID < diffs[j].ResourceID
	})
}
