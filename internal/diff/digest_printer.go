package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/acme/driftctl-diff/internal/drift"
)

// DigestPrinter writes the computed digest to a writer in a human-readable
// one-liner format suitable for scripting and CI pipelines.
type DigestPrinter struct {
	digest  *Digest
	writer  io.Writer
}

// NewDigestPrinter creates a DigestPrinter. If writer is nil, os.Stdout is
// used.
func NewDigestPrinter(digest *Digest, writer io.Writer) *DigestPrinter {
	if writer == nil {
		writer = os.Stdout
	}
	return &DigestPrinter{digest: digest, writer: writer}
}

// Print computes the digest for results and writes it to the configured
// writer.
func (p *DigestPrinter) Print(results []drift.DriftResult) error {
	hash := p.digest.Compute(results)
	_, err := fmt.Fprintf(p.writer, "drift-digest: %s\n", hash)
	return err
}
