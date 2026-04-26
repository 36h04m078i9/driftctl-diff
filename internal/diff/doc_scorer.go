// Package diff provides utilities for working with drift results.
//
// # Scorer
//
// The Scorer assigns a numeric weight to each drifted resource based on the
// number and kind of attribute changes detected.
//
// Changed attributes are weighted by ScorerOptions.WeightChanged (default 1).
// Missing attributes are weighted by ScorerOptions.WeightMissing (default 2)
// because a missing attribute typically represents a more severe deviation
// from the desired state.
//
// Results are returned sorted in descending order of total score, making it
// easy to identify the most severely drifted resources at a glance.
//
// Usage:
//
//	scorer := diff.NewScorer(diff.DefaultScorerOptions())
//	scores := scorer.Score(results)
//	diff.NewScorerPrinter(os.Stdout).Print(scores)
package diff
