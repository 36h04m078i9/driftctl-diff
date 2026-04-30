package diff

import (
	"fmt"
	"io"
	"os"
)

// LinkerPrinter renders a LinkerResult to a writer.
type LinkerPrinter struct {
	w io.Writer
}

// NewLinkerPrinter creates a LinkerPrinter. If w is nil stdout is used.
func NewLinkerPrinter(w io.Writer) *LinkerPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &LinkerPrinter{w: w}
}

// Print writes a human-readable summary of the LinkerResult.
func (p *LinkerPrinter) Print(lr LinkerResult) {
	fmt.Fprintf(p.w, "Resources: %d\n", len(lr.Results))
	if len(lr.Links) == 0 {
		fmt.Fprintln(p.w, "No links found.")
		return
	}
	fmt.Fprintf(p.w, "Links (%d):\n", len(lr.Links))
	for _, lnk := range lr.Links {
		fmt.Fprintf(p.w, "  %s <-> %s  [shared: %s]\n",
			lnk.SourceID, lnk.TargetID, lnk.SharedAttribute)
	}
}
