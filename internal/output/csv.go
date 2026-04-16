package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/snyk/driftctl-diff/internal/drift"
)

// CSVFormatter writes drift results as a CSV file.
type CSVFormatter struct {
	w io.Writer
}

// NewCSVFormatter returns a CSVFormatter writing to w.
// If w is nil it defaults to os.Stdout.
func NewCSVFormatter(w io.Writer) *CSVFormatter {
	if w == nil {
		w = os.Stdout
	}
	return &CSVFormatter{w: w}
}

// Format writes changes in CSV format.
func (f *CSVFormatter) Format(changes []drift.ResourceDiff) error {
	cw := csv.NewWriter(f.w)
	defer cw.Flush()

	if err := cw.Write([]string{"resource_type", "resource_id", "attribute", "kind", "state_value", "live_value"}); err != nil {
		return fmt.Errorf("csv header: %w", err)
	}

	for _, rd := range changes {
		for _, ch := range rd.Changes {
			row := []string{
				rd.ResourceType,
				rd.ResourceID,
				ch.Attribute,
				ch.Kind.String(),
				ch.StateValue,
				ch.LiveValue,
			}
			if err := cw.Write(row); err != nil {
				return fmt.Errorf("csv row: %w", err)
			}
		}
	}

	if err := cw.Error(); err != nil {
		return fmt.Errorf("csv flush: %w", err)
	}
	return nil
}
