// Package diff provides utilities for comparing, rendering, and analysing
// infrastructure drift results.
//
// The Scorer type assigns a weighted numeric score to each drifted resource
// based on the kinds of attribute changes detected:
//
//   - KindMissing  → weight 3 (Critical)
//   - KindChanged  → weight 2 (Warning)
//   - KindAdded    → weight 1 (Info)
//
// Results are returned sorted in descending order of total score so that the
// most severely drifted resources appear first.
//
// ScorerPrinter renders the scores as a tab-aligned table suitable for
// terminal output.
package diff
