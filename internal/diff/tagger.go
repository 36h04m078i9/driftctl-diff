package diff

import "github.com/driftctl-diff/internal/drift"

// Tag is a key-value label attached to a drift result for downstream filtering or display.
type Tag struct {
	Key   string
	Value string
}

// TaggedResult wraps a DriftResult with a set of tags.
type TaggedResult struct {
	drift.DriftResult
	Tags []Tag
}

// TaggerOptions controls which tags are applied.
type TaggerOptions struct {
	EnvTag      string // e.g. "production"
	RegionTag   string // e.g. "us-east-1"
	CustomTags  map[string]string
}

// Tagger attaches metadata tags to drift results.
type Tagger struct {
	opts TaggerOptions
}

// NewTagger returns a Tagger with the given options.
func NewTagger(opts TaggerOptions) *Tagger {
	return &Tagger{opts: opts}
}

// Tag annotates each DriftResult with configured tags and returns TaggedResults.
func (t *Tagger) Tag(results []drift.DriftResult) []TaggedResult {
	out := make([]TaggedResult, 0, len(results))
	for _, r := range results {
		tags := []Tag{}
		if t.opts.EnvTag != "" {
			tags = append(tags, Tag{Key: "env", Value: t.opts.EnvTag})
		}
		if t.opts.RegionTag != "" {
			tags = append(tags, Tag{Key: "region", Value: t.opts.RegionTag})
		}
		for k, v := range t.opts.CustomTags {
			tags = append(tags, Tag{Key: k, Value: v})
		}
		out = append(out, TaggedResult{DriftResult: r, Tags: tags})
	}
	return out
}

// FilterByTag returns only those TaggedResults that carry a tag matching key=value.
func FilterByTag(results []TaggedResult, key, value string) []TaggedResult {
	var out []TaggedResult
	for _, r := range results {
		for _, t := range r.Tags {
			if t.Key == key && t.Value == value {
				out = append(out, r)
				break
			}
		}
	}
	return out
}
