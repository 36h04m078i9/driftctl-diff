// Package metrics provides lightweight, thread-safe counters that track
// statistics gathered during a driftctl-diff scan run, including resource
// counts, drift counts, attribute checks, fetch errors, and elapsed time.
//
// Usage:
//
//	col := metrics.New()
//	col.IncResources(10)
//	col.IncDrifted(2)
//	snap := col.Snapshot()
//	metrics.NewPrinter(os.Stdout).Print(snap)
package metrics
