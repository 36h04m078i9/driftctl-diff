package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/owner/driftctl-diff/internal/drift"
)

// PatchExporter writes patch output to a file or writer.
type PatchExporter struct {
	patcher *Patcher
	out     io.Writer
}

// NewPatchExporter creates a PatchExporter with the given options.
func NewPatchExporter(opts PatchOptions, out io.Writer) *PatchExporter {
	if out == nil {
		out = os.Stdout
	}
	return &PatchExporter{
		patcher: NewPatcher(opts, out),
		out:     out,
	}
}

// Export writes patch output for results, prefixed with a run header.
func (pe *PatchExporter) Export(results []drift.ResourceDiff) error {
	if _, err := fmt.Fprintf(pe.out, "# driftctl-diff patch\n# resources: %d\n", len(results)); err != nil {
		return err
	}
	return pe.patcher.Generate(results)
}
