package diff_test

import (
	"testing"

	"github.com/driftctl/driftctl-diff/internal/diff"
	"github.com/driftctl/driftctl-diff/internal/drift"
)

func makeStatsResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceID:   "res-1",
			ResourceType: "aws_instance",
			Changes: []drift.AttributeChange{
				{Attribute: "ami", Kind: drift.Changed, StateValue: "old", LiveValue: "new"},
				{Attribute: "tags", Kind: drift.Deleted, StateValue: "v", LiveValue: ""},
			},
		},
		{
			ResourceID:   "res-2",
			ResourceType: "aws_s3_bucket",
			Changes:      []drift.AttributeChange{},
		},
		{
			ResourceID:   "res-3",
			ResourceType: "aws_instance",
			Changes: []drift.AttributeChange{
				{Attribute: "subnet_id", Kind: drift.Added, StateValue: "", LiveValue: "subnet-abc"},
			},
		},
	}
}

func TestCompute_TotalResources(t *testing.T) {
	s := diff.Compute(makeStatsResults())
	if s.TotalResources != 3 {
		t.Fatalf("expected 3 total resources, got %d", s.TotalResources)
	}
}

func TestCompute_DriftedResources(t *testing.T) {
	s := diff.Compute(makeStatsResults())
	if s.DriftedResources != 2 {
		t.Fatalf("expected 2 drifted resources, got %d", s.DriftedResources)
	}
}

func TestCompute_AttributeCounts(t *testing.T) {
	s := diff.Compute(makeStatsResults())
	if s.ChangedAttributes != 1 {
		t.Errorf("expected 1 changed, got %d", s.ChangedAttributes)
	}
	if s.DeletedAttributes != 1 {
		t.Errorf("expected 1 deleted, got %d", s.DeletedAttributes)
	}
	if s.AddedAttributes != 1 {
		t.Errorf("expected 1 added, got %d", s.AddedAttributes)
	}
}

func TestCompute_Empty(t *testing.T) {
	s := diff.Compute(nil)
	if s.TotalResources != 0 || s.DriftedResources != 0 {
		t.Error("expected zero stats for empty input")
	}
}

func TestDriftPercent_NoResources(t *testing.T) {
	s := diff.Stats{}
	if s.DriftPercent() != 0 {
		t.Error("expected 0 percent for empty stats")
	}
}

func TestDriftPercent_Calculated(t *testing.T) {
	s := diff.Compute(makeStatsResults())
	got := s.DriftPercent()
	want := 200.0 / 3.0
	if got < want-0.01 || got > want+0.01 {
		t.Errorf("expected ~%.2f%%, got %.2f%%", want, got)
	}
}
