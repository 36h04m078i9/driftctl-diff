package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/acme/driftctl-diff/internal/drift"
)

// MarkdownFormatter renders drift results as a Markdown document.
type MarkdownFormatter struct{}

// NewMarkdownFormatter returns a new MarkdownFormatter.
func NewMarkdownFormatter() *MarkdownFormatter {
	return &MarkdownFormatter{}
}

// Format writes drift results to w in Markdown format.
func (f *MarkdownFormatter) Format(w io.Writer, results []drift.ResourceDiff) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "## Drift Report\n\n✅ No drift detected.")
		return err
	}

	var sb strings.Builder
	sb.WriteString("## Drift Report\n\n")
	sb.WriteString(fmt.Sprintf("**%d resource(s) drifted.**\n\n", len(results)))

	for _, r := range results {
		sb.WriteString(fmt.Sprintf("### `%s` — `%s`\n\n", r.ResourceType, r.ResourceID))
		sb.WriteString("| Attribute | State Value | Live Value |\n")
		sb.WriteString("|-----------|-------------|------------|\n")
		for _, ch := range r.Changes {
			sb.WriteString(fmt.Sprintf("| `%s` | `%v` | `%v` |\n",
				ch.Attribute, ch.StateValue, ch.LiveValue))
		}
		sb.WriteString("\n")
	}

	_, err := fmt.Fprint(w, sb.String())
	return err
}
