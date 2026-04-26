package diff

import (
	"sort"

	"github.com/acme/driftctl-diff/internal/drift"
)

// Score holds the drift score for a single resource.
type Score struct {
	ResourceType string
	ResourceID   string
	Total        int
	Changed      int
	Missing      int
}

// ScorerOptions configures scoring behaviour.
type ScorerOptions struct {
	// WeightChanged is the multiplier applied to changed-attribute count.
	WeightChanged int
	// WeightMissing is the multiplier applied to missing-attribute count.
	WeightMissing int
}

// DefaultScorerOptions returns sensible defaults.
func DefaultScorerOptions() ScorerOptions {
	return ScorerOptions{WeightChanged: 1, WeightMissing: 2}
}

// Scorer computes a numeric drift score per resource.
type Scorer struct {
	opts ScorerOptions
}

// NewScorer constructs a Scorer with the given options.
func NewScorer(opts ScorerOptions) *Scorer {
	return &Scorer{opts: opts}
}

// Score computes scores for all results, sorted descending by Total.
func (s *Scorer) Score(results []drift.Result) []Score {
	scores := make([]Score, 0, len(results))
	for _, r := range results {
		var changed, missing int
		for _, c := range r.Changes {
			switch c.Kind {
			case drift.KindChanged:
				changed++
			case drift.KindMissing:
				missing++
			}
		}
		total := changed*s.opts.WeightChanged + missing*s.opts.WeightMissing
		scores = append(scores, Score{
			ResourceType: r.ResourceType,
			ResourceID:   r.ResourceID,
			Total:        total,
			Changed:      changed,
			Missing:      missing,
		})
	}
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Total > scores[j].Total
	})
	return scores
}
