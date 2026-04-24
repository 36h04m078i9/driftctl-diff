package diff

import (
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeReducerResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceID:   "bucket-1",
			ResourceType: "aws_s3_bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "acl", Got: "public", Want: "private", Kind: drift.KindChanged},
				{Attribute: "versioning", Got: "false", Want: "true", Kind: drift.KindChanged},
			},
		},
		{
			ResourceID:   "sg-1",
			ResourceType: "aws_security_group",
			Changes: []drift.AttributeChange{
				{Attribute: "ingress", Got: "", Want: "0.0.0.0/0", Kind: drift.KindMissing},
			},
		},
		{
			ResourceID:   "role-1",
			ResourceType: "aws_iam_role",
			Changes:      []drift.AttributeChange{},
		},
	}
}

func TestReducer_NoOptions_ReturnsAll(t *testing.T) {
	r := NewReducer(ReduceOptions{})
	out := r.Reduce(makeReducerResults())
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
}

func TestReducer_KeepTopN_LimitsCount(t *testing.T) {
	r := NewReducer(ReduceOptions{KeepTopN: 1})
	out := r.Reduce(makeReducerResults())
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	// The resource with the most changes should be first.
	if out[0].ResourceID != "bucket-1" {
		t.Errorf("expected bucket-1, got %s", out[0].ResourceID)
	}
}

func TestReducer_OnlyKinds_FiltersUnmatched(t *testing.T) {
	r := NewReducer(ReduceOptions{OnlyKinds: []drift.ChangeKind{drift.KindMissing}})
	out := r.Reduce(makeReducerResults())
	for _, res := range out {
		for _, ch := range res.Changes {
			if ch.Kind != drift.KindMissing {
				t.Errorf("unexpected kind %v in result", ch.Kind)
			}
		}
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 result after kind filter, got %d", len(out))
	}
}

func TestReducer_KeepTopN_LargerThanSlice_ReturnsAll(t *testing.T) {
	r := NewReducer(ReduceOptions{KeepTopN: 100})
	out := r.Reduce(makeReducerResults())
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
}

func TestReducer_EmptyInput_ReturnsEmpty(t *testing.T) {
	r := NewReducer(ReduceOptions{KeepTopN: 5})
	out := r.Reduce(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %d", len(out))
	}
}

func TestReducer_SortedDescendingByChangeCount(t *testing.T) {
	r := NewReducer(ReduceOptions{})
	out := r.Reduce(makeReducerResults())
	for i := 1; i < len(out); i++ {
		if len(out[i-1].Changes) < len(out[i].Changes) {
			t.Errorf("results not sorted descending at index %d", i)
		}
	}
}
