package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/snyk/driftctl-diff/internal/drift"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBold   = "\033[1m"
)

// Formatter renders drift results to a writer.
type Formatter struct {
	w     io.Writer
	color bool
}

// NewFormatter creates a new Formatter writing to w.
func NewFormatter(w io.Writer, color bool) *Formatter {
	return &Formatter{w: w, color: color}
}

// Render writes a human-readable diff of the provided drift changes.
func (f *Formatter) Render(changes []drift.Change) {
	if len(changes) == 0 {
		fmt.Fprintf(f.w, "%sNo drift detected.%s\n", f.bold(), colorReset)
		return
	}

	// Group changes by resource ID
	grouped := make(map[string][]drift.Change)
	for _, c := range changes {
		grouped[c.ResourceID] = append(grouped[c.ResourceID], c)
	}

	for resourceID, drifts := range grouped {
		fmt.Fprintf(f.w, "\n%s~ %s%s\n", f.yellow(), resourceID, colorReset)
		fmt.Fprintf(f.w, "%s\n", strings.Repeat("-", 40))
		for _, d := range drifts {
			fmt.Fprintf(f.w, "  %s- %s: %v%s\n", f.red(), d.Attribute, d.Expected, colorReset)
			fmt.Fprintf(f.w, "  %s+ %s: %v%s\n", f.green(), d.Attribute, d.Actual, colorReset)
		}
	}

	fmt.Fprintf(f.w, "\n%sSummary: %d resource(s) drifted, %d change(s) total.%s\n",
		f.bold(), len(grouped), len(changes), colorReset)
}

func (f *Formatter) red() string {
	if f.color {
		return colorRed
	}
	return ""
}

func (f *Formatter) green() string {
	if f.color {
		return colorGreen
	}
	return ""
}

func (f *Formatter) yellow() string {
	if f.color {
		return colorYellow
	}
	return ""
}

func (f *Formatter) bold() string {
	if f.color {
		return colorBold
	}
	return ""
}
