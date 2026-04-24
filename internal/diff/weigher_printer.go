package diff

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// WeigherPrinter renders a slice of WeighedResult in a tabular format.
type WeigherPrinter struct {
	w io.Writer
}

// NewWeigherPrinter creates a WeigherPrinter that writes to w.
// If w is nil, os.Stdout is used.
func NewWeigherPrinter(w io.Writer) *WeigherPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &WeigherPrinter{w: w}
}

// Print writes the weighed results table to the configured writer.
func (p *WeigherPrinter) Print(results []WeighedResult) {
	if len(results) == 0 {
		fmt.Fprintln(p.w, "No drift detected.")
		return
	}

	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "RESOURCE TYPE\tRESOURCE ID\tCHANGES\tWEIGHT")
	fmt.Fprintln(tw, "-------------\t-----------\t-------\t------")
	for _, wr := range results {
		fmt.Fprintf(tw, "%s\t%s\t%d\t%.2f\n",
			wr.Result.ResourceType,
			wr.Result.ResourceID,
			len(wr.Result.Changes),
			wr.Weight,
		)
	}
	tw.Flush()
}
