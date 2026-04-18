package diff

import (
	"sort"

	"github.com/snyk/driftctl-diff/internal/drift"
)

// Score holds a drift severity score for a single resource.
type Score struct {
	ResourceType string
	ResourceID   string
	Total        int
	Critical     int
	Warning      int
	Info         int
}

// Scorer computes a weighted drift score for each resource.
type Scorer struct {
	weights map[drift.ChangeKind]int
}

// NewScorer returns a Scorer with default weights.
func NewScorer() *Scorer {
	return &Scorer{
		weights: map[drift.ChangeKind]int{
			drift.KindChanged: 2,
			drift.KindMissing: 3,
			drift.KindAdded:   1,
		},
	}
}

// Score computes a Score for each DriftResult, sorted descending by Total.
func (s *Scorer) Score(results []drift.DriftResult) []Score {
	scores := make([]Score, 0, len(results))
	for _, r := range results {
		sc := Score{
			ResourceType: r.ResourceType,
			ResourceID:   r.ResourceID,
		}
		for _, c := range r.Changes {
			w := s.weights[c.Kind]
			sc.Total += w
			switch {
			case w >= 3:
				sc.Critical++
			case w == 2:
				sc.Warning++
			default:
				sc.Info++
			}
		}
		scores = append(scores, sc)
	}
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Total > scores[j].Total
	})
	return scores
}
