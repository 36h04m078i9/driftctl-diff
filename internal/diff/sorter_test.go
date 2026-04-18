package diff

import (
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeSortDriftResults() []drift.DriftResult {
	return []drift.DriftResult{
		{ResourceType: "aws_s3_bucket", ResourceID: "bucket-b", Changes: []drift.AttributeChange{{Attribute: "a"}}},
		{ResourceType: "aws_instance", ResourceID: "instance-a", Changes: []drift.AttributeChange{{Attribute: "a"}, {Attribute: "b"}}},
		{ResourceType: "aws_vpc", ResourceID: "vpc-c", Changes: []drift.AttributeChange{}},
	}
}

func TestSort_ByResourceType_Ascending(t *testing.T) {
	s := NewSorter(SortByResourceType, Ascending)
	results := s.Sort(makeSortDriftResults())

	if results[0].ResourceType != "aws_instance" {
		t.Errorf("expected aws_instance first, got %s", results[0].ResourceType)
	}
	if results[2].ResourceType != "aws_vpc" {
		t.Errorf("expected aws_vpc last, got %s", results[2].ResourceType)
	}
}

func TestSort_ByResourceType_Descending(t *testing.T) {
	s := NewSorter(SortByResourceType, Descending)
	results := s.Sort(makeSortDriftResults())

	if results[0].ResourceType != "aws_vpc" {
		t.Errorf("expected aws_vpc first, got %s", results[0].ResourceType)
	}
}

func TestSort_ByResourceID_Ascending(t *testing.T) {
	s := NewSorter(SortByResourceID, Ascending)
	results := s.Sort(makeSortDriftResults())

	if results[0].ResourceID != "bucket-b" {
		t.Errorf("expected bucket-b first, got %s", results[0].ResourceID)
	}
}

func TestSort_ByChangeCount_Ascending(t *testing.T) {
	s := NewSorter(SortByChangeCount, Ascending)
	results := s.Sort(makeSortDriftResults())

	if len(results[0].Changes) != 0 {
		t.Errorf("expected 0 changes first, got %d", len(results[0].Changes))
	}
	if len(results[2].Changes) != 2 {
		t.Errorf("expected 2 changes last, got %d", len(results[2].Changes))
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	original := makeSortDriftResults()
	firstID := original[0].ResourceID
	s := NewSorter(SortByResourceID, Ascending)
	s.Sort(original)

	if original[0].ResourceID != firstID {
		t.Error("original slice was mutated")
	}
}
