package diff

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// EnricherPrinter renders EnrichedResult slices in a tabular format.
type EnricherPrinter struct {
	w io.Writer
}

// NewEnricherPrinter creates a printer that writes to w.
// If w is nil it defaults to os.Stdout.
func NewEnricherPrinter(w io.Writer) *EnricherPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &EnricherPrinter{w: w}
}

// Print writes the enriched results table.
func (p *EnricherPrinter) Print(results []EnrichedResult) {
	if len(results) == 0 {
		fmt.Fprintln(p.w, "No enriched drift results.")
		return
	}
	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "RESOURCE TYPE\tRESOURCE ID\tCHANGES\tURL")
	fmt.Fprintln(tw, "-------------\t-----------\t-------\t---")
	for _, r := range results {
		url := r.ResourceURL
		if url == "" {
			url = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\n",
			r.ResourceType, r.ResourceID, r.ChangeCount, url)
	}
	_ = tw.Flush()
}
