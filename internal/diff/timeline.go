package diff

import (
	"sort"
	"time"

	"github.com/acme/driftctl-diff/internal/drift"
)

// TimelineEntry records drift results captured at a specific point in time.
type TimelineEntry struct {
	CapturedAt time.Time
	Results    []drift.ResourceDiff
}

// Timeline holds an ordered sequence of drift snapshots.
type Timeline struct {
	entries []TimelineEntry
}

// NewTimeline creates an empty Timeline.
func NewTimeline() *Timeline {
	return &Timeline{}
}

// Add appends a new entry to the timeline.
func (t *Timeline) Add(capturedAt time.Time, results []drift.ResourceDiff) {
	t.entries = append(t.entries, TimelineEntry{
		CapturedAt: capturedAt,
		Results:    results,
	})
	sort.Slice(t.entries, func(i, j int) bool {
		return t.entries[i].CapturedAt.Before(t.entries[j].CapturedAt)
	})
}

// Entries returns all timeline entries in chronological order.
func (t *Timeline) Entries() []TimelineEntry {
	out := make([]TimelineEntry, len(t.entries))
	copy(out, t.entries)
	return out
}

// Latest returns the most recent entry, or false if empty.
func (t *Timeline) Latest() (TimelineEntry, bool) {
	if len(t.entries) == 0 {
		return TimelineEntry{}, false
	}
	return t.entries[len(t.entries)-1], true
}

// Len returns the number of entries.
func (t *Timeline) Len() int {
	return len(t.entries)
}

// Since returns all entries captured at or after the given time.
func (t *Timeline) Since(since time.Time) []TimelineEntry {
	var out []TimelineEntry
	for _, e := range t.entries {
		if !e.CapturedAt.Before(since) {
			out = append(out, e)
		}
	}
	return out
}
