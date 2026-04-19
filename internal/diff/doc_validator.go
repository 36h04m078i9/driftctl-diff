// Package diff provides utilities for comparing, analysing, and rendering
// infrastructure drift results produced by driftctl-diff.
//
// # Validator
//
// The Validator checks a slice of drift.ResourceDiff values for structural
// integrity before they are passed to downstream formatters or reporters.
//
// Errors are raised for results with empty ResourceID or ResourceType fields.
// Warnings are emitted when the number of attribute changes on a single
// resource exceeds the configured threshold (default 500).
//
// Usage:
//
//	v := diff.NewValidator(200)
//	vr, err := v.Validate(results)
//	if err != nil { ... }
//	diff.NewValidatorPrinter(os.Stderr).Print(vr)
package diff
