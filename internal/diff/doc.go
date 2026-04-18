// Package diff provides utilities for working with drift results, including
// pagination, searching, sorting, grouping, highlighting, annotation,
// exporting, filtering, statistics, and comparison between drift snapshots.
//
// The Comparator type allows callers to compare two sets of drift results —
// for example a previously recorded baseline against the current scan — and
// surface which resources have newly drifted, which have been resolved, and
// which continue to drift unchanged.
package diff
