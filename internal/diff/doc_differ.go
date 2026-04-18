// Package diff provides utilities for comparing, rendering, and navigating
// infrastructure drift results produced by the detector.
//
// The Differ type renders a unified-diff style output for a single
// drift.DriftResult, showing removed (state) values prefixed with '-' and
// live values prefixed with '+', grouped by attribute name.
//
// Example usage:
//
//	opts := diff.DefaultDiffOptions()
//	opts.Color = true
//	d := diff.NewDiffer(opts, os.Stdout)
//	if err := d.Render(result); err != nil {
//		log.Fatal(err)
//	}
package diff
