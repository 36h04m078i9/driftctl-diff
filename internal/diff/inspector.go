package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/acme/driftctl-diff/internal/drift"
)

// Inspector provides detailed per-attribute inspection of a single drifted resource.
type Inspector struct {
	w io.Writer
}

// InspectResult holds the inspection output for one resource.
type InspectResult struct {
	ResourceID   string
	ResourceType string
	TotalChanges int
	Lines        []string
}

// NewInspector returns an Inspector that writes to w (defaults to stdout).
func NewInspector(w io.Writer) *Inspector {
	if w == nil {
		w = os.Stdout
	}
	return &Inspector{w: w}
}

// Inspect finds the resource matching id inside results and returns an InspectResult.
// Returns an error if the resource is not found.
func (ins *Inspector) Inspect(results []drift.ResourceDiff, id string) (*InspectResult, error) {
	for _, r := range results {
		if r.ResourceID != id {
			continue
		}
		lines := make([]string, 0, len(r.Changes))
		for _, c := range r.Changes {
			lines = append(lines, fmt.Sprintf("  [%s] %s: %q -> %q", c.Kind, c.Attribute, c.StateValue, c.LiveValue))
		}
		return &InspectResult{
			ResourceID:   r.ResourceID,
			ResourceType: r.ResourceType,
			TotalChanges: len(r.Changes),
			Lines:        lines,
		}, nil
	}
	return nil, fmt.Errorf("resource %q not found in results", id)
}

// Print writes the InspectResult to the configured writer.
func (ins *Inspector) Print(res *InspectResult) {
	if res == nil {
		fmt.Fprintln(ins.w, "no inspection result")
		return
	}
	fmt.Fprintf(ins.w, "Resource: %s (%s)\n", res.ResourceID, res.ResourceType)
	fmt.Fprintf(ins.w, "Changes : %d\n", res.TotalChanges)
	if len(res.Lines) == 0 {
		fmt.Fprintln(ins.w, "  (no changes)")
		return
	}
	fmt.Fprintln(ins.w, strings.Join(res.Lines, "\n"))
}
