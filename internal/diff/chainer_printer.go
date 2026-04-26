package diff

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// ChainerPrinter renders a ChainResult pipeline summary to a writer.
type ChainerPrinter struct {
	w io.Writer
}

// NewChainerPrinter returns a ChainerPrinter. If w is nil it defaults to
// os.Stdout.
func NewChainerPrinter(w io.Writer) *ChainerPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &ChainerPrinter{w: w}
}

// Print writes a human-readable table of each pipeline step and the number
// of drift results that remained after it ran.
func (p *ChainerPrinter) Print(cr ChainResult) {
	if len(cr.Steps) == 0 {
		fmt.Fprintln(p.w, "no pipeline steps were executed")
		return
	}

	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "STEP\tRESULTS REMAINING")
	for _, s := range cr.Steps {
		fmt.Fprintf(tw, "%s\t%d\n", s.Name, s.Count)
	}
	_ = tw.Flush()

	fmt.Fprintf(p.w, "\nfinal result count: %d\n", len(cr.Results))
}
