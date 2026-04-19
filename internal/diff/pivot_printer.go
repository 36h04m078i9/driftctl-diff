package diff

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// PivotPrinter renders a []PivotEntry as a tab-aligned table.
type PivotPrinter struct {
	w io.Writer
}

// NewPivotPrinter returns a PivotPrinter writing to w. If w is nil stdout is
// used.
func NewPivotPrinter(w io.Writer) *PivotPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &PivotPrinter{w: w}
}

// Print writes the pivot table to the configured writer.
func (pp *PivotPrinter) Print(entries []PivotEntry) {
	if len(entries) == 0 {
		fmt.Fprintln(pp.w, "no attribute-level drift detected")
		return
	}

	tw := tabwriter.NewWriter(pp.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ATTRIBUTE\tRESOURCE TYPE\tRESOURCE ID\tSTATE VALUE\tLIVE VALUE")
	fmt.Fprintln(tw, "---------\t-------------\t-----------\t-----------\t----------")
	for _, e := range entries {
		for _, r := range e.Resources {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
				e.Attribute, r.ResourceType, r.ResourceID, r.StateValue, r.LiveValue)
		}
	}
	_ = tw.Flush()
}
