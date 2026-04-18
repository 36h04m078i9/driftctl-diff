package diff

import (
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeMergeResult(resType, resID string, changes ...drift.AttributeChange) drift.DriftResult {
	return drift.DriftResult{
		ResourceType: resType,
		ResourceID:   resID,
		Changes:      changes,
	}
}

func makeChange(attr, want, got string) drift.AttributeChange {
	return drift.AttributeChange{Attribute: attr, Expected: want, Actual: got}
}

func TestMerge_EmptyInputs(t *testing.T) {
	m := NewMerger(MergeOptions{})
	out := m.Merge(nil, nil)
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d", len(out))
	}
}

func TestMerge_NoDuplicates(t *testing.T) {
	left := []drift.DriftResult{
		makeMergeResult("aws_s3_bucket", "bucket-a", makeChange("acl", "private", "public")),
	}
	right := []drift.DriftResult{
		makeMergeResult("aws_instance", "i-123", makeChange("ami", "ami-old", "ami-new")),
	}
	m := NewMerger(MergeOptions{})
	out := m.Merge(left, right)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestMerge_DeduplicatesSameAttribute(t *testing.T) {
	ch := makeChange("acl", "private", "public")
	left := []drift.DriftResult{makeMergeResult("aws_s3_bucket", "bucket-a", ch)}
	right := []drift.DriftResult{makeMergeResult("aws_s3_bucket", "bucket-a", ch)}

	m := NewMerger(MergeOptions{})
	out := m.Merge(left, right)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if len(out[0].Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(out[0].Changes))
	}
}

func TestMerge_PreferLeft_KeepsLeftValue(t *testing.T) {
	left := []drift.DriftResult{makeMergeResult("aws_s3_bucket", "b", makeChange("acl", "private", "LEFT"))}
	right := []drift.DriftResult{makeMergeResult("aws_s3_bucket", "b", makeChange("acl", "private", "RIGHT"))}

	m := NewMerger(MergeOptions{PreferLeft: true})
	out := m.Merge(left, right)
	if out[0].Changes[0].Actual != "LEFT" {
		t.Fatalf("expected LEFT, got %s", out[0].Changes[0].Actual)
	}
}

func TestMerge_PreferRight_KeepsRightValue(t *testing.T) {
	left := []drift.DriftResult{makeMergeResult("aws_s3_bucket", "b", makeChange("acl", "private", "LEFT"))}
	right := []drift.DriftResult{makeMergeResult("aws_s3_bucket", "b", makeChange("acl", "private", "RIGHT"))}

	m := NewMerger(MergeOptions{PreferLeft: false})
	out := m.Merge(left, right)
	if out[0].Changes[0].Actual != "RIGHT" {
		t.Fatalf("expected RIGHT, got %s", out[0].Changes[0].Actual)
	}
}

func TestMerge_SortedOutput(t *testing.T) {
	left := []drift.DriftResult{
		makeMergeResult("aws_s3_bucket", "z-bucket"),
		makeMergeResult("aws_instance", "i-999"),
	}
	m := NewMerger(MergeOptions{})
	out := m.Merge(left, nil)
	if out[0].ResourceType != "aws_instance" {
		t.Fatalf("expected aws_instance first, got %s", out[0].ResourceType)
	}
}
