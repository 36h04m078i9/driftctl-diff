package diff

import (
	"fmt"
	"strings"

	"github.com/acme/driftctl-diff/internal/drift"
)

// Label represents a human-readable tag attached to a drift result.
type Label struct {
	Key   string
	Value string
}

// LabeledResult wraps a DriftResult with additional labels.
type LabeledResult struct {
	drift.DriftResult
	Labels []Label
}

// Labeler attaches metadata labels to drift results.
type Labeler struct {
	rules []labelRule
}

type labelRule struct {
	key     string
	matchFn func(drift.DriftResult) (string, bool)
}

// NewLabeler constructs a Labeler with default rules.
func NewLabeler() *Labeler {
	l := &Labeler{}
	l.rules = []labelRule{
		{
			key: "resource_type",
			matchFn: func(r drift.DriftResult) (string, bool) {
				return r.ResourceType, r.ResourceType != ""
			},
		},
		{
			key: "drift_kind",
			matchFn: func(r drift.DriftResult) (string, bool) {
				if len(r.Changes) == 0 {
					return "", false
				}
				kinds := make(map[string]struct{})
				for _, c := range r.Changes {
					kinds[string(c.Kind)] = struct{}{}
				}
				parts := make([]string, 0, len(kinds))
				for k := range kinds {
					parts = append(parts, k)
				}
				return strings.Join(parts, ","), true
			},
		},
		{
			key: "change_count",
			matchFn: func(r drift.DriftResult) (string, bool) {
				return fmt.Sprintf("%d", len(r.Changes)), true
			},
		},
	}
	return l
}

// Label applies all rules to each result and returns LabeledResults.
func (l *Labeler) Label(results []drift.DriftResult) []LabeledResult {
	out := make([]LabeledResult, 0, len(results))
	for _, r := range results {
		lr := LabeledResult{DriftResult: r}
		for _, rule := range l.rules {
			if v, ok := rule.matchFn(r); ok {
				lr.Labels = append(lr.Labels, Label{Key: rule.key, Value: v})
			}
		}
		out = append(out, lr)
	}
	return out
}
