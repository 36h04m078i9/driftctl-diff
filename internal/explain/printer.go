package explain

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Printer writes explanations to an io.Writer.
type Printer struct {
	w io.Writer
}

// NewPrinter returns a Printer writing to w. If w is nil, os.Stdout is used.
func NewPrinter(w io.Writer) *Printer {
	if w == nil {
		w = os.Stdout
	}
	return &Printer{w: w}
}

// Print writes all explanations in a human-readable format.
func (p *Printer) Print(exps []Explanation) {
	if len(exps) == 0 {
		fmt.Fprintln(p.w, "No drift explanations.")
		return
	}
	for _, ex := range exps {
		fmt.Fprintf(p.w, "[%s] %s (%s)\n  %s\n",
			strings.ToUpper(string(ex.Severity)),
			ex.ResourceID,
			ex.Attribute,
			ex.Message,
		)
	}
}
