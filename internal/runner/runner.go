// Package runner orchestrates the full drift detection pipeline:
// load config → parse state → fetch live attrs → detect drift → output results.
package runner

import (
	"fmt"
	"io"

	"github.com/acme/driftctl-diff/internal/cache"
	"github.com/acme/driftctl-diff/internal/config"
	"github.com/acme/driftctl-diff/internal/drift"
	"github.com/acme/driftctl-diff/internal/filter"
	"github.com/acme/driftctl-diff/internal/output"
	"github.com/acme/driftctl-diff/internal/provider"
	"github.com/acme/driftctl-diff/internal/state"
	"github.com/acme/driftctl-diff/internal/summary"
)

// Runner executes the end-to-end drift detection workflow.
type Runner struct {
	cfg      *config.Config
	out      io.Writer
}

// New creates a Runner using the supplied config and output destination.
func New(cfg *config.Config, out io.Writer) *Runner {
	return &Runner{cfg: cfg, out: out}
}

// Run performs drift detection and writes human-readable output.
// It returns an error if any stage of the pipeline fails.
func (r *Runner) Run(statePath string) error {
	// 1. Parse Terraform state.
	p := state.NewParser()
	resources, err := p.Parse(statePath)
	if err != nil {
		return fmt.Errorf("parsing state: %w", err)
	}

	// 2. Build provider registry with caching.
	reg := provider.NewRegistry()
	awsProv, err := provider.NewAWSProvider(r.cfg)
	if err != nil {
		return fmt.Errorf("initialising AWS provider: %w", err)
	}
	c := cache.New()
	cached := cache.NewCachedProvider(awsProv, c)
	reg.Register("aws", cached)

	// 3. Apply resource filter.
	f := filter.New(r.cfg.Filters)

	// 4. Detect drift.
	det := drift.NewDetector(reg)
	var results []drift.Result
	for _, res := range resources {
		if !f.Allow(res.Type, res.ID) {
			continue
		}
		result, err := det.Detect(res)
		if err != nil {
			return fmt.Errorf("detecting drift for %s/%s: %w", res.Type, res.ID, err)
		}
		results = append(results, result)
	}

	// 5. Format and write output.
	colorizer := output.NewColorizer(r.cfg.Color)
	fmt := output.NewFormatter(colorizer)
	w := output.NewWriter(fmt, r.out)
	if err := w.Write(results); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	// 6. Print summary.
	sum := summary.Compute(results)
	printer := summary.NewPrinter(r.out)
	printer.Print(sum)

	return nil
}
