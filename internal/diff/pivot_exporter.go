package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/acme/driftctl-diff/internal/drift"
)

// PivotExporter exports pivot data in text or JSON format.
type PivotExporter struct {
	pivot   *Pivot
	printer *PivotPrinter
}

// NewPivotExporter returns a PivotExporter writing to w.
func NewPivotExporter(w io.Writer) *PivotExporter {
	if w == nil {
		w = os.Stdout
	}
	return &PivotExporter{
		pivot:   NewPivot(),
		printer: NewPivotPrinter(w),
	}
}

// Export writes drift results in the requested format ("text" or "json").
func (pe *PivotExporter) Export(results []drift.ResourceDiff, format string, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}
	entries := pe.pivot.Compute(results)
	switch format {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(entries)
	case "text", "":
		NewPivotPrinter(w).Print(entries)
		return nil
	default:
		return fmt.Errorf("pivot exporter: unsupported format %q", format)
	}
}
