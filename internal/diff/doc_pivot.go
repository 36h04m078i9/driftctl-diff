// Package diff provides utilities for comparing, rendering and analysing
// infrastructure drift results.
//
// # Pivot
//
// The Pivot type inverts the standard resource-centric view of drift results
// so that each row represents an attribute and each column a resource. This
// makes it easy to spot when the same attribute (e.g. "region") has drifted
// across many resources at once.
//
// Usage:
//
//	p := diff.NewPivot()
//	entries := p.Compute(results)
//	diff.NewPivotPrinter(os.Stdout).Print(entries)
package diff
