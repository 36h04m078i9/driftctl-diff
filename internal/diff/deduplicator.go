package diff

import (
	"fmt"

	"github.com/acme/driftctl-diff/internal/drift"
)

// DeduplicateOptions controls deduplication behaviour.
type DeduplicateOptions struct {
	// IgnoreCase treats attribute names as case-insensitive when deduplicating.
	IgnoreCase bool
}

// Deduplicator removes duplicate drift results across multiple sources.
type Deduplicator struct {
	opts DeduplicateOptions
}

// NewDeduplicator returns a Deduplicator with the given options.
func NewDeduplicator(opts DeduplicateOptions) *Deduplicator {
	return &Deduplicator{opts: opts}
}

// Deduplicate returns a slice of DriftResult with duplicate entries removed.
// Two results are considered duplicates when they share the same resource type,
// resource ID, attribute name (and optionally case-insensitive match), and
// change kind.
func (d *Deduplicator) Deduplicate(results []drift.DriftResult) []drift.DriftResult {
	seen := make(map[string]struct{})
	out := make([]drift.DriftResult, 0, len(results))

	for _, r := range results {
		uniqChanges := make([]drift.AttributeChange, 0, len(r.Changes))
		for _, c := range r.Changes {
			attr := c.Attribute
			if d.opts.IgnoreCase {
				attr = toLower(attr)
			}
			k := fmt.Sprintf("%s|%s|%s|%s", r.ResourceType, r.ResourceID, attr, c.Kind)
			if _, exists := seen[k]; exists {
				continue
			}
			seen[k] = struct{}{}
			uniqChanges = append(uniqChanges, c)
		}
		if len(uniqChanges) == 0 {
			continue
		}
		copy := r
		copy.Changes = uniqChanges
		out = append(out, copy)
	}
	return out
}

func toLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}
