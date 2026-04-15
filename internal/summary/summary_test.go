package summary_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
	"github.com/acme/driftctl-diff/internal/summary"
)

func makeResults(drifted, clean int) []drift.ResourceDiff {
	var results []drift.ResourceDiff
	for i := 0; i < drifted; i++ {
		results = append(results, drift.ResourceDiff{
			ResourceID:   fmt.Sprintf("res-%d", i),
			ResourceType: "aws_instance",
			Changes: []drift.Change{
				{Attribute: "ami", StateValue: "old", LiveValue: "new"},
			},
		})
	}
	for i := 0; i < clean; i++ {
		results = append(results, drift.ResourceDiff{
			ResourceID:   fmt.Sprintf("clean-%d", i),
			ResourceType: "aws_instance",
			Changes:      nil,
		})
	}
	return results
}

func TestCompute_NoDrift(t *testing.T) {
	stats := summary.Compute(makeResults(0, 3))
	if stats.TotalResources != 3 {
		t.Errorf("expected 3 total, got %d", stats.TotalResources)
	}
	if stats.DriftedResources != 0 {
		t.Errorf("expected 0 drifted, got %d", stats.DriftedResources)
	}
	if stats.CleanResources != 3 {
		t.Errorf("expected 3 clean, got %d", stats.CleanResources)
	}
	if stats.TotalChanges != 0 {
		t.Errorf("expected 0 changes, got %d", stats.TotalChanges)
	}
}

func TestCompute_WithDrift(t *testing.T) {
	stats := summary.Compute(makeResults(2, 1))
	if stats.TotalResources != 3 {
		t.Errorf("expected 3 total, got %d", stats.TotalResources)
	}
	if stats.DriftedResources != 2 {
		t.Errorf("expected 2 drifted, got %d", stats.DriftedResources)
	}
	if stats.TotalChanges != 2 {
		t.Errorf("expected 2 changes, got %d", stats.TotalChanges)
	}
}

func TestPrinter_Print_ContainsKeyFields(t *testing.T) {
	var buf bytes.Buffer
	p := summary.NewPrinter(&buf)
	p.Print(summary.Stats{
		TotalResources:   5,
		DriftedResources: 2,
		CleanResources:   3,
		TotalChanges:     4,
	})
	out := buf.String()
	for _, want := range []string{"5", "2", "3", "4", "Summary"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}
