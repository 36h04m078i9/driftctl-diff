// Package diff provides utilities for computing, rendering, and managing
// infrastructure drift results.
//
// # Archiver
//
// The Archiver and ArchiverPrinter types allow drift results to be persisted
// to disk and later reviewed.
//
// Basic usage:
//
//	archiver := diff.NewArchiver("/var/lib/driftctl/archives")
//
//	// Save current drift results with an optional label.
//	path, err := archiver.Save(results, "nightly-run")
//
//	// List all saved archives.
//	paths, err := archiver.List()
//
//	// Load a specific archive entry.
//	entry, err := archiver.Load(paths[0])
//
//	// Print a summary table of all archives.
//	printer := diff.NewArchiverPrinter(archiver, os.Stdout)
//	printer.Print()
package diff
