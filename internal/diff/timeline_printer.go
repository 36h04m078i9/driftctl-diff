package diff

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// TimelinePrinter renders a Timeline to a human-readable table.
type TimelinePrinter struct {
	w io.Writer
}

// NewTimelinePrinter creates a TimelinePrinter writing to w.
// If w is nil, os.Stdout is used.
func NewTimelinePrinter(w io.Writer) *TimelinePrinter {
	if w == nil {
		w = os.Stdout
	}
	return &TimelinePrinter{w: w}
}

// Print writes the timeline summary table.
func (p *TimelinePrinter) Print(tl *Timeline) {
	entries := tl.Entries()
	if len(entries) == 0 {
		fmt.Fprintln(p.w, "No timeline entries recorded.")
		return
	}

	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "#\tCaptured At\tDrifted Resources\tTotal Changes")
	fmt.Fprintln(tw, "--\t------------\t-----------------\t-------------")
	for i, e := range entries {
		drifted := 0
		totalChanges := 0
		for _, r := range e.Results {
			if len(r.Changes) > 0 {
				drifted++
				totalChanges += len(r.Changes)
			}
		}
		fmt.Fprintf(tw, "%d\t%s\t%d\t%d\n",
			i+1,
			e.CapturedAt.Format("2006-01-02 15:04:05"),
			drifted,
			totalChanges,
		)
	}
	tw.Flush()
}
