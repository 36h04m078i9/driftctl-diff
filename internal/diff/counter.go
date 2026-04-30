package diff

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/acme/driftctl-diff/internal/drift"
)

// CounterEntry holds the drift count for a single resource type.
type CounterEntry struct {
	ResourceType string
	Total        int
	Changed      int
	Added        int
	Deleted      int
}

// Counter aggregates change counts grouped by resource type.
type Counter struct {
	entries map[string]*CounterEntry
}

// NewCounter returns an initialised Counter.
func NewCounter() *Counter {
	return &Counter{entries: make(map[string]*CounterEntry)}
}

// Count processes a slice of DriftResult and populates the counter.
func (c *Counter) Count(results []drift.DriftResult) {
	for _, r := range results {
		e, ok := c.entries[r.ResourceType]
		if !ok {
			e = &CounterEntry{ResourceType: r.ResourceType}
			c.entries[r.ResourceType] = e
		}
		e.Total++
		for _, ch := range r.Changes {
			switch ch.Kind {
			case drift.KindChanged:
				e.Changed++
			case drift.KindAdded:
				e.Added++
			case drift.KindDeleted:
				e.Deleted++
			}
		}
	}
}

// Entries returns a stable, sorted slice of CounterEntry values.
func (c *Counter) Entries() []CounterEntry {
	out := make([]CounterEntry, 0, len(c.entries))
	for _, e := range c.entries {
		out = append(out, *e)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ResourceType < out[j].ResourceType
	})
	return out
}

// Print writes a summary table to w (defaults to os.Stdout when nil).
func (c *Counter) Print(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	entries := c.Entries()
	if len(entries) == 0 {
		fmt.Fprintln(w, "no drift detected")
		return
	}
	fmt.Fprintf(w, "%-30s %8s %8s %8s %8s\n", "RESOURCE TYPE", "TOTAL", "CHANGED", "ADDED", "DELETED")
	for _, e := range entries {
		fmt.Fprintf(w, "%-30s %8d %8d %8d %8d\n", e.ResourceType, e.Total, e.Changed, e.Added, e.Deleted)
	}
}
