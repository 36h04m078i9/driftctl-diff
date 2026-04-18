package diff

import (
	"strings"
	"testing"

	"github.com/driftctl-diff/internal/drift"
)

func makeTruncResults(n int) []drift.ResourceDiff {
	out := make([]drift.ResourceDiff, n)
	for i := range out {
		out[i] = drift.ResourceDiff{
			ResourceID:   "res-" + string(rune('a'+i)),
			ResourceType: "aws_instance",
			Changes: []drift.AttributeChange{
				{Attribute: "attr", Want: "old", Got: "new", Kind: drift.KindChanged},
			},
		}
	}
	return out
}

func TestTruncate_NoLimits_ReturnsAll(t *testing.T) {
	tr := NewTruncator(TruncateOptions{})
	results := makeTruncResults(5)
	out, truncated := tr.Truncate(results)
	if truncated {
		t.Fatal("expected no truncation")
	}
	if len(out) != 5 {
		t.Fatalf("expected 5, got %d", len(out))
	}
}

func TestTruncate_MaxResults_LimitsSlice(t *testing.T) {
	tr := NewTruncator(TruncateOptions{MaxResults: 3})
	out, truncated := tr.Truncate(makeTruncResults(6))
	if !truncated {
		t.Fatal("expected truncation")
	}
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
}

func TestTruncate_MaxChanges_LimitsChanges(t *testing.T) {
	results := []drift.ResourceDiff{{
		ResourceID:   "r1",
		ResourceType: "aws_s3_bucket",
		Changes: []drift.AttributeChange{
			{Attribute: "a", Want: "1", Got: "2", Kind: drift.KindChanged},
			{Attribute: "b", Want: "3", Got: "4", Kind: drift.KindChanged},
			{Attribute: "c", Want: "5", Got: "6", Kind: drift.KindChanged},
		},
	}}
	tr := NewTruncator(TruncateOptions{MaxChanges: 2})
	out, truncated := tr.Truncate(results)
	if !truncated {
		t.Fatal("expected truncation")
	}
	if len(out[0].Changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(out[0].Changes))
	}
}

func TestTruncate_MaxValueLen_TruncatesStrings(t *testing.T) {
	long := strings.Repeat("x", 200)
	results := []drift.ResourceDiff{{
		ResourceID:   "r1",
		ResourceType: "aws_iam_policy",
		Changes: []drift.AttributeChange{
			{Attribute: "policy", Want: long, Got: long, Kind: drift.KindChanged},
		},
	}}
	tr := NewTruncator(TruncateOptions{MaxValueLen: 50})
	out, truncated := tr.Truncate(results)
	if !truncated {
		t.Fatal("expected truncation")
	}
	if len(out[0].Changes[0].Got) != 53 { // 50 + len("...")
		t.Fatalf("unexpected length: %d", len(out[0].Changes[0].Got))
	}
	if !strings.HasSuffix(out[0].Changes[0].Got, "...") {
		t.Fatal("expected ellipsis suffix")
	}
}

func TestTruncate_DefaultOptions_NoPanic(t *testing.T) {
	tr := NewTruncator(DefaultTruncateOptions())
	_, _ = tr.Truncate(makeTruncResults(10))
}
