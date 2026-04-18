package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/owner/driftctl-diff/internal/drift"
)

// Exporter writes drift results to an external destination in a chosen format.
type Exporter struct {
	format string
	dest   io.Writer
}

// NewExporter creates an Exporter that writes to dest using the given format.
// Supported formats: "json", "text". If dest is nil, os.Stdout is used.
func NewExporter(format string, dest io.Writer) *Exporter {
	if dest == nil {
		dest = os.Stdout
	}
	return &Exporter{format: format, dest: dest}
}

// Export serialises results according to the configured format.
func (e *Exporter) Export(results []drift.ResourceDiff) error {
	switch e.format {
	case "json":
		return e.exportJSON(results)
	case "text", "":
		return e.exportText(results)
	default:
		return fmt.Errorf("unsupported export format: %q", e.format)
	}
}

func (e *Exporter) exportJSON(results []drift.ResourceDiff) error {
	enc := json.NewEncoder(e.dest)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

func (e *Exporter) exportText(results []drift.ResourceDiff) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(e.dest, "No drift detected.")
		return err
	}
	for _, r := range results {
		if _, err := fmt.Fprintf(e.dest, "[%s] %s\n", r.ResourceType, r.ResourceID); err != nil {
			return err
		}
		for _, c := range r.Changes {
			if _, err := fmt.Fprintf(e.dest, "  ~ %s: %q -> %q\n", c.Attribute, c.StateValue, c.LiveValue); err != nil {
				return err
			}
		}
	}
	return nil
}
