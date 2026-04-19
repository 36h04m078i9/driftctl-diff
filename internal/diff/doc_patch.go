// Package diff provides utilities for comparing, rendering, and exporting
// infrastructure drift results.
//
// The patch sub-feature generates unified-diff style output from drift results,
// making it easy to pipe into standard patch tooling or store as artefacts.
//
// Basic usage:
//
//	opts := diff.DefaultPatchOptions()
//	p := diff.NewPatcher(opts, os.Stdout)
//	p.Generate(results)
//
// To include a run header use PatchExporter:
//
//	pe := diff.NewPatchExporter(opts, os.Stdout)
//	pe.Export(results)
package diff
