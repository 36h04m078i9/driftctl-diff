package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/acme/driftctl-diff/internal/drift"
)

// PipelinePrinter renders a summary of a Pipeline execution to an io.Writer.
type PipelinePrinter struct {
	w io.Writer
}

// NewPipelinePrinter returns a PipelinePrinter that writes to w.
// If w is nil it defaults to os.Stdout.
func NewPipelinePrinter(w io.Writer) *PipelinePrinter {
	if w == nil {
		w = os.Stdout
	}
	return &PipelinePrinter{w: w}
}

// Print writes a human-readable summary of results produced by the pipeline.
func (pp *PipelinePrinter) Print(results []drift.DriftResult, steps int) {
	fmt.Fprintf(pp.w, "Pipeline steps executed : %d\n", steps)
	fmt.Fprintf(pp.w, "Resources after pipeline: %d\n", len(results))

	drifted := 0
	for _, r := range results {
		if len(r.Changes) > 0 {
			drifted++
		}
	}
	fmt.Fprintf(pp.w, "Drifted resources       : %d\n", drifted)
}
