// Package diff provides utilities for processing, transforming, and presenting
// infrastructure drift results produced by the detector.
//
// # Reducer
//
// The Reducer trims a drift result set to a focused subset using configurable
// criteria:
//
//   - KeepTopN – retain only the N resources with the most attribute changes.
//     Useful for triage workflows where operators want to address the noisiest
//     resources first.
//
//   - OnlyKinds – keep only results that contain at least one change of the
//     specified ChangeKind (e.g. KindMissing, KindChanged).  Useful for
//     narrowing output to a particular class of drift.
//
// Usage:
//
//	reducer := diff.NewReducer(diff.ReduceOptions{
//	    KeepTopN:  10,
//	    OnlyKinds: []drift.ChangeKind{drift.KindChanged},
//	})
//	focused := reducer.Reduce(allResults)
package diff
