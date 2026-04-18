// Package diff provides utilities for comparing, rendering, and tracking
// infrastructure drift results over time.
//
// The Timeline type records successive drift snapshots so operators can
// observe how drift evolves across multiple runs of driftctl-diff.
//
// Basic usage:
//
//	tl := diff.NewTimeline()
//	tl.Add(time.Now(), results)
//
//	printer := diff.NewTimelinePrinter(os.Stdout)
//	printer.Print(tl)
package diff
