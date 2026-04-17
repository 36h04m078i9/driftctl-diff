package diff

import (
	"testing"

	"github.com/snyk/driftctl-diff/internal/drift"
)

func makeDriftResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "tags", Kind: drift.KindChanged, StateValue: "a", LiveValue: "b"},
				{Attribute: "acl", Kind: drift.KindMissing, StateValue: "private", LiveValue: ""},
			},
		},
		{
			ResourceType: "aws_instance",
			ResourceID:   "web-server",
			Changes: []drift.AttributeChange{
				{Attribute: "instance_type", Kind: drift.KindChanged, StateValue: "t2.micro", LiveValue: "t3.small"},
			},
		},
	}
}

func TestSearch_NoFilter_ReturnsAll(t *testing.T) {
	s := NewSearcher(SearchFilter{})
	results := makeDriftResults()
	out := s.Search(results)
	if len(out) != 0 {
		t.Fatalf("expected 0 results with empty filter (no changes matched), got %d", len(out))
	}
}

func TestSearch_ByResourceType(t *testing.T) {
	s := NewSearcher(SearchFilter{ResourceType: "aws_s3_bucket", Attribute: "tags"})
	out := s.Search(makeDriftResults())
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].ResourceType != "aws_s3_bucket" {
		t.Errorf("unexpected type %s", out[0].ResourceType)
	}
}

func TestSearch_ByResourceID_Partial(t *testing.T) {
	s := NewSearcher(SearchFilter{ResourceID: "web", Attribute: "instance_type"})
	out := s.Search(makeDriftResults())
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].ResourceID != "web-server" {
		t.Errorf("unexpected id %s", out[0].ResourceID)
	}
}

func TestSearch_ByKind_FiltersMismatched(t *testing.T) {
	s := NewSearcher(SearchFilter{KindSet: true, Kind: drift.KindMissing})
	out := s.Search(makeDriftResults())
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Changes[0].Kind != drift.KindMissing {
		t.Errorf("expected KindMissing, got %v", out[0].Changes[0].Kind)
	}
}

func TestSearch_ByAttribute_NarrowsChanges(t *testing.T) {
	s := NewSearcher(SearchFilter{Attribute: "acl"})
	out := s.Search(makeDriftResults())
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if len(out[0].Changes) != 1 || out[0].Changes[0].Attribute != "acl" {
		t.Errorf("unexpected changes: %+v", out[0].Changes)
	}
}
