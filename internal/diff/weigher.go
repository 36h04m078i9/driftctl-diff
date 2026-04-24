package diff

import (
	"sort"

	"github.com/acme/driftctl-diff/internal/drift"
)

// WeighOptions controls how drift results are weighted.
type WeighOptions struct {
	// SensitiveMultiplier scales the weight of changes to sensitive attributes.
	SensitiveMultiplier float64
	// MissingMultiplier scales the weight of missing-attribute changes.
	MissingMultiplier float64
}

// DefaultWeighOptions returns sensible defaults.
func DefaultWeighOptions() WeighOptions {
	return WeighOptions{
		SensitiveMultiplier: 2.0,
		MissingMultiplier: 1.5,
	}
}

// WeighedResult pairs a DriftResult with a computed weight.
type WeighedResult struct {
	Result drift.DriftResult
	Weight float64
}

// Weigher assigns a numeric weight to each DriftResult so that the most
// impactful drift surfaces first.
type Weigher struct {
	opts             WeighOptions
	sensitiveAttrs   map[string]bool
}

// NewWeigher creates a Weigher with the given options and a default set of
// sensitive attribute names.
func NewWeigher(opts WeighOptions) *Weigher {
	return &Weigher{
		opts: opts,
		sensitiveAttrs: map[string]bool{
			"password": true,
			"secret":   true,
			"token":    true,
			"key":      true,
		},
	}
}

// Weigh scores every result and returns them sorted by descending weight.
func (w *Weigher) Weigh(results []drift.DriftResult) []WeighedResult {
	weighed := make([]WeighedResult, 0, len(results))
	for _, r := range results {
		weighed = append(weighed, WeighedResult{
			Result: r,
			Weight: w.score(r),
		})
	}
	sort.Slice(weighed, func(i, j int) bool {
		return weighed[i].Weight > weighed[j].Weight
	})
	return weighed
}

func (w *Weigher) score(r drift.DriftResult) float64 {
	var total float64
	for _, ch := range r.Changes {
		switch ch.Kind {
		case drift.Missing:
			total += 1.0 * w.opts.MissingMultiplier
		default:
			total += 1.0
		}
		if w.sensitiveAttrs[ch.Attribute] {
			total += 1.0 * (w.opts.SensitiveMultiplier - 1.0)
		}
	}
	return total
}
