package diff

import (
	"fmt"
	"strings"

	"github.com/acme/driftctl-diff/internal/drift"
)

// EnrichOptions controls how results are enriched.
type EnrichOptions struct {
	// AddResourceURL generates a console URL for each resource when true.
	AddResourceURL bool
	// URLTemplate is a format string accepting (resourceType, resourceID).
	// Defaults to a generic placeholder when empty.
	URLTemplate string
	// AddChangeCount attaches the number of drifted attributes as metadata.
	AddChangeCount bool
}

// DefaultEnrichOptions returns sensible defaults.
func DefaultEnrichOptions() EnrichOptions {
	return EnrichOptions{
		AddResourceURL: false,
		URLTemplate:    "https://console.example.com/%s/%s",
		AddChangeCount: true,
	}
}

// EnrichedResult wraps a DriftResult with additional metadata.
type EnrichedResult struct {
	drift.DriftResult
	ResourceURL string
	ChangeCount int
	Metadata    map[string]string
}

// Enricher attaches extra metadata to drift results.
type Enricher struct {
	opts EnrichOptions
}

// NewEnricher creates an Enricher with the given options.
func NewEnricher(opts EnrichOptions) *Enricher {
	return &Enricher{opts: opts}
}

// Enrich iterates over results and produces EnrichedResult values.
func (e *Enricher) Enrich(results []drift.DriftResult) []EnrichedResult {
	out := make([]EnrichedResult, 0, len(results))
	for _, r := range results {
		er := EnrichedResult{
			DriftResult: r,
			Metadata:    make(map[string]string),
		}
		if e.opts.AddChangeCount {
			er.ChangeCount = len(r.Changes)
			er.Metadata["change_count"] = fmt.Sprintf("%d", er.ChangeCount)
		}
		if e.opts.AddResourceURL {
			tmpl := e.opts.URLTemplate
			if tmpl == "" {
				tmpl = "https://console.example.com/%s/%s"
			}
			er.ResourceURL = fmt.Sprintf(tmpl,
				strings.ReplaceAll(r.ResourceType, "_", "-"),
				r.ResourceID,
			)
			er.Metadata["resource_url"] = er.ResourceURL
		}
		out = append(out, er)
	}
	return out
}
