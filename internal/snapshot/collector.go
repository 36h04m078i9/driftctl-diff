package snapshot

import (
	"fmt"

	"github.com/example/driftctl-diff/internal/provider"
	"github.com/example/driftctl-diff/internal/state"
)

// Collector fetches live attributes for every resource in a parsed state
// and assembles them into a Snapshot.
type Collector struct {
	registry *provider.Registry
}

// NewCollector returns a Collector backed by the given provider registry.
func NewCollector(registry *provider.Registry) *Collector {
	return &Collector{registry: registry}
}

// Collect iterates over state resources, fetches live attributes via the
// registry, and returns a populated Snapshot. Errors for individual
// resources are collected and returned as a combined error.
func (c *Collector) Collect(resources []state.Resource) (*Snapshot, error) {
	snap := New()
	var errs []string

	for _, res := range resources {
		attrs, err := c.registry.FetchAttributes(res.Type, res.ID)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s.%s: %v", res.Type, res.ID, err))
			continue
		}
		snap.Add(resourceKey(res), attrs)
	}

	if len(errs) > 0 {
		return snap, fmt.Errorf("collector: %d resource(s) failed: %v", len(errs), errs)
	}
	return snap, nil
}

func resourceKey(r state.Resource) string {
	return r.Type + "." + r.Name
}
