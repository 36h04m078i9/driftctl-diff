// Package summary aggregates drift detection results into human-readable
// statistics. It exposes Stats, which is computed from a slice of
// drift.ResourceDiff values, and Printer, which formats those statistics
// to any io.Writer.
//
// Typical usage:
//
//	results := detector.Detect(stateResources, liveResources)
//	stats := summary.Compute(results)
//	printer := summary.NewPrinter(os.Stdout)
//	printer.Print(stats)
package summary
