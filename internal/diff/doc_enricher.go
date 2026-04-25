// Package diff provides utilities for processing and presenting infrastructure
// drift results produced by driftctl-diff.
//
// # Enricher
//
// The Enricher attaches supplementary metadata to a slice of [drift.DriftResult]
// values before they are rendered or exported. Supported enrichments include:
//
//   - Change count: the number of drifted attributes per resource, stored both
//     as a typed field (ChangeCount) and in the Metadata map under the key
//     "change_count".
//
//   - Resource URL: a console hyperlink constructed from a configurable
//     URLTemplate format string. The default template produces a generic
//     placeholder URL; operators should override it with their cloud
//     provider's console base URL.
//
// Use [NewEnricher] with [DefaultEnrichOptions] as a starting point, then
// adjust individual fields to suit your pipeline.
//
// Use [NewEnricherPrinter] to render [EnrichedResult] slices as a human-readable
// tab-separated table to any [io.Writer].
package diff
