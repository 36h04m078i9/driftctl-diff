package diff

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/acme/driftctl-diff/internal/drift"
)

// TransformerPrinter renders a before/after summary of a Transformer run to
// an io.Writer, showing how many resources and attributes were affected.
type TransformerPrinter struct {
	w io.Writer
}

// NewTransformerPrinter returns a TransformerPrinter that writes to w.
// If w is nil it defaults to os.Stdout.
func NewTransformerPrinter(w io.Writer) *TransformerPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &TransformerPrinter{w: w}
}

// Print writes a human-readable summary comparing before and after slices.
func (p *TransformerPrinter) Print(before, after []drift.DriftResult) {
	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "FIELD\tBEFORE\tAFTER")
	fmt.Fprintf(tw, "Resources\t%d\t%d\n", len(before), len(after))
	fmt.Fprintf(tw, "Attribute changes\t%d\t%d\n", totalChanges(before), totalChanges(after))
	tw.Flush()
}

func totalChanges(results []drift.DriftResult) int {
	n := 0
	for _, r := range results {
		n += len(r.Changes)
	}
	return n
}
