// Package diff provides utilities for processing, analysing, and rendering
// infrastructure drift results produced by the drift detector.
//
// # Transformer
//
// The Transformer type applies deterministic, non-destructive mutations to a
// slice of drift.DriftResult values before they reach downstream consumers
// such as formatters or exporters.
//
// Available mutations (controlled via TransformOptions):
//
//   - PrefixResourceIDs  – prepend a fixed string to every resource ID.
//   - SuffixResourceTypes – append a fixed string to every resource type.
//   - UpperCaseAttributes – convert attribute keys to UPPER_CASE.
//   - DropEmptyValues     – remove attribute changes where both the state
//     value and the live value are empty strings.
//
// The original input slice is never modified; Transform always returns a
// freshly allocated slice.
package diff
