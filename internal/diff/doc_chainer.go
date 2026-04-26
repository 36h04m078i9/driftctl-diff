// Package diff provides utilities for processing, transforming, and rendering
// infrastructure drift results produced by driftctl-diff.
//
// # Chainer
//
// The Chainer type executes an ordered pipeline of named transformation steps
// over a slice of drift.DriftResult values. Each step is a pure function that
// receives the current slice and returns a (possibly filtered or modified)
// replacement slice.
//
// Steps are composable: you can assemble pipelines from the existing
// Filter, Truncator, Normalizer, Reducer, Pruner, and other primitives in this
// package by wrapping their methods as ChainStep.ApplyFn values.
//
// The ChainResult returned by Run contains both the final results and a
// per-step count summary, which can be rendered with ChainerPrinter for
// debugging or audit purposes.
package diff
