package diff

import (
	"fmt"
	"strings"

	"github.com/acme/driftctl-diff/internal/drift"
)

// RewriteRule defines a transformation applied to attribute values during rewriting.
type RewriteRule struct {
	// ResourceType to match ("*" for all).
	ResourceType string
	// Attribute name to match ("*" for all).
	Attribute string
	// Find is the substring to replace.
	Find string
	// Replace is the replacement string.
	Replace string
}

// RewriteOptions controls Rewriter behaviour.
type RewriteOptions struct {
	Rules []RewriteRule
}

// DefaultRewriteOptions returns an empty RewriteOptions.
func DefaultRewriteOptions() RewriteOptions {
	return RewriteOptions{}
}

// Rewriter applies string-substitution rules to drift result attribute values.
type Rewriter struct {
	opts RewriteOptions
}

// NewRewriter constructs a Rewriter with the provided options.
func NewRewriter(opts RewriteOptions) *Rewriter {
	return &Rewriter{opts: opts}
}

// Rewrite returns a new slice of DriftResults with attribute values rewritten
// according to the configured rules. Original results are not mutated.
func (r *Rewriter) Rewrite(results []drift.DriftResult) []drift.DriftResult {
	out := make([]drift.DriftResult, len(results))
	for i, res := range results {
		out[i] = drift.DriftResult{
			ResourceID:   res.ResourceID,
			ResourceType: res.ResourceType,
			Changes:      r.rewriteChanges(res.ResourceType, res.Changes),
		}
	}
	return out
}

func (r *Rewriter) rewriteChanges(resourceType string, changes []drift.AttributeChange) []drift.AttributeChange {
	out := make([]drift.AttributeChange, len(changes))
	for i, ch := range changes {
		out[i] = drift.AttributeChange{
			Attribute: ch.Attribute,
			Kind:      ch.Kind,
			WantValue: r.applyRules(resourceType, ch.Attribute, ch.WantValue),
			GotValue:  r.applyRules(resourceType, ch.Attribute, ch.GotValue),
		}
	}
	return out
}

func (r *Rewriter) applyRules(resourceType, attribute, value string) string {
	for _, rule := range r.opts.Rules {
		if !matchesPattern(rule.ResourceType, resourceType) {
			continue
		}
		if !matchesPattern(rule.Attribute, attribute) {
			continue
		}
		value = strings.ReplaceAll(value, rule.Find, rule.Replace)
	}
	return value
}

func matchesPattern(pattern, value string) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(value, strings.TrimSuffix(pattern, "*"))
	}
	return fmt.Sprintf("%s", pattern) == value
}
