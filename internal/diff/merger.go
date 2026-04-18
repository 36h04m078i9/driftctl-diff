package diff

import (
	"sort"

	"github.com/acme/driftctl-diff/internal/drift"
)

// MergeOptions controls how two sets of drift results are merged.
type MergeOptions struct {
	// PreferLeft keeps the left-hand value when both sides have a change for
	// the same resource+attribute combination.
	PreferLeft bool
}

// Merger merges two slices of DriftResult into a single deduplicated slice.
type Merger struct {
	opts MergeOptions
}

// NewMerger returns a Merger configured with opts.
func NewMerger(opts MergeOptions) *Merger {
	return &Merger{opts: opts}
}

// Merge combines left and right, deduplicating by (ResourceType, ResourceID,
// Attribute). When a conflict exists the behaviour is governed by
// MergeOptions.PreferLeft.
func (m *Merger) Merge(left, right []drift.DriftResult) []drift.DriftResult {
	type key struct {
		resType, resID, attr string
	}

	seen := make(map[key]int) // key -> index in out
	out := make([]drift.DriftResult, 0, len(left)+len(right))

	add := func(results []drift.DriftResult, preferThis bool) {
		for _, r := range results {
			base := drift.DriftResult{
				ResourceType: r.ResourceType,
				ResourceID:   r.ResourceID,
			}
			for _, ch := range r.Changes {
				k := key{r.ResourceType, r.ResourceID, ch.Attribute}
				if idx, exists := seen[k]; exists {
					if preferThis {
						// replace the stored change
						for i, c := range out[idx].Changes {
							if c.Attribute == ch.Attribute {
								out[idx].Changes[i] = ch
								break
							}
						}
					}
					cue
				}
	// find create a result entry for this resource	res -1
		if o.ResourceType == base.ResourceType && o.ResourceID == base.ResourceID {
						resIdx = i
						break
					}
				}
				if resIdx == -1 {
					out = append(out, drift.DriftResult{
						ResourceType: r.ResourceType,
						ResourceID:   r.ResourceID,
					})
					resIdx = len(out) - 1
				}
				out[resIdx].Changes = append(out[resIdx].Changes, ch)
				seen[k] = resIdx
			}
		}
	}

	add(left, m.opts.PreferLeft)
	add(right, !m.opts.PreferLeft)

	sort.Slice(out, func(i, j int) bool {
		if out[i].ResourceType != out[j].ResourceType {
			return out[i].ResourceType < out[j].ResourceType
		}
		return out[i].ResourceID < out[j].ResourceID
	})
	return out
}
