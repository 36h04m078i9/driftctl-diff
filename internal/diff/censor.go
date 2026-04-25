package diff

import (
	"strings"

	"github.com/acme/driftctl-diff/internal/drift"
)

// DefaultCensorOptions returns CensorOptions with a standard set of sensitive
// attribute name patterns that should have their values masked in output.
func DefaultCensorOptions() CensorOptions {
	return CensorOptions{
		Patterns:    []string{"password", "secret", "token", "key", "private"},
		Replacement: "[REDACTED]",
	}
}

// CensorOptions controls how the Censor masks sensitive attribute values.
type CensorOptions struct {
	// Patterns is a list of lowercase substrings; any attribute whose name
	// contains one of these substrings will have its values replaced.
	Patterns []string
	// Replacement is the string written in place of the real value.
	Replacement string
}

// Censor masks sensitive attribute values in drift results so that secrets are
// never surfaced in plain-text output.
type Censor struct {
	opts CensorOptions
}

// NewCensor creates a Censor with the given options.  Call
// DefaultCensorOptions() to get a sensible starting point.
func NewCensor(opts CensorOptions) *Censor {
	return &Censor{opts: opts}
}

// Apply returns a copy of results where every change whose attribute name
// matches a sensitive pattern has its Got and Want values replaced with the
// configured replacement string.
func (c *Censor) Apply(results []drift.ResourceDiff) []drift.ResourceDiff {
	out := make([]drift.ResourceDiff, len(results))
	for i, r := range results {
		copy := r
		copy.Changes = make([]drift.Change, len(r.Changes))
		for j, ch := range r.Changes {
			if c.isSensitive(ch.Attribute) {
				ch.Got = c.opts.Replacement
				ch.Want = c.opts.Replacement
			}
			copy.Changes[j] = ch
		}
		out[i] = copy
	}
	return out
}

func (c *Censor) isSensitive(attr string) bool {
	lower := strings.ToLower(attr)
	for _, p := range c.opts.Patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}
