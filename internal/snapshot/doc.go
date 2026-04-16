// Package snapshot provides types and helpers for capturing a point-in-time
// view of live cloud resource attributes.
//
// A Snapshot can be built from live infrastructure via a Collector, persisted
// to disk with SaveTo, and reloaded with LoadFrom. Snapshots are used by the
// drift detector to compare against Terraform state without requiring a live
// cloud API call on every run.
//
// Typical usage:
//
//	collector := snapshot.NewCollector(registry)
//	snap, err := collector.Collect(resources)
//	if err != nil { ... }
//	if err := snap.SaveTo("snapshot.json"); err != nil { ... }
package snapshot
