// Package diff provides utilities for comparing, rendering, and analysing
// infrastructure drift results.
//
// # Digest
//
// The Digest type produces a stable, order-independent SHA-256 fingerprint
// for a set of DriftResult values. It is designed for use in CI pipelines
// where you need to detect whether the drift profile has changed between two
// runs without performing a full structural comparison.
//
// Usage:
//
//	opts := diff.DefaultDigestOptions()
//	d    := diff.NewDigest(opts)
//	hash := d.Compute(results)
//
// DigestPrinter wraps a Digest and writes the hash to any io.Writer in a
// human-readable one-liner format:
//
//	drift-digest: <hex>
package diff
