package diff

import (
	"fmt"

	"github.com/acme/driftctl-diff/internal/drift"
)

// PipelineStep is a function that transforms a slice of DriftResults.
type PipelineStep func([]drift.DriftResult) ([]drift.DriftResult, error)

// Pipeline executes an ordered sequence of PipelineSteps, passing the output
// of each step as the input to the next.
type Pipeline struct {
	steps []PipelineStep
}

// NewPipeline returns an empty Pipeline.
func NewPipeline(steps ...PipelineStep) *Pipeline {
	return &Pipeline{steps: steps}
}

// Add appends a step to the pipeline.
func (p *Pipeline) Add(step PipelineStep) {
	p.steps = append(p.steps, step)
}

// Run executes all steps in order and returns the final result.
// If any step returns an error the pipeline is aborted and the error is
// returned with the step index embedded in the message.
func (p *Pipeline) Run(input []drift.DriftResult) ([]drift.DriftResult, error) {
	current := input
	for i, step := range p.steps {
		var err error
		current, err = step(current)
		if err != nil {
			return nil, fmt.Errorf("pipeline step %d: %w", i, err)
		}
	}
	return current, nil
}

// Len returns the number of steps registered in the pipeline.
func (p *Pipeline) Len() int { return len(p.steps) }
