package diff

import (
	"testing"
	"time"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeTimelineDiffs(id string) []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{ResourceID: id, ResourceType: "aws_instance"},
	}
}

func TestTimeline_StartsEmpty(t *testing.T) {
	tl := NewTimeline()
	if tl.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", tl.Len())
	}
}

func TestTimeline_AddIncreasesLen(t *testing.T) {
	tl := NewTimeline()
	tl.Add(time.Now(), makeTimelineDiffs("i-1"))
	if tl.Len() != 1 {
		t.Fatalf("expected 1, got %d", tl.Len())
	}
}

func TestTimeline_EntriesAreChronological(t *testing.T) {
	tl := NewTimeline()
	now := time.Now()
	tl.Add(now.Add(2*time.Hour), makeTimelineDiffs("i-3"))
	tl.Add(now, makeTimelineDiffs("i-1"))
	tl.Add(now.Add(time.Hour), makeTimelineDiffs("i-2"))

	entries := tl.Entries()
	for i := 1; i < len(entries); i++ {
		if entries[i].CapturedAt.Before(entries[i-1].CapturedAt) {
			t.Errorf("entry %d is before entry %d", i, i-1)
		}
	}
}

func TestTimeline_Latest_Empty(t *testing.T) {
	tl := NewTimeline()
	_, ok := tl.Latest()
	if ok {
		t.Fatal("expected false for empty timeline")
	}
}

func TestTimeline_Latest_ReturnsMostRecent(t *testing.T) {
	tl := NewTimeline()
	now := time.Now()
	tl.Add(now, makeTimelineDiffs("i-1"))
	tl.Add(now.Add(time.Hour), makeTimelineDiffs("i-2"))

	entry, ok := tl.Latest()
	if !ok {
		t.Fatal("expected entry")
	}
	if entry.Results[0].ResourceID != "i-2" {
		t.Errorf("expected i-2, got %s", entry.Results[0].ResourceID)
	}
}

func TestTimeline_Entries_ReturnsCopy(t *testing.T) {
	tl := NewTimeline()
	tl.Add(time.Now(), makeTimelineDiffs("i-1"))
	e1 := tl.Entries()
	e1[0].ResourceID = "mutated"
	e2 := tl.Entries()
	if e2[0].CapturedAt.IsZero() {
		t.Error("entry should not be zeroed")
	}
}
