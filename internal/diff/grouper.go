package diff

import "github.com/driftctl/driftctl-diff/internal/drift"

// Group holds drift results sharing the same resource type.
type Group struct {
	ResourceType string
	Results      []drift.Result
}

// Grouper groups drift results by resource type.
type Grouper struct{}

// NewGrouper returns a new Grouper.
func NewGrouper() *Grouper {
	return &Grouper{}
}

// GroupByType partitions results into groups keyed by ResourceType.
// Results within each group preserve their original order.
func (g *Grouper) GroupByType(results []drift.Result) []Group {
	index := make(map[string]int)
	groups := []Group{}

	for _, r := range results {
		if i, ok := index[r.ResourceType]; ok {
			groups[i].Results = append(groups[i].Results, r)
		} else {
			index[r.ResourceType] = len(groups)
			groups = append(groups, Group{
				ResourceType: r.ResourceType,
				Results:      []drift.Result{r},
			})
		}
	}

	return groups
}

// GroupByKind partitions results by the first change kind found in each result.
// Results with no changes are placed under the "none" key.
func (g *Grouper) GroupByKind(results []drift.Result) map[string][]drift.Result {
	out := make(map[string][]drift.Result)

	for _, r := range results {
		key := "none"
		if len(r.Changes) > 0 {
			key = string(r.Changes[0].Kind)
		}
		out[key] = append(out[key], r)
	}

	return out
}
