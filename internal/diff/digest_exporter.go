package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/acme/driftctl-diff/internal/drift"
)

// DigestReport is the JSON-serialisable representation of a digest export.
type DigestReport struct {
	GeneratedAt string `json:"generated_at"`
	Digest      string `json:"digest"`
	ResourceCount int  `json:"resource_count"`
	Drifted     bool   `json:"drifted"`
}

// DigestExporter writes a DigestReport as JSON to an io.Writer.
type DigestExporter struct {
	digest *Digest
	writer io.Writer
}

// NewDigestExporter creates a DigestExporter. A nil writer defaults to
// os.Stdout.
func NewDigestExporter(digest *Digest, writer io.Writer) *DigestExporter {
	if writer == nil {
		writer = os.Stdout
	}
	return &DigestExporter{digest: digest, writer: writer}
}

// Export serialises the digest report for results to JSON and writes it to
// the configured writer.
func (e *DigestExporter) Export(results []drift.DriftResult) error {
	drifted := false
	for _, r := range results {
		if len(r.Changes) > 0 {
			drifted = true
			break
		}
	}

	report := DigestReport{
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		Digest:        e.digest.Compute(results),
		ResourceCount: len(results),
		Drifted:       drifted,
	}

	enc := json.NewEncoder(e.writer)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return fmt.Errorf("digest exporter: %w", err)
	}
	return nil
}
