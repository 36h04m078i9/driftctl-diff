package diff

import (
	"fmt"

	"github.com/acme/driftctl-diff/internal/drift"
)

// ChainStep is a named transformation applied to a slice of DriftResults.
type ChainStep struct {
	Name    string
	ApplyFn func([]drift.DriftResult) ([]drift.DriftResult, error)
}

// ChainResult holds the output of a full pipeline run.
type ChainResult struct {
	Results []drift.DriftResult
	// Steps records how many results remained after each named step.
	Steps []ChainStepSummary
}

// ChainStepSummary captures the count of results after a single step.
type ChainStepSummary struct {
	Name  string
	Count int
}

// Chainer executes an ordered sequence of ChainSteps against drift results.
type Chainer struct {
	steps []ChainStep
}

// NewChainer returns a Chainer with the given steps.
func NewChainer(steps []ChainStep) *Chainer {
	return &Chainer{steps: steps}
}

// Run executes each step in order, passing the output of one step as the
// input to the next. It returns a ChainResult that includes per-step counts
// so callers can observe how the pipeline shaped the data.
func (c *Chainer) Run(input []drift.DriftResult) (ChainResult, error) {
	current := input
	summaries := make([]ChainStepSummary, 0, len(c.steps))

	for _, step := range c.steps {
		var err error
		current, err = step.ApplyFn(current)
		if err != nil {
			return ChainResult{}, fmt.Errorf("chainer step %q: %w", step.Name, err)
		}
		summaries = append(summaries, ChainStepSummary{
			Name:  step.Name,
			Count: len(current),
		})
	}

	return ChainResult{
		Results: current,
		Steps:   summaries,
	}, nil
}
