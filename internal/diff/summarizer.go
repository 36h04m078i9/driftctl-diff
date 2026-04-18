package diff

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/owner/driftctl-diff/internal/drift"
)

// SummaryOptions controls what the Summarizer includes in its output.
type SummaryOptions struct {
	TopN int // number of top drifted resource types to show; 0 = all
}

// TypeSummary holds aggregated drift info for a single resource type.
type TypeSummary struct {
	ResourceType   string
	ResourceCount  int
	ChangedAttrs   int
}

// Summarizer aggregates drift results by resource type.
type Summarizer struct {
	opts SummaryOptions
}

// NewSummarizer returns a Summarizer with the given options.
func NewSummarizer(opts SummaryOptions) *Summarizer {
	if opts.TopN < 0 {
		opts.TopN = 0
	}
	return &Summarizer{opts: opts}
}

// Summarize computes per-type summaries from drift results.
func (s *Summarizer) Summarize(results []drift.ResourceDiff) []TypeSummary {
	index := map[string]*TypeSummary{}
	for _, r := range results {
		ts, ok := index[r.ResourceType]
		if !ok {
			ts = &TypeSummary{ResourceType: r.ResourceType}
			index[r.ResourceType] = ts
		}
		ts.ResourceCount++
		ts.ChangedAttrs += len(r.Changes)
	}
	out := make([]TypeSummary, 0, len(index))
	for _, ts := range index {
		out = append(out, *ts)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].ChangedAttrs != out[j].ChangedAttrs {
			return out[i].ChangedAttrs > out[j].ChangedAttrs
		}
		return out[i].ResourceType < out[j].ResourceType
	})
	if s.opts.TopN > 0 && len(out) > s.opts.TopN {
		out = out[:s.opts.TopN]
	}
	return out
}

// Print writes a human-readable summary table to w (defaults to stdout).
func (s *Summarizer) Print(w io.Writer, summaries []TypeSummary) {
	if w == nil {
		w = os.Stdout
	}
	if len(summaries) == 0 {
		fmt.Fprintln(w, "No drift detected.")
		return
	}
	fmt.Fprintf(w, "%-35s %10s %15s\n", "RESOURCE TYPE", "RESOURCES", "CHANGED ATTRS")
	fmt.Fprintf(w, "%-35s %10s %15s\n", "-------------", "---------", "-------------")
	for _, ts := range summaries {
		fmt.Fprintf(w, "%-35s %10d %15d\n", ts.ResourceType, ts.ResourceCount, ts.ChangedAttrs)
	}
}
