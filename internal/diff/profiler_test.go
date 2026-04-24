package diff

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/owner/driftctl-diff/internal/drift"
)

func makeProfilerResults(n int) []drift.ResourceDiff {
	results := make([]drift.ResourceDiff, n)
	for i := 0; i < n; i++ {
		results[i] = drift.ResourceDiff{
			ResourceID:   fmt.Sprintf("res-%d", i),
			ResourceType: "aws_instance",
			Changes: []drift.AttributeChange{
				{Attribute: "ami", Got: "old", Want: "new"},
			},
		}
	}
	return results
}

func TestProfiler_StartsEmpty(t *testing.T) {
	p := NewProfiler()
	if len(p.Entries()) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(p.Entries()))
	}
}

func TestProfiler_RecordIncreasesCount(t *testing.T) {
	p := NewProfiler()
	p.Record("first", makeProfilerResults(3), 10*time.Millisecond)
	if len(p.Entries()) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(p.Entries()))
	}
}

func TestProfiler_EntriesSortedByDurationDesc(t *testing.T) {
	p := NewProfiler()
	p.Record("fast", nil, 5*time.Millisecond)
	p.Record("slow", nil, 500*time.Millisecond)
	p.Record("medium", nil, 50*time.Millisecond)

	entries := p.Entries()
	if entries[0].Label != "slow" {
		t.Errorf("expected first entry to be 'slow', got %q", entries[0].Label)
	}
	if entries[2].Label != "fast" {
		t.Errorf("expected last entry to be 'fast', got %q", entries[2].Label)
	}
}

func TestProfiler_ChangeCountSumsChanges(t *testing.T) {
	p := NewProfiler()
	p.Record("op", makeProfilerResults(4), 20*time.Millisecond)

	entries := p.Entries()
	if entries[0].ChangeCount != 4 {
		t.Errorf("expected 4 changes, got %d", entries[0].ChangeCount)
	}
}

func TestProfiler_Print_ContainsLabel(t *testing.T) {
	p := NewProfiler()
	p.Record("my-operation", makeProfilerResults(2), 15*time.Millisecond)

	var buf bytes.Buffer
	p.Print(&buf)

	if !strings.Contains(buf.String(), "my-operation") {
		t.Errorf("expected output to contain label, got:\n%s", buf.String())
	}
}

func TestProfiler_Print_NilWriter_DefaultsToStdout(t *testing.T) {
	p := NewProfiler()
	p.Record("noop", nil, time.Millisecond)
	// Should not panic.
	p.Print(nil)
}
