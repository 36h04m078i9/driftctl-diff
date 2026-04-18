package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/acme/driftctl-diff/internal/drift"
)

// DiffOptions controls how the unified diff is rendered.
type DiffOptions struct {
	ContextLines int  // number of unchanged lines shown around each change
	Color        bool // whether to emit ANSI colour codes
}

// DefaultDiffOptions returns sensible defaults.
func DefaultDiffOptions() DiffOptions {
	return DiffOptions{ContextLines: 3, Color: false}
}

// Differ renders a unified-style diff for a single DriftResult.
type Differ struct {
	opts DiffOptions
	out  io.Writer
}

// NewDiffer creates a Differ. If w is nil it defaults to os.Stdout.
func NewDiffer(opts DiffOptions, w io.Writer) *Differ {
	if w == nil {
		w = os.Stdout
	}
	return &Differ{opts: opts, out: w}
}

// Render writes a unified diff for the given result to the writer.
func (d *Differ) Render(result drift.DriftResult) error {
	if len(result.Changes) == 0 {
		_, err := fmt.Fprintln(d.out, "(no drift)")
		return err
	}

	header := fmt.Sprintf("--- %s/%s (state)\n+++ %s/%s (live)",
		result.ResourceType, result.ResourceID,
		result.ResourceType, result.ResourceID)
	_, err := fmt.Fprintln(d.out, header)
	if err != nil {
		return err
	}

	for _, ch := range result.Changes {
		lines := d.formatChange(ch)
		for _, l := range lines {
			if _, err := fmt.Fprintln(d.out, l); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *Differ) formatChange(ch drift.AttributeChange) []string {
	red := func(s string) string {
		if d.opts.Color {
			return "\033[31m" + s + "\033[0m"
		}
		return s
	}
	green := func(s string) string {
		if d.opts.Color {
			return "\033[32m" + s + "\033[0m"
		}
		return s
	}

	var out []string
	out = append(out, fmt.Sprintf("@@ attribute: %s @@", ch.Attribute))
	if strings.TrimSpace(fmt.Sprintf("%v", ch.StateValue)) != "" {
		out = append(out, red(fmt.Sprintf("-  %v", ch.StateValue)))
	}
	if strings.TrimSpace(fmt.Sprintf("%v", ch.LiveValue)) != "" {
		out = append(out, green(fmt.Sprintf("+  %v", ch.LiveValue)))
	}
	return out
}
