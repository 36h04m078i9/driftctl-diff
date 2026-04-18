package diff

import (
	"testing"

	"github.com/nikoksr/driftctl-diff/internal/drift"
)

func makeSplitterResults() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceID:   "aws_s3_bucket.added",
			ResourceType: "aws_s3_bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "tags", Kind: drift.KindAdded},
			},
		},
		{
			ResourceID:   "aws_s3_bucket.deleted",
			ResourceType: "aws_s3_bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "region", Kind: drift.KindDeleted},
			},
		},
		{
			ResourceID:   "aws_instance.changed",
			ResourceType: "aws_instance",
			Changes: []drift.AttributeChange{
				{Attribute: "ami", Kind: drift.KindChanged, WantValue: "old", GotValue: "new"},
			},
		},
		{
			ResourceID:   "aws_instance.mixed",
			ResourceType: "aws_instance",
			Changes: []drift.AttributeChange{
				{Attribute: "ami", Kind: drift.KindAdded},
				{Attribute: "type", Kind: drift.KindDeleted},
			},
		},
	}
}

func TestSplitter_AddedBucket(t *testing.T) {
	s := NewSplitter(SplitOptions{})
	res := s.Split(makeSplitterResults())
	if len(res.Added) != 1 || res.Added[0].ResourceID != "aws_s3_bucket.added" {
		t.Fatalf("expected 1 added resource, got %v", res.Added)
	}
}

func TestSplitter_DeletedBucket(t *testing.T) {
	s := NewSplitter(SplitOptions{})
	res := s.Split(makeSplitterResults())
	if len(res.Deleted) != 1 || res.Deleted[0].ResourceID != "aws_s3_bucket.deleted" {
		t.Fatalf("expected 1 deleted resource, got %v", res.Deleted)
	}
}

func TestSplitter_ChangedBucket(t *testing.T) {
	s := NewSplitter(SplitOptions{})
	res := s.Split(makeSplitterResults())
	// changed + mixed both land in Changed
	if len(res.Changed) != 2 {
		t.Fatalf("expected 2 changed resources, got %d", len(res.Changed))
	}
}

func TestSplitter_EmptyInput(t *testing.T) {
	s := NewSplitter(SplitOptions{})
	res := s.Split(nil)
	if len(res.Added)+len(res.Deleted)+len(res.Changed) != 0 {
		t.Fatal("expected all buckets empty")
	}
}

func TestSplitter_NoChanges_GoesToChanged(t *testing.T) {
	s := NewSplitter(SplitOptions{})
	input := []drift.ResourceDiff{{ResourceID: "r", ResourceType: "t", Changes: nil}}
	res := s.Split(input)
	if len(res.Changed) != 1 {
		t.Fatalf("expected resource with no changes in Changed bucket")
	}
}
