// Package diff provides utilities for processing and presenting infrastructure
// drift results.
//
// # Pipeline
//
// Pipeline composes an ordered sequence of PipelineStep functions, each of
// which accepts a []drift.DriftResult and returns a transformed slice (or an
// error). Steps are executed in registration order; the first error aborts
// the run.
//
// Typical usage:
//
//	p := diff.NewPipeline(
//		diff.NewNormalizer(opts).Normalize,
//		diff.NewFilter(fopts).Filter,
//		diff.NewReducer(ropts).Reduce,
//	)
//	out, err := p.Run(results)
package diff
