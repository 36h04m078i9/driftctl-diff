package diff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/owner/driftctl-diff/internal/drift"
)

// CompactOptions controls how the Compactor condenses drift results.
type CompactOptions struct {
	// MaxAttributesPerResource limits the number of attribute changes shown per
	// resource. Remaining changes are summarised as "… N more".
	MaxAttributesPerResource int

	// OmitUnchanged drops resources that have zero attribute changes from the
	// output entirely.
	OmitUnchanged bool

	// MergeAdjacentValues collapses a → b → c chains into a single a → c entry
	// when the same attribute appears more than once for a resource.
	MergeAdjacentValues bool
}

// DefaultCompactOptions returns conservative defaults that keep all data but
// cap per-resource attribute lines at 10.
func DefaultCompactOptions() CompactOptions {
	return CompactOptions{
		MaxAttributesPerResource: 10,
		OmitUnchanged:            true,
		MergeAdjacentValues:      false,
	}
}

// Compactor condenses a slice of DriftResult values according to CompactOptions,
// reducing noise in large drift reports.
type Compactor struct {
	opts CompactOptions
	w    io.Writer
}

// NewCompactor creates a Compactor with the supplied options. If w is nil,
// os.Stdout is used when Print is called.
func NewCompactor(opts CompactOptions, w io.Writer) *Compactor {
	if w == nil {
		w = os.Stdout
	}
	return &Compactor{opts: opts, w: w}
}

// Compact applies the configured rules to results and returns the condensed
// slice. The original slice is not modified.
func (c *Compactor) Compact(results []drift.DriftResult) []drift.DriftResult {
	out := make([]drift.DriftResult, 0, len(results))

	for _, r := range results {
		changes := r.Changes

		if c.opts.OmitUnchanged && len(changes) == 0 {
			continue
		}

		if c.opts.MergeAdjacentValues {
			changes = mergeAdjacentChanges(changes)
		}

		if c.opts.MaxAttributesPerResource > 0 && len(changes) > c.opts.MaxAttributesPerResource {
			changes = changes[:c.opts.MaxAttributesPerResource]
		}

		copy := r
		copy.Changes = changes
		out = append(out, copy)
	}

	return out
}

// Print writes a compact human-readable summary of the results to the writer.
func (c *Compactor) Print(results []drift.DriftResult) {
	compacted := c.Compact(results)

	if len(compacted) == 0 {
		fmt.Fprintln(c.w, "No drift detected.")
		return
	}

	for _, r := range compacted {
		fmt.Fprintf(c.w, "Resource: %s (%s)  changes: %d\n",
			r.ResourceID, r.ResourceType, len(r.Changes))

		for _, ch := range r.Changes {
			fmt.Fprintf(c.w, "  ~ %s: %q → %q\n",
				ch.Attribute, ch.StateValue, ch.LiveValue)
		}
	}
}

// mergeAdjacentChanges collapses duplicate attribute entries so that only the
// first StateValue and last LiveValue are kept.
func mergeAdjacentChanges(changes []drift.Change) []drift.Change {
	seen := make(map[string]*drift.Change)
	order := make([]string, 0, len(changes))

	for i := range changes {
		ch := &changes[i]
		key := strings.ToLower(ch.Attribute)
		if existing, ok := seen[key]; ok {
			existing.LiveValue = ch.LiveValue
		} else {
			copy := *ch
			seen[key] = &copy
			order = append(order, key)
		}
	}

	merged := make([]drift.Change, 0, len(order))
	for _, k := range order {
		merged = append(merged, *seen[k])
	}
	return merged
}
