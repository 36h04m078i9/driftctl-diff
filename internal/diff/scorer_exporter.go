package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// ScorerExporter writes scores to an io.Writer in the requested format.
type ScorerExporter struct {
	w io.Writer
}

// NewScorerExporter creates a ScorerExporter. If w is nil stdout is used.
func NewScorerExporter(w io.Writer) *ScorerExporter {
	if w == nil {
		w = os.Stdout
	}
	return &ScorerExporter{w: w}
}

// Export writes scores in the specified format ("json" or "text").
func (e *ScorerExporter) Export(scores []Score, format string) error {
	switch format {
	case "json":
		return e.exportJSON(scores)
	case "text", "":
		NewScorerPrinter(e.w).Print(scores)
		return nil
	default:
		return fmt.Errorf("unsupported scorer export format: %q", format)
	}
}

func (e *ScorerExporter) exportJSON(scores []Score) error {
	enc := json.NewEncoder(e.w)
	enc.SetIndent("", "  ")
	return enc.Encode(scores)
}
