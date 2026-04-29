package diff

import (
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/acme/driftctl-diff/internal/drift"
)

// RouterPrinter prints a summary of how results were routed.
type RouterPrinter struct {
	w io.Writer
}

// NewRouterPrinter returns a RouterPrinter that writes to w.
// If w is nil, os.Stdout is used.
func NewRouterPrinter(w io.Writer) *RouterPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &RouterPrinter{w: w}
}

// Print writes a table showing each route label and the count of resources
// dispatched to it, based on the provided results and label key.
func (p *RouterPrinter) Print(results []drift.DriftResult, labelKey string) {
	if labelKey == "" {
		labelKey = "env"
	}

	counts := make(map[string]int)
	for _, r := range results {
		label := "(default)"
		if r.Metadata != nil {
			if v, ok := r.Metadata[labelKey]; ok && v != "" {
				label = v
			}
		}
		counts[label]++
	}

	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ROUTE\tRESOURCES")
	for _, k := range keys {
		fmt.Fprintf(tw, "%s\t%d\n", k, counts[k])
	}
	_ = tw.Flush()
}
