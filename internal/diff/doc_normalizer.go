// Package diff provides utilities for processing, rendering, and analysing
// infrastructure drift results produced by the detector.
//
// # Normalizer
//
// The Normalizer removes cosmetic differences from a slice of DriftResult
// values before they are rendered or reported. This prevents noise caused by
// insignificant formatting variations such as leading/trailing whitespace,
// differing letter case, or unnecessary quotation marks around values.
//
// Usage:
//
//	opts := diff.NormalizeOptions{
//		TrimSpace:   true,
//		LowerCase:   false,
//		StripQuotes: true,
//	}
//	n := diff.NewNormalizer(opts)
//	clean := n.Normalize(results)
package diff
