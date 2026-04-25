// Package diff provides utilities for processing and transforming drift results.
package diff

import "github.com/acme/driftctl-diff/internal/drift"

// CompactOptions controls the behaviour of the Compactor.
type CompactOptions struct {
	// MergeAdjacentResources combines consecutive DriftResult entries that share
	// the same ResourceID and ResourceType into a single entry.
	MergeAdjacentResources bool

	// DeduplicateChanges removes duplicate Change entries within a single result.
	DeduplicateChanges bool

	// DropEmptyResults removes results that contain no changes.
	DropEmptyResults bool
}

// DefaultCompactOptions returns a CompactOptions with all flags disabled so
// that the Compactor acts as a pass-through by default.
func DefaultCompactOptions() CompactOptions {
	return CompactOptions{}
}

// Compactor applies a configurable set of compaction passes to a slice of
// DriftResult values, reducing noise and merging related entries.
type Compactor struct {
	opts CompactOptions
}

// NewCompactor creates a Compactor with the supplied options.
func NewCompactor(opts CompactOptions) *Compactor {
	return &Compactor{opts: opts}
}

// Compact runs all enabled compaction passes and returns the resulting slice.
func (c *Compactor) Compact(results []drift.DriftResult) []drift.DriftResult {
	if len(results) == 0 {
		return []drift.DriftResult{}
	}

	out := make([]drift.DriftResult, len(results))
	copy(out, results)

	if c.opts.DeduplicateChanges {
		out = deduplicateChangesInResults(out)
	}

	if c.opts.MergeAdjacentResources {
		out = mergeAdjacentChanges(out)
	}

	if c.opts.DropEmptyResults {
		out = dropEmpty(out)
	}

	return out
}

// mergeAdjacentChanges combines consecutive results that share the same
// ResourceID and ResourceType.
func mergeAdjacentChanges(results []drift.DriftResult) []drift.DriftResult {
	if len(results) == 0 {
		return results
	}

	merged := []drift.DriftResult{results[0]}

	for i := 1; i < len(results); i++ {
		cur := results[i]
		last := &merged[len(merged)-1]

		if cur.ResourceID == last.ResourceID && cur.ResourceType == last.ResourceType {
			last.Changes = append(last.Changes, cur.Changes...)
		} else {
			merged = append(merged, cur)
		}
	}

	return merged
}

// deduplicateChangesInResults removes duplicate Change entries within each
// DriftResult, comparing by Attribute, Got, and Want.
func deduplicateChangesInResults(results []drift.DriftResult) []drift.DriftResult {
	for i := range results {
		seen := make(map[string]struct{})
		uniq := results[i].Changes[:0]

		for _, ch := range results[i].Changes {
			key := ch.Attribute + "|" + ch.Got + "|" + ch.Want
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			uniq = append(uniq, ch)
		}

		results[i].Changes = uniq
	}

	return results
}

// dropEmpty filters out results that carry no changes.
func dropEmpty(results []drift.DriftResult) []drift.DriftResult {
	out := results[:0]
	for _, r := range results {
		if len(r.Changes) > 0 {
			out = append(out, r)
		}
	}
	return out
}
