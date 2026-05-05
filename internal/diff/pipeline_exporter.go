package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/acme/driftctl-diff/internal/drift"
)

// PipelineExporter writes the output of a Pipeline run to an io.Writer in
// either "text" or "json" format.
type PipelineExporter struct {
	w      io.Writer
	format string
}

// NewPipelineExporter returns a PipelineExporter.
// format must be "text" or "json". w defaults to os.Stdout when nil.
func NewPipelineExporter(w io.Writer, format string) *PipelineExporter {
	if w == nil {
		w = os.Stdout
	}
	return &PipelineExporter{w: w, format: format}
}

// Export writes results to the configured writer.
func (pe *PipelineExporter) Export(results []drift.DriftResult) error {
	switch pe.format {
	case "json":
		enc := json.NewEncoder(pe.w)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	case "text", "":
		for _, r := range results {
			fmt.Fprintf(pe.w, "[%s] %s  changes=%d\n",
				r.ResourceType, r.ResourceID, len(r.Changes))
		}
		return nil
	default:
		return fmt.Errorf("unsupported format %q", pe.format)
	}
}
