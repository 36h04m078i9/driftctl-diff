package diff

import (
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeFilterResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "acl", Got: "public", Want: "private", Kind: drift.KindChanged},
			},
		},
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-abc123",
			Changes: []drift.AttributeChange{
				{Attribute: "instance_type", Got: "", Want: "t3.micro", Kind: drift.KindMissing},
				{Attribute: "ami", Got: "ami-new", Want: "ami-old", Kind: drift.KindChanged},
			},
		},
		{
			ResourceType: "aws_security_group",
			ResourceID:   "sg-xyz",
			Changes:      []drift.AttributeChange{},
		},
	}
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	f := NewFilter(FilterOptions{})
	results := makeFilterResults()
	got := f.Apply(results)
	if len(got) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(got))
	}
}

func TestFilter_ByResourceType(t *testing.T) {
	f := NewFilter(FilterOptions{ResourceType: "aws_instance"})
	got := f.Apply(makeFilterResults())
	if len(got) != 1 || got[0].ResourceType != "aws_instance" {
		t.Fatalf("unexpected results: %+v", got)
	}
}

func TestFilter_ByResourceID_Partial(t *testing.T) {
	f := NewFilter(FilterOptions{ResourceID: "bucket"})
	got := f.Apply(makeFilterResults())
	if len(got) != 1 || got[0].ResourceID != "my-bucket" {
		t.Fatalf("unexpected results: %+v", got)
	}
}

func TestFilter_ByKind_Missing(t *testing.T) {
	f := NewFilter(FilterOptions{Kinds: []drift.ChangeKind{drift.KindMissing}})
	got := f.Apply(makeFilterResults())
	if len(got) != 1 || got[0].ResourceID != "i-abc123" {
		t.Fatalf("unexpected results: %+v", got)
	}
}

func TestFilter_MinChanges(t *testing.T) {
	f := NewFilter(FilterOptions{MinChanges: 2})
	got := f.Apply(makeFilterResults())
	if len(got) != 1 || got[0].ResourceID != "i-abc123" {
		t.Fatalf("expected only i-abc123, got %+v", got)
	}
}

func TestFilter_CombinedOptions(t *testing.T) {
	f := NewFilter(FilterOptions{
		ResourceType: "aws_instance",
		Kinds:        []drift.ChangeKind{drift.KindChanged},
	})
	got := f.Apply(makeFilterResults())
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
}
