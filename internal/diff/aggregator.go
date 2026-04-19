package diff

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/owner/driftctl-diff/internal/drift"
)

// AggregateResult holds aggregated drift data grouped by resource type and kind.
type AggregateResult struct {
	ResourceType string
	Kind         drift.ChangeKind
	Count        int
	Attributes   []string
}

// Aggregator groups drift results by resource type and change kind.
type Aggregator struct {
	w io.Writer
}

// NewAggregator returns an Aggregator that writes to w, defaulting to stdout.
func NewAggregator(w io.Writer) *Aggregator {
	if w == nil {
		w = os.Stdout
	}
	return &Aggregator{w: w}
}

// Aggregate groups the provided drift results and returns a sorted slice of AggregateResult.
func (a *Aggregator) Aggregate(results []drift.ResourceDiff) []AggregateResult {
	type key struct {
		ResourceType string
		Kind         drift.ChangeKind
	}

	buckets := make(map[key]*AggregateResult)

	for _, r := range results {
		for _, c := range r.Changes {
			k := key{ResourceType: r.ResourceType, Kind: c.Kind}
			if _, ok := buckets[k]; !ok {
				buckets[k] = &AggregateResult{
					ResourceType: r.ResourceType,
					Kind:         c.Kind,
				}
			}
			buckets[k].Count++
			buckets[k].Attributes = append(buckets[k].Attributes, c.Attribute)
		}
	}

	out := make([]AggregateResult, 0, len(buckets))
	for _, v := range buckets {
		out = append(out, *v)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].ResourceType != out[j].ResourceType {
			return out[i].ResourceType < out[j].ResourceType
		}
		return string(out[i].Kind) < string(out[j].Kind)
	})
	return out
}

// Print writes a human-readable aggregation summary to the writer.
func (a *Aggregator) Print(results []drift.ResourceDiff) {
	agg := a.Aggregate(results)
	if len(agg) == 0 {
		fmt.Fprintln(a.w, "no drift detected")
		return
	}
	fmt.Fprintln(a.w, "drift aggregation:")
	for _, r := range agg {
		fmt.Fprintf(a.w, "  %-30s %-10s %d change(s)\n", r.ResourceType, r.Kind, r.Count)
	}
}
