package diff

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/acme/driftctl-diff/internal/drift"
)

// RewriterPrinter renders a summary of rewritten drift results to a writer.
type RewriterPrinter struct {
	w io.Writer
}

// NewRewriterPrinter creates a RewriterPrinter. If w is nil, os.Stdout is used.
func NewRewriterPrinter(w io.Writer) *RewriterPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &RewriterPrinter{w: w}
}

// Print writes a human-readable table of rewritten results.
func (p *RewriterPrinter) Print(results []drift.DriftResult) {
	if len(results) == 0 {
		fmt.Fprintln(p.w, "No drift results to display.")
		return
	}
	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "RESOURCE TYPE\tRESOURCE ID\tATTRIBUTE\tWANT\tGOT")
	fmt.Fprintln(tw, "-------------\t-----------\t---------\t----\t---")
	for _, res := range results {
		for _, ch := range res.Changes {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
				res.ResourceType,
				res.ResourceID,
				ch.Attribute,
				ch.WantValue,
				ch.GotValue,
			)
		}
	}
	tw.Flush()
}
