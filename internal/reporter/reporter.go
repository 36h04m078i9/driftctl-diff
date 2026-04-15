// Package reporter provides structured reporting of drift detection results
// to various output targets such as JSON files or stdout.
package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/acme/driftctl-diff/internal/drift"
	"github.com/acme/driftctl-diff/internal/summary"
)

// Report holds the full output of a drift detection run.
type Report struct {
	GeneratedAt time.Time          `json:"generated_at"`
	Summary     summary.Result     `json:"summary"`
	Changes     []drift.DriftResult `json:"changes"`
}

// Reporter writes a Report to a destination.
type Reporter struct {
	dest io.Writer
}

// New creates a Reporter that writes to dest.
// If dest is nil, os.Stdout is used.
func New(dest io.Writer) *Reporter {
	if dest == nil {
		dest = os.Stdout
	}
	return &Reporter{dest: dest}
}

// WriteJSON serialises the report as indented JSON.
func (r *Reporter) WriteJSON(report Report) error {
	enc := json.NewEncoder(r.dest)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return fmt.Errorf("reporter: encode JSON: %w", err)
	}
	return nil
}

// Build constructs a Report from drift results and a pre-computed summary.
func Build(changes []drift.DriftResult, sum summary.Result) Report {
	return Report{
		GeneratedAt: time.Now().UTC(),
		Summary:     sum,
		Changes:     changes,
	}
}
