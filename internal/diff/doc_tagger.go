// Package diff provides utilities for comparing, annotating, and rendering
// infrastructure drift results.
//
// # Tagger
//
// The Tagger type attaches metadata tags (key-value pairs) to drift results,
// enabling downstream consumers to filter, route, or display results based on
// environment, region, or arbitrary custom labels.
//
// Usage:
//
//	opts := diff.TaggerOptions{
//		EnvTag:    "production",
//		RegionTag: "us-east-1",
//		CustomTags: map[string]string{"team": "platform"},
//	}
//	tagger := diff.NewTagger(opts)
//	tagged := tagger.Tag(results)
//	filtered := diff.FilterByTag(tagged, "env", "production")
package diff
