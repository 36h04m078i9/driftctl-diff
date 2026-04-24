package diff

import (
	"math/rand"
	"sort"

	"github.com/acme/driftctl-diff/internal/drift"
)

// SampleOptions controls how results are sampled.
type SampleOptions struct {
	// MaxResults is the maximum number of DriftResults to return.
	// Zero means no limit.
	MaxResults int
	// Seed is used for reproducible random sampling. Zero uses a random seed.
	Seed int64
	// Deterministic sorts results by ResourceID before sampling so output is
	// stable across runs when a fixed Seed is provided.
	Deterministic bool
}

// DefaultSampleOptions returns a SampleOptions with sensible defaults.
func DefaultSampleOptions() SampleOptions {
	return SampleOptions{
		MaxResults:    0,
		Seed:          42,
		Deterministic: true,
	}
}

// Sampler draws a random subset of DriftResults.
type Sampler struct {
	opts SampleOptions
}

// NewSampler creates a Sampler with the provided options.
func NewSampler(opts SampleOptions) *Sampler {
	return &Sampler{opts: opts}
}

// Sample returns up to opts.MaxResults entries from results.
// When MaxResults is zero or greater than len(results) all results are
// returned. The selection is pseudo-random, seeded by opts.Seed.
func (s *Sampler) Sample(results []drift.DriftResult) []drift.DriftResult {
	if len(results) == 0 {
		return results
	}

	// Work on a copy so we never mutate the caller's slice.
	copied := make([]drift.DriftResult, len(results))
	copy(copied, results)

	if s.opts.Deterministic {
		sort.Slice(copied, func(i, j int) bool {
			if copied[i].ResourceType != copied[j].ResourceType {
				return copied[i].ResourceType < copied[j].ResourceType
			}
			return copied[i].ResourceID < copied[j].ResourceID
		})
	}

	max := s.opts.MaxResults
	if max <= 0 || max >= len(copied) {
		return copied
	}

	//nolint:gosec // predictable seed is intentional for reproducibility
	r := rand.New(rand.NewSource(s.opts.Seed))
	r.Shuffle(len(copied), func(i, j int) {
		copied[i], copied[j] = copied[j], copied[i]
	})

	return copied[:max]
}
