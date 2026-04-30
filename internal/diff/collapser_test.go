package diff

import (
	"testing"

	"github.com/driftctl-diff/internal/drift"
)

func makeCollapserResults(n int, changesEach int) []drift.ResourceDiff {
	results := make([]drift.ResourceDiff, n)
	for i := 0; i < n; i++ {
		changes := make([]drift.AttributeChange, changesEach)
		for j := 0; j < changesEach; j++ {
			changes[j] = drift.AttributeChange{
				Attribute: "attr",
				Expected:  "want",
				Actual:    "got",
				Kind:      drift.ChangeKindChanged,
			}
		}
		results[i] = drift.ResourceDiff{
			ResourceID:   "res",
			ResourceType: "aws_instance",
			Changes:      changes,
		}
	}
	return results
}

func TestCollapser_NoOptions_ReturnsAll(t *testing.T) {
	c := NewCollapser(DefaultCollapserOptions())
	input := makeCollapserResults(3, 4)
	out := c.Collapse(input)
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
	if len(out[0].Changes) != 4 {
		t.Fatalf("expected 4 changes, got %d", len(out[0].Changes))
	}
}

func TestCollapser_MaxChanges_CapsPerResource(t *testing.T) {
	opts := DefaultCollapserOptions()
	opts.MaxChangesPerResource = 2
	c := NewCollapser(opts)
	input := makeCollapserResults(2, 5)
	out := c.Collapse(input)
	for _, r := range out {
		if len(r.Changes) > 2 {
			t.Errorf("expected at most 2 changes, got %d", len(r.Changes))
		}
	}
}

func TestCollapser_OmitUnchanged_DropsEmpty(t *testing.T) {
	opts := DefaultCollapserOptions()
	opts.OmitUnchanged = true
	c := NewCollapser(opts)
	input := []drift.ResourceDiff{
		{ResourceID: "a", ResourceType: "aws_s3_bucket", Changes: nil},
		{ResourceID: "b", ResourceType: "aws_s3_bucket", Changes: []drift.AttributeChange{
			{Attribute: "tags", Expected: "x", Actual: "y", Kind: drift.ChangeKindChanged},
		}},
	}
	out := c.Collapse(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].ResourceID != "b" {
		t.Errorf("expected resource b, got %s", out[0].ResourceID)
	}
}

func TestCollapser_EmptyInput_ReturnsEmpty(t *testing.T) {
	c := NewCollapser(DefaultCollapserOptions())
	out := c.Collapse([]drift.ResourceDiff{})
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d", len(out))
	}
}

func TestCollapser_DoesNotMutateOriginal(t *testing.T) {
	opts := DefaultCollapserOptions()
	opts.MaxChangesPerResource = 1
	c := NewCollapser(opts)
	input := makeCollapserResults(1, 3)
	origLen := len(input[0].Changes)
	c.Collapse(input)
	if len(input[0].Changes) != origLen {
		t.Errorf("original slice was mutated")
	}
}
