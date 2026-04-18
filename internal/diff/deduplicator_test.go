package diff

import (
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeDedupResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "region", Kind: drift.KindChanged, StateValue: "us-east-1", LiveValue: "eu-west-1"},
				{Attribute: "region", Kind: drift.KindChanged, StateValue: "us-east-1", LiveValue: "eu-west-1"},
			},
		},
	}
}

func TestDeduplicator_RemovesDuplicateChanges(t *testing.T) {
	d := NewDeduplicator(DeduplicateOptions{})
	results := makeDedupResults()
	out := d.Deduplicate(results)

	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if len(out[0].Changes) != 1 {
		t.Fatalf("expected 1 change after dedup, got %d", len(out[0].Changes))
	}
}

func TestDeduplicator_KeepsDistinctChanges(t *testing.T) {
	d := NewDeduplicator(DeduplicateOptions{})
	results := []drift.DriftResult{
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "region", Kind: drift.KindChanged},
				{Attribute: "acl", Kind: drift.KindChanged},
			},
		},
	}
	out := d.Deduplicate(results)
	if len(out[0].Changes) != 2 {
		t.Fatalf("expected 2 distinct changes, got %d", len(out[0].Changes))
	}
}

func TestDeduplicator_IgnoreCase_DeduplicatesMixedCase(t *testing.T) {
	d := NewDeduplicator(DeduplicateOptions{IgnoreCase: true})
	results := []drift.DriftResult{
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-123",
			Changes: []drift.AttributeChange{
				{Attribute: "Tags", Kind: drift.KindChanged},
				{Attribute: "tags", Kind: drift.KindChanged},
			},
		},
	}
	out := d.Deduplicate(results)
	if len(out[0].Changes) != 1 {
		t.Fatalf("expected 1 change with IgnoreCase, got %d", len(out[0].Changes))
	}
}

func TestDeduplicator_EmptyInput_ReturnsEmpty(t *testing.T) {
	d := NewDeduplicator(DeduplicateOptions{})
	out := d.Deduplicate(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d", len(out))
	}
}

func TestDeduplicator_DropsResultWithNoChangesAfterDedup(t *testing.T) {
	d := NewDeduplicator(DeduplicateOptions{})
	results := []drift.DriftResult{
		{
			ResourceType: "aws_vpc",
			ResourceID:   "vpc-1",
			Changes: []drift.AttributeChange{
				{Attribute: "cidr", Kind: drift.KindChanged},
				{Attribute: "cidr", Kind: drift.KindChanged},
			},
		},
	}
	out := d.Deduplicate(results)
	if len(out) != 1 {
		t.Fatalf("expected result to be kept with 1 change, got %d results", len(out))
	}
}
