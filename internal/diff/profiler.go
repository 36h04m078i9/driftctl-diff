package diff

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/owner/driftctl-diff/internal/drift"
)

// ProfileEntry records timing and size information for a single diff run.
type ProfileEntry struct {
	Label        string
	Duration     time.Duration
	ResourceCount int
	ChangeCount  int
}

// Profiler collects performance profiles across one or more diff operations.
type Profiler struct {
	entries []ProfileEntry
}

// NewProfiler returns an initialised Profiler.
func NewProfiler() *Profiler {
	return &Profiler{}
}

// Record appends a new ProfileEntry derived from the provided results and
// elapsed duration. label is a human-readable name for the operation.
func (p *Profiler) Record(label string, results []drift.ResourceDiff, d time.Duration) {
	changes := 0
	for _, r := range results {
		changes += len(r.Changes)
	}
	p.entries = append(p.entries, ProfileEntry{
		Label:        label,
		Duration:     d,
		ResourceCount: len(results),
		ChangeCount:  changes,
	})
}

// Entries returns a copy of all recorded entries sorted by duration descending.
func (p *Profiler) Entries() []ProfileEntry {
	out := make([]ProfileEntry, len(p.entries))
	copy(out, p.entries)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Duration > out[j].Duration
	})
	return out
}

// Print writes a human-readable profile table to w. If w is nil, os.Stdout is
// used.
func (p *Profiler) Print(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprintln(w, "=== Diff Profile ===")
	fmt.Fprintf(w, "%-30s %12s %10s %10s\n", "Label", "Duration", "Resources", "Changes")
	fmt.Fprintln(w, "--------------------------------------------------------------")
	for _, e := range p.Entries() {
		fmt.Fprintf(w, "%-30s %12s %10d %10d\n",
			e.Label, e.Duration.Round(time.Millisecond), e.ResourceCount, e.ChangeCount)
	}
}
