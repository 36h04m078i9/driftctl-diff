package diff

import (
	"strings"

	"github.com/acme/driftctl-diff/internal/drift"
)

// TransformOptions controls how results are transformed.
type TransformOptions struct {
	// PrefixResourceIDs prepends a string to every resource ID.
	PrefixResourceIDs string
	// SuffixResourceTypes appends a string to every resource type.
	SuffixResourceTypes string
	// UpperCaseAttributes converts all attribute keys to upper-case.
	UpperCaseAttributes bool
	// DropEmptyValues removes attribute changes where both values are empty.
	DropEmptyValues bool
}

// DefaultTransformOptions returns a TransformOptions with no mutations applied.
func DefaultTransformOptions() TransformOptions {
	return TransformOptions{}
}

// Transformer applies deterministic, reversible mutations to a slice of
// DriftResult values. It is useful for normalising results before they are
// handed off to downstream consumers (formatters, exporters, etc.).
type Transformer struct {
	opts TransformOptions
}

// NewTransformer creates a Transformer with the supplied options.
func NewTransformer(opts TransformOptions) *Transformer {
	return &Transformer{opts: opts}
}

// Transform returns a new slice of DriftResult with all configured mutations
// applied. The original slice is never modified.
func (t *Transformer) Transform(results []drift.DriftResult) []drift.DriftResult {
	out := make([]drift.DriftResult, 0, len(results))
	for _, r := range results {
		out = append(out, t.transformOne(r))
	}
	return out
}

func (t *Transformer) transformOne(r drift.DriftResult) drift.DriftResult {
	copy := drift.DriftResult{
		ResourceID:   r.ResourceID,
		ResourceType: r.ResourceType,
		Changes:      make([]drift.AttributeChange, 0, len(r.Changes)),
	}
	if t.opts.PrefixResourceIDs != "" {
		copy.ResourceID = t.opts.PrefixResourceIDs + copy.ResourceID
	}
	if t.opts.SuffixResourceTypes != "" {
		copy.ResourceType = copy.ResourceType + t.opts.SuffixResourceTypes
	}
	for _, ch := range r.Changes {
		if t.opts.DropEmptyValues && ch.StateValue == "" && ch.LiveValue == "" {
			continue
		}
		key := ch.Attribute
		if t.opts.UpperCaseAttributes {
			key = strings.ToUpper(key)
		}
		copy.Changes = append(copy.Changes, drift.AttributeChange{
			Attribute:  key,
			StateValue: ch.StateValue,
			LiveValue:  ch.LiveValue,
			Kind:       ch.Kind,
		})
	}
	return copy
}
