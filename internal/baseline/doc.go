// Package baseline manages a persisted set of acknowledged drift entries.
//
// A baseline allows operators to mark known infrastructure drift as
// intentional, suppressing it from future diff output until it changes again.
//
// Usage:
//
//	b, err := baseline.LoadFrom(".driftctl-baseline.json")
//	if err != nil { /* handle */ }
//	filtered := b.Filter(driftResults)
//
// To acknowledge new drift:
//
//	b.Add(resourceType, resourceID, attribute)
//	_ = b.SaveTo(".driftctl-baseline.json")
package baseline
