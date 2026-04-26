package diff

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeChainerResults(n int) []drift.DriftResult {
	out := make([]drift.DriftResult, n)
	for i := range out {
		out[i] = drift.DriftResult{
			ResourceID:   fmt.Sprintf("res-%d", i),
			ResourceType: "aws_s3_bucket",
		}
	}
	return out
}

func TestChainer_EmptySteps_ReturnsInput(t *testing.T) {
	input := makeChainerResults(3)
	c := NewChainer(nil)
	cr, err := c.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cr.Results) != 3 {
		t.Errorf("expected 3 results, got %d", len(cr.Results))
	}
	if len(cr.Steps) != 0 {
		t.Errorf("expected 0 step summaries, got %d", len(cr.Steps))
	}
}

func TestChainer_SingleStep_FiltersResults(t *testing.T) {
	input := makeChainerResults(5)
	step := ChainStep{
		Name: "keep-first-two",
		ApplyFn: func(rs []drift.DriftResult) ([]drift.DriftResult, error) {
			if len(rs) > 2 {
				return rs[:2], nil
			}
			return rs, nil
		},
	}
	c := NewChainer([]ChainStep{step})
	cr, err := c.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cr.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(cr.Results))
	}
	if cr.Steps[0].Count != 2 {
		t.Errorf("step summary count mismatch: got %d", cr.Steps[0].Count)
	}
}

func TestChainer_MultipleSteps_CountsArePerStep(t *testing.T) {
	input := makeChainerResults(6)
	steps := []ChainStep{
		{Name: "half", ApplyFn: func(rs []drift.DriftResult) ([]drift.DriftResult, error) {
			return rs[:len(rs)/2], nil
		}},
		{Name: "all", ApplyFn: func(rs []drift.DriftResult) ([]drift.DriftResult, error) {
			return rs, nil
		}},
	}
	c := NewChainer(steps)
	cr, err := c.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cr.Steps[0].Count != 3 {
		t.Errorf("step 0 count: want 3, got %d", cr.Steps[0].Count)
	}
	if cr.Steps[1].Count != 3 {
		t.Errorf("step 1 count: want 3, got %d", cr.Steps[1].Count)
	}
}

func TestChainer_StepError_ReturnsWrappedError(t *testing.T) {
	input := makeChainerResults(2)
	boom := errors.New("boom")
	step := ChainStep{
		Name:    "fail",
		ApplyFn: func(rs []drift.DriftResult) ([]drift.DriftResult, error) { return nil, boom },
	}
	c := NewChainer([]ChainStep{step})
	_, err := c.Run(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "fail") {
		t.Errorf("error should mention step name: %v", err)
	}
	if !errors.Is(err, boom) {
		t.Errorf("error should wrap original: %v", err)
	}
}

func TestChainerPrinter_NoSteps_PrintsMessage(t *testing.T) {
	var buf bytes.Buffer
	p := NewChainerPrinter(&buf)
	p.Print(ChainResult{})
	if !strings.Contains(buf.String(), "no pipeline steps") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestChainerPrinter_WithSteps_ContainsStepName(t *testing.T) {
	var buf bytes.Buffer
	p := NewChainerPrinter(&buf)
	p.Print(ChainResult{
		Results: makeChainerResults(2),
		Steps: []ChainStepSummary{
			{Name: "filter-step", Count: 2},
		},
	})
	if !strings.Contains(buf.String(), "filter-step") {
		t.Errorf("expected step name in output: %q", buf.String())
	}
	if !strings.Contains(buf.String(), "final result count: 2") {
		t.Errorf("expected final count in output: %q", buf.String())
	}
}
