package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/acme/driftctl-diff/internal/drift"
)

// TableFormatter renders drift results as an aligned table.
type TableFormatter struct {
	colorizer *Colorizer
}

// NewTableFormatter creates a TableFormatter with the given colorizer.
func NewTableFormatter(c *Colorizer) *TableFormatter {
	return &TableFormatter{colorizer: c}
}

// Write outputs a tabular summary of drift changes to w.
func (t *TableFormatter) Write(w io.Writer, results []drift.ResourceDiff) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "No drift detected.")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	header := fmt.Sprintf("%s\t%s\t%s\t%s",
		"RESOURCE TYPE", "RESOURCE ID", "ATTRIBUTE", "CHANGE")
	fmt.Fprintln(tw, header)
	fmt.Fprintln(tw, strings.Repeat("-", 72))

	for _, rd := range results {
		for _, ch := range rd.Changes {
			kind := formatKind(ch)
			change := fmt.Sprintf("%s -> %s", ch.StateValue, ch.LiveValue)
			if ch.StateValue == "" {
				change = fmt.Sprintf("(missing) -> %s", ch.LiveValue)
			} else if ch.LiveValue == "" {
				change = fmt.Sprintf("%s -> (missing)", ch.StateValue)
			}
			line := fmt.Sprintf("%s\t%s\t%s\t%s",
				rd.ResourceType, rd.ResourceID, ch.Attribute, change)
			line = t.colorizer.Colorize(kind, line)
			fmt.Fprintln(tw, line)
		}
	}
	return nil
}

func formatKind(ch drift.Change) string {
	if ch.StateValue == "" || ch.LiveValue == "" {
		return "missing"
	}
	return "changed"
}
