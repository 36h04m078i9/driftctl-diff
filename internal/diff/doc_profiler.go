// Package diff provides utilities for computing, rendering, and analysing
// infrastructure drift between Terraform state and live cloud resources.
//
// # Profiler
//
// The Profiler type collects lightweight performance profiles across one or
// more diff operations. Each call to Record captures the operation label,
// elapsed wall-clock duration, resource count, and total attribute-change
// count for a slice of ResourceDiff values.
//
// Entries are returned sorted by duration descending so that the slowest
// operations surface first, making it easy to identify bottlenecks during
// large-scale drift scans.
//
// Example usage:
//
//	p := diff.NewProfiler()
//	start := time.Now()
//	results := detector.Detect(stateResources, liveResources)
//	p.Record("detect", results, time.Since(start))
//	p.Print(os.Stdout)
package diff
