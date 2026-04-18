package diff

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// ScorerPrinter renders a slice of Score to a writer.
type ScorerPrinter struct {
	w io.Writer
}

// NewScorerPrinter returns a ScorerPrinter. If w is nil, os.Stdout is used.
func NewScorerPrinter(w io.Writer) *ScorerPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &ScorerPrinter{w: w}
}

// Print writes a formatted score table.
func (p *ScorerPrinter) Print(scores []Score) {
	if len(scores) == 0 {
		fmt.Fprintln(p.w, "no drift scores to display")
		return
	}
	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "RESOURCE TYPE\tRESOURCE ID\tSCORE\tCRITICAL\tWARNING\tINFO")
	fmt.Fprintln(tw, "-------------\t-----------\t-----\t--------\t-------\t----")
	for _, s := range scores {
		fmt.Fprintf(tw, "%s\t%s\t%d\t%d\t%d\t%d\n",
			s.ResourceType, s.ResourceID, s.Total, s.Critical, s.Warning, s.Info)
	}
	_ = tw.Flush()
}
