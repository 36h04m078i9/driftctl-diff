package diff

import (
	"fmt"
	"sort"

	"github.com/owner/driftctl-diff/internal/drift"
)

// LinkerOptions controls how drift results are linked together.
type LinkerOptions struct {
	// LinkByAttribute groups results that share a common attribute name.
	LinkByAttribute string
	// MaxLinks caps the number of linked peers per result.
	MaxLinks int
}

// DefaultLinkerOptions returns sensible defaults.
func DefaultLinkerOptions() LinkerOptions {
	return LinkerOptions{
		LinkByAttribute: "",
		MaxLinks:        10,
	}
}

// Link represents a directed relationship between two drift results.
type Link struct {
	SourceID string
	TargetID string
	SharedAttribute string
}

// LinkerResult pairs the original results with the computed links.
type LinkerResult struct {
	Results []drift.ResourceDiff
	Links   []Link
}

// Linker finds relationships between drifted resources.
type Linker struct {
	opts LinkerOptions
}

// NewLinker creates a Linker with the provided options.
func NewLinker(opts LinkerOptions) *Linker {
	if opts.MaxLinks <= 0 {
		opts.MaxLinks = DefaultLinkerOptions().MaxLinks
	}
	return &Linker{opts: opts}
}

// Link analyses results and returns a LinkerResult with discovered links.
func (l *Linker) Link(results []drift.ResourceDiff) LinkerResult {
	if len(results) == 0 || l.opts.LinkByAttribute == "" {
		return LinkerResult{Results: results}
	}

	// index resource IDs by the value of the target attribute
	index := make(map[string][]string)
	for _, r := range results {
		for _, c := range r.Changes {
			if c.Attribute == l.opts.LinkByAttribute {
				key := fmt.Sprintf("%v", c.LiveValue)
				index[key] = append(index[key], r.ResourceID)
			}
		}
	}

	var links []Link
	for _, ids := range index {
		sort.Strings(ids)
		for i := 0; i < len(ids) && i < l.opts.MaxLinks; i++ {
			for j := i + 1; j < len(ids) && j <= i+l.opts.MaxLinks; j++ {
				links = append(links, Link{
					SourceID:        ids[i],
					TargetID:        ids[j],
					SharedAttribute: l.opts.LinkByAttribute,
				})
			}
		}
	}

	return LinkerResult{Results: results, Links: links}
}
