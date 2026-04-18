// Package diff provides utilities for processing and presenting infrastructure
// drift results produced by the detector.
//
// The Labeler sub-feature attaches structured key/value labels to each
// DriftResult, making it easy for downstream consumers (formatters, exporters,
// policy engines) to filter or annotate output without re-inspecting raw
// change slices.
//
// Usage:
//
//	labeler := diff.NewLabeler()
//	labeled := labeler.Label(results)
//	for _, lr := range labeled {
//		fmt.Println(lr.ResourceID, lr.Labels)
//	}
package diff
