package diff

import (
	"strings"

	"github.com/acme/driftctl-diff/internal/drift"
)

// NormalizeOptions controls how attribute values are normalized before comparison.
type NormalizeOptions struct {
	TrimSpace   bool
	LowerCase   bool
	StripQuotes bool
}

// DefaultNormalizeOptions returns sensible defaults.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		TrimSpace:   true,
		LowerCase:   false,
		StripQuotes: false,
	}
}

// Normalizer applies value normalization to drift results so that cosmetic
// differences (whitespace, quoting) do not surface as real drift.
type Normalizer struct {
	opts NormalizeOptions
}

// NewNormalizer creates a Normalizer with the given options.
func NewNormalizer(opts NormalizeOptions) *Normalizer {
	return &Normalizer{opts: opts}
}

// Normalize returns a new slice of DriftResult with attribute values normalized.
// Changes whose Got and Want become equal after normalization are dropped.
func (n *Normalizer) Normalize(results []drift.DriftResult) []drift.DriftResult {
	out := make([]drift.DriftResult, 0, len(results))
	for _, r := range results {
		filtered := make([]drift.AttributeChange, 0, len(r.Changes))
		for _, c := range r.Changes {
			got := n.normalize(c.Got)
			want := n.normalize(c.Want)
			if got == want {
				continue
			}
			filtered = append(filtered, drift.AttributeChange{
				Attribute: c.Attribute,
				Got:       got,
				Want:      want,
				Kind:      c.Kind,
			})
		}
		if len(filtered) == 0 {
			continue
		}
		out = append(out, drift.DriftResult{
			ResourceID:   r.ResourceID,
			ResourceType: r.ResourceType,
			Changes:      filtered,
		})
	}
	return out
}

func (n *Normalizer) normalize(v string) string {
	if n.opts.TrimSpace {
		v = strings.TrimSpace(v)
	}
	if n.opts.LowerCase {
		v = strings.ToLower(v)
	}
	if n.opts.StripQuotes {
		v = strings.Trim(v, `"`)
	}
	return v
}
