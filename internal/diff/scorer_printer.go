package diff

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// ScorerPrinter renders a score table to a writer.
type ScorerPrinter struct {
	w io.Writer
}

// NewScorerPrinter creates a ScorerPrinter. If w is nil stdout is used.
func NewScorerPrinter(w io.Writer) *ScorerPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &ScorerPrinter{w: w}
}

// Print writes the score table.
func (p *ScorerPrinter) Print(scores []Score) {
	if len(scores) == 0 {
		fmt.Fprintln(p.w, "No drift scores to display.")
		return
	}
	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "RESOURCE TYPE\tRESOURCE ID\tCHANGED\tMISSING\tSCORE")
	for _, s := range scores {
		fmt.Fprintf(tw, "%s\t%s\t%d\t%d\t%d\n",
			s.ResourceType, s.ResourceID, s.Changed, s.Missing, s.Total)
	}
	_ = tw.Flush()
}
