// Package summary provides aggregation and reporting of drift detection results.
package summary

import (
	"fmt"
	"io"

	"github.com/acme/driftctl-diff/internal/drift"
)

// Stats holds aggregated counts from a drift detection run.
type Stats struct {
	TotalResources int
	DriftedResources int
	CleanResources   int
	TotalChanges     int
}

// Compute derives Stats from a slice of drift.Change slices keyed by resource.
func Compute(results []drift.ResourceDiff) Stats {
	s := Stats{
		TotalResources: len(results),
	}
	for _, rd := range results {
		if len(rd.Changes) > 0 {
			s.DriftedResources++
			s.TotalChanges += len(rd.Changes)
		} else {
			s.CleanResources++
		}
	}
	return s
}

// Printer writes a human-readable summary to an io.Writer.
type Printer struct {
	w io.Writer
}

// NewPrinter creates a Printer that writes to w.
func NewPrinter(w io.Writer) *Printer {
	return &Printer{w: w}
}

// Print writes the summary stats to the underlying writer.
func (p *Printer) Print(s Stats) {
	fmt.Fprintf(p.w, "\nSummary:\n")
	fmt.Fprintf(p.w, "  Resources scanned : %d\n", s.TotalResources)
	fmt.Fprintf(p.w, "  Drifted           : %d\n", s.DriftedResources)
	fmt.Fprintf(p.w, "  In sync           : %d\n", s.CleanResources)
	fmt.Fprintf(p.w, "  Total changes     : %d\n", s.TotalChanges)
}
