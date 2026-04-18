package diff

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ClassifierPrinter renders Classifications in a human-readable format.
type ClassifierPrinter struct {
	w io.Writer
}

// NewClassifierPrinter creates a ClassifierPrinter writing to w.
// If w is nil it defaults to os.Stdout.
func NewClassifierPrinter(w io.Writer) *ClassifierPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &ClassifierPrinter{w: w}
}

// Print writes a formatted table of classifications.
func (p *ClassifierPrinter) Print(classifications []Classification) {
	if len(classifications) == 0 {
		fmt.Fprintln(p.w, "No classifications to display.")
		return
	}
	fmt.Fprintf(p.w, "%-12s %-30s %-10s %s\n", "SEVERITY", "RESOURCE", "TYPE", "REASON")
	fmt.Fprintln(p.w, strings.Repeat("-", 72))
	for _, c := range classifications {
		fmt.Fprintf(p.w, "%-12s %-30s %-10s %s\n",
			strings.ToUpper(string(c.Severity)),
			c.Result.ResourceID,
			c.Result.ResourceType,
			c.Reason,
		)
	}
}
