package diff

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"github.com/acme/driftctl-diff/internal/drift"
)

// DigestOptions configures how a Digest is computed.
type DigestOptions struct {
	// IncludeKind includes the change kind in the hash input when true.
	IncludeKind bool
}

// DefaultDigestOptions returns sensible defaults.
func DefaultDigestOptions() DigestOptions {
	return DigestOptions{IncludeKind: true}
}

// Digest produces a stable SHA-256 fingerprint for a slice of DriftResult.
// Two slices with the same resources and attributes (regardless of order)
// will produce the same digest, making it suitable for change-detection
// across runs.
type Digest struct {
	opts DigestOptions
}

// NewDigest creates a Digest with the supplied options.
func NewDigest(opts DigestOptions) *Digest {
	return &Digest{opts: opts}
}

// Compute returns a hex-encoded SHA-256 hash of the provided results.
func (d *Digest) Compute(results []drift.DriftResult) string {
	if len(results) == 0 {
		return emptySHA256()
	}

	lines := make([]string, 0, len(results))
	for _, r := range results {
		for _, c := range r.Changes {
			parts := []string{r.ResourceType, r.ResourceID, c.Attribute, c.Expected, c.Actual}
			if d.opts.IncludeKind {
				parts = append(parts, string(c.Kind))
			}
			lines = append(lines, strings.Join(parts, "|"))
		}
	}

	sort.Strings(lines)
	h := sha256.Sum256([]byte(strings.Join(lines, "\n")))
	return fmt.Sprintf("%x", h)
}

func emptySHA256() string {
	h := sha256.Sum256([]byte{})
	return fmt.Sprintf("%x", h)
}
