package diff

import (
	"fmt"

	"github.com/acme/driftctl-diff/internal/drift"
)

// ClonerOptions controls deep-copy behaviour.
type ClonerOptions struct {
	// CopyMetadata controls whether the Metadata map is deep-copied.
	// When false the cloned result shares the original map reference.
	CopyMetadata bool
}

// DefaultClonerOptions returns sensible defaults.
func DefaultClonerOptions() ClonerOptions {
	return ClonerOptions{
		CopyMetadata: true,
	}
}

// Cloner produces an independent deep copy of a slice of DriftResults so that
// downstream pipeline stages cannot accidentally mutate earlier stages.
type Cloner struct {
	opts ClonerOptions
}

// NewCloner constructs a Cloner with the supplied options.
func NewCloner(opts ClonerOptions) *Cloner {
	return &Cloner{opts: opts}
}

// Clone returns a deep copy of results. An error is returned only when a
// result contains a nil ResourceID, which is treated as invalid input.
func (c *Cloner) Clone(results []drift.DriftResult) ([]drift.DriftResult, error) {
	out := make([]drift.DriftResult, 0, len(results))
	for i, r := range results {
		if r.ResourceID == "" {
			return nil, fmt.Errorf("cloner: result at index %d has empty ResourceID", i)
		}

		cloned := drift.DriftResult{
			ResourceID:   r.ResourceID,
			ResourceType: r.ResourceType,
			Changes:      cloneChanges(r.Changes),
		}

		if c.opts.CopyMetadata && r.Metadata != nil {
			m := make(map[string]string, len(r.Metadata))
			for k, v := range r.Metadata {
				m[k] = v
			}
			cloned.Metadata = m
		} else {
			cloned.Metadata = r.Metadata
		}

		out = append(out, cloned)
	}
	return out, nil
}

func cloneChanges(changes []drift.AttributeChange) []drift.AttributeChange {
	if changes == nil {
		return nil
	}
	copied := make([]drift.AttributeChange, len(changes))
	copy(copied, changes)
	return copied
}
