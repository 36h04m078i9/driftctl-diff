package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/owner/driftctl-diff/internal/drift"
)

// PatchOptions controls patch generation behaviour.
type PatchOptions struct {
	ContextLines int
	Header       bool
}

// DefaultPatchOptions returns sensible defaults.
func DefaultPatchOptions() PatchOptions {
	return PatchOptions{ContextLines: 3, Header: true}
}

// Patcher generates unified-diff style patch output from drift results.
type Patcher struct {
	opts PatchOptions
	out  io.Writer
}

// NewPatcher creates a Patcher writing to out (defaults to stdout).
func NewPatcher(opts PatchOptions, out io.Writer) *Patcher {
	if out == nil {
		out = os.Stdout
	}
	return &Patcher{opts: opts, out: out}
}

// Generate writes a patch-style representation of results to the writer.
func (p *Patcher) Generate(results []drift.ResourceDiff) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(p.out, "No drift detected.")
		return err
	}
	for _, r := range results {
		if p.opts.Header {
			fmt.Fprintf(p.out, "--- terraform: %s/%s\n", r.ResourceType, r.ResourceID)
			fmt.Fprintf(p.out, "+++ live:      %s/%s\n", r.ResourceType, r.ResourceID)
		}
		for _, c := range r.Changes {
			p.writeHunk(c)
		}
	}
	return nil
}

func (p *Patcher) writeHunk(c drift.AttributeChange) {
	ctx := strings.Repeat(" ", p.opts.ContextLines)
	_ = ctx
	fmt.Fprintf(p.out, "@@ .%s @@\n", c.Attribute)
	fmt.Fprintf(p.out, "-%s\n", c.TerraformValue)
	fmt.Fprintf(p.out, "+%s\n", c.LiveValue)
}
