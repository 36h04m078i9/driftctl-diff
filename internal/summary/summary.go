// Package summary computes and prints high-level drift statistics.
package summary

import (
	"fmt"
	"io"
	"os"

	"github.com/acme/driftctl-diff/internal/drift"
)

// Result holds aggregate counts for a drift detection run.
type Result struct {
	Total   int `json:"total"`
	Drifted int `json:"drifted"`
	Clean   int `json:"clean"`
}

// Compute derives a Result from a slice of DriftResults.
func Compute(results []drift.DriftResult) Result {
	r := Result{Total: len(results)}
	for _, res := range results {
		if len(res.Changes) > 0 {
			r.Drifted++
		} else {
			r.Clean++
		}
	}
	return r
}

// Printer writes a human-readable summary line.
type Printer struct {
	dest io.Writer
}

// NewPrinter creates a Printer writing to dest (defaults to os.Stdout).
func NewPrinter(dest io.Writer) *Printer {
	if dest == nil {
		dest = os.Stdout
	}
	return &Printer{dest: dest}
}

// Print writes a one-line summary of the result.
func (p *Printer) Print(r Result) {
	status := "✔  No drift detected"
	if r.Drifted > 0 {
		status = fmt.Sprintf("✘  Drift detected in %d resource(s)", r.Drifted)
	}
	fmt.Fprintf(p.dest, "%s  (total: %d, clean: %d, drifted: %d)\n",
		status, r.Total, r.Clean, r.Drifted)
}
