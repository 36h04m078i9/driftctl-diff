package diff

import (
	"errors"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makePipelineResults(n int) []drift.DriftResult {
	out := make([]drift.DriftResult, n)
	for i := range out {
		out[i] = drift.DriftResult{ResourceID: fmt.Sprintf("res-%d", i), ResourceType: "aws_instance"}
	}
	return out
}

func TestPipeline_NoSteps_ReturnsInput(t *testing.T) {
	p := NewPipeline()
	input := makePipelineResults(3)
	out, err := p.Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 results, got %d", len(out))
	}
}

func TestPipeline_SingleStep_TransformsResults(t *testing.T) {
	p := NewPipeline(func(r []drift.DriftResult) ([]drift.DriftResult, error) {
		return r[:1], nil
	})
	out, err := p.Run(makePipelineResults(4))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 result, got %d", len(out))
	}
}

func TestPipeline_StepError_AbortsWithIndex(t *testing.T) {
	p := NewPipeline(
		func(r []drift.DriftResult) ([]drift.DriftResult, error) { return r, nil },
		func(r []drift.DriftResult) ([]drift.DriftResult, error) {
			return nil, errors.New("boom")
		},
	)
	_, err := p.Run(makePipelineResults(2))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "step 1") {
		t.Errorf("expected step index in error, got: %v", err)
	}
}

func TestPipeline_MultipleSteps_ChainedCorrectly(t *testing.T) {
	calls := 0
	step := func(r []drift.DriftResult) ([]drift.DriftResult, error) {
		calls++
		return r, nil
	}
	p := NewPipeline(step, step, step)
	_, err := p.Run(makePipelineResults(1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 step calls, got %d", calls)
	}
}

func TestPipelinePrinter_NilWriter_DefaultsToStdout(t *testing.T) {
	pp := NewPipelinePrinter(nil)
	if pp.w == nil {
		t.Error("writer should default to os.Stdout, got nil")
	}
}

func TestPipelinePrinter_Print_ContainsStepCount(t *testing.T) {
	var buf strings.Builder
	pp := NewPipelinePrinter(&buf)
	pp.Print(makePipelineResults(2), 5)
	if !strings.Contains(buf.String(), "5") {
		t.Errorf("expected step count in output, got: %s", buf.String())
	}
}
