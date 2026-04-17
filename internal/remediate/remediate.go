// Package remediate suggests remediation steps for detected drift.
package remediate

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/owner/driftctl-diff/internal/drift"
)

// Suggestion holds a human-readable remediation hint for a drifted resource.
type Suggestion struct {
	ResourceType string
	ResourceID   string
	Lines        []string
}

// Suggester produces remediation suggestions from drift results.
type Suggester struct {
	w io.Writer
}

// New returns a Suggester writing to w. If w is nil, os.Stdout is used.
func New(w io.Writer) *Suggester {
	if w == nil {
		w = os.Stdout
	}
	return &Suggester{w: w}
}

// Suggest generates Suggestion values for each drifted result.
func Suggest(results []drift.ResourceDiff) []Suggestion {
	var out []Suggestion
	for _, r := range results {
		if len(r.Changes) == 0 {
			continue
		}
		s := Suggestion{
			ResourceType: r.ResourceType,
			ResourceID:   r.ResourceID,
		}
		for _, c := range r.Changes {
			s.Lines = append(s.Lines, fmt.Sprintf(
				"  attribute %q: run `terraform apply` to reconcile (state=%q, live=%q)",
				c.Attribute, c.StateValue, c.LiveValue,
			))
		}
		out = append(out, s)
	}
	return out
}

// Print writes suggestions to the Suggester's writer.
func (s *Suggester) Print(suggestions []Suggestion) {
	if len(suggestions) == 0 {
		fmt.Fprintln(s.w, "No remediation needed — resources are in sync.")
		return
	}
	for _, sg := range suggestions {
		fmt.Fprintf(s.w, "[%s] %s\n", sg.ResourceType, sg.ResourceID)
		fmt.Fprintln(s.w, strings.Join(sg.Lines, "\n"))
	}
}
