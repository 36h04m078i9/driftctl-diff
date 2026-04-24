// Package diff provides utilities for processing and presenting infrastructure
// drift results produced by the detector.
//
// The Labeler sub-feature attaches structured key/value labels to each
// DriftResult, making it easy for downstream consumers (formatters, exporters,
// policy engines) to filter or annotate output without re-inspecting raw
// change slices.
//
// # Labeler
//
// NewLabeler creates a Labeler with default label rules. Labels are derived
// from the DriftResult fields according to the following conventions:
//
//   - "resource_type": the type of the drifted resource (e.g. "aws_s3_bucket")
//   - "drift_kind":    one of "missing", "unmanaged", or "changed"
//   - "severity":      inferred from drift_kind ("high" for missing/unmanaged,
//     "low" for changed)
//
// Usage:
//
//	labeler := diff.NewLabeler()
//	labeled := labeler.Label(results)
//	for _, lr := range labeled {
//		fmt.Println(lr.ResourceID, lr.Labels)
//	}
package diff
