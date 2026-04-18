package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/acme/driftctl-diff/internal/drift"
)

// RenderOptions controls how diffs are rendered.
type RenderOptions struct {
	Color   bool
	Context int // number of unchanged lines of context (reserved for future use)
}

// Renderer renders a slice of DriftResult in a unified-diff-like format.
type Renderer struct {
	opts        RenderOptions
	highlighter *Highlighter
	w           io.Writer
}

// NewRenderer creates a Renderer writing to w (defaults to os.Stdout).
func NewRenderer(opts RenderOptions, w io.Writer) *Renderer {
	if w == nil {
		w = os.Stdout
	}
	return &Renderer{
		opts:        opts,
		highlighter: NewHighlighter(opts.Color),
		w:           w,
	}
}

// Render writes a human-readable diff of results to the writer.
func (r *Renderer) Render(results []drift.DriftResult) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(r.w, "No drift detected.")
		return err
	}
	for _, res := range results {
		header := fmt.Sprintf("--- %s/%s (state)\n+++ %s/%s (live)",
			res.ResourceType, res.ResourceID,
			res.ResourceType, res.ResourceID)
		if _, err := fmt.Fprintln(r.w, header); err != nil {
			return err
		}
		for _, ch := range res.Changes {
			lines := r.highlighter.Highlight(ch)
			if _, err := fmt.Fprintln(r.w, strings.Join(lines, "\n")); err != nil {
				return err
			}
		}
	}
	return nil
}
