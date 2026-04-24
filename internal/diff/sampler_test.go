package diff

import (
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeSamplerResults(n int) []drift.DriftResult {
	results := make([]drift.DriftResult, n)
	for i := 0; i < n; i++ {
		results[i] = drift.DriftResult{
			ResourceType: "aws_instance",
			ResourceID:   fmt.Sprintf("res-%02d", i),
			Changes: []drift.AttributeChange{
				{Attribute: "ami", StateValue: "old", LiveValue: "new"},
			},
		}
	}
	return results
}

func TestSampler_NoLimit_ReturnsAll(t *testing.T) {
	opts := DefaultSampleOptions()
	opts.MaxResults = 0
	s := NewSampler(opts)

	results := makeSamplerResults(10)
	got := s.Sample(results)

	if len(got) != 10 {
		t.Fatalf("expected 10 results, got %d", len(got))
	}
}

func TestSampler_MaxResults_LimitsCount(t *testing.T) {
	opts := DefaultSampleOptions()
	opts.MaxResults = 3
	s := NewSampler(opts)

	results := makeSamplerResults(10)
	got := s.Sample(results)

	if len(got) != 3 {
		t.Fatalf("expected 3 results, got %d", len(got))
	}
}

func TestSampler_MaxResultsLargerThanSlice_ReturnsAll(t *testing.T) {
	opts := DefaultSampleOptions()
	opts.MaxResults = 50
	s := NewSampler(opts)

	results := makeSamplerResults(5)
	got := s.Sample(results)

	if len(got) != 5 {
		t.Fatalf("expected 5 results, got %d", len(got))
	}
}

func TestSampler_EmptyInput_ReturnsEmpty(t *testing.T) {
	s := NewSampler(DefaultSampleOptions())
	got := s.Sample([]drift.DriftResult{})
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %d elements", len(got))
	}
}

func TestSampler_Deterministic_SameSeedSameOrder(t *testing.T) {
	opts := DefaultSampleOptions()
	opts.MaxResults = 4
	opts.Seed = 99
	opts.Deterministic = true
	s := NewSampler(opts)

	results := makeSamplerResults(10)
	first := s.Sample(results)
	second := s.Sample(results)

	if len(first) != len(second) {
		t.Fatalf("lengths differ: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if first[i].ResourceID != second[i].ResourceID {
			t.Errorf("index %d: %s != %s", i, first[i].ResourceID, second[i].ResourceID)
		}
	}
}

func TestSampler_DoesNotMutateInput(t *testing.T) {
	opts := DefaultSampleOptions()
	opts.MaxResults = 2
	s := NewSampler(opts)

	results := makeSamplerResults(5)
	originalIDs := make([]string, len(results))
	for i, r := range results {
		originalIDs[i] = r.ResourceID
	}

	s.Sample(results)

	for i, r := range results {
		if r.ResourceID != originalIDs[i] {
			t.Errorf("input mutated at index %d: got %s, want %s", i, r.ResourceID, originalIDs[i])
		}
	}
}
