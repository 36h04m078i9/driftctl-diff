package metrics

import (
	"fmt"
	"io"
	"os"
)

// Printer writes a Counters summary to a writer.
type Printer struct {
	w io.Writer
}

// NewPrinter creates a Printer. If w is nil it defaults to os.Stdout.
func NewPrinter(w io.Writer) *Printer {
	if w == nil {
		w = os.Stdout
	}
	return &Printer{w: w}
}

// Print writes human-readable metrics to the underlying writer.
func (p *Printer) Print(c Counters) {
	fmt.Fprintf(p.w, "Scan duration   : %s\n", c.Duration.Round(1000000))
	fmt.Fprintf(p.w, "Resources total : %d\n", c.ResourcesTotal)
	fmt.Fprintf(p.w, "Resources drifted: %d\n", c.ResourcesDrifted)
	fmt.Fprintf(p.w, "Attributes checked: %d\n", c.AttributesChecked)
	fmt.Fprintf(p.w, "Fetch errors    : %d\n", c.FetchErrors)
}
