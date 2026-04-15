// Package filter provides resource filtering capabilities for drift detection.
// It allows users to include or exclude specific resources by type or ID pattern.
package filter

import (
	"strings"
)

// Rule defines a single filter rule with a resource type and optional ID glob.
type Rule struct {
	ResourceType string
	ResourceID   string // supports "*" wildcard suffix, e.g. "aws_s3_bucket.*"
}

// Filter holds include and exclude rule sets.
type Filter struct {
	includes []Rule
	excludes []Rule
}

// New returns an empty Filter.
func New() *Filter {
	return &Filter{}
}

// AddInclude registers a rule that explicitly allows matching resources.
func (f *Filter) AddInclude(resourceType, resourceID string) {
	f.includes = append(f.includes, Rule{ResourceType: resourceType, ResourceID: resourceID})
}

// AddExclude registers a rule that suppresses matching resources from drift output.
func (f *Filter) AddExclude(resourceType, resourceID string) {
	f.excludes = append(f.excludes, Rule{ResourceType: resourceType, ResourceID: resourceID})
}

// Allow returns true when a resource should be included in drift analysis.
// Exclude rules take precedence over include rules.
func (f *Filter) Allow(resourceType, resourceID string) bool {
	for _, r := range f.excludes {
		if matchRule(r, resourceType, resourceID) {
			return false
		}
	}
	if len(f.includes) == 0 {
		return true
	}
	for _, r := range f.includes {
		if matchRule(r, resourceType, resourceID) {
			return true
		}
	}
	return false
}

// matchRule checks whether a resource matches a single Rule.
// ResourceType must match exactly; ResourceID supports a trailing "*" wildcard.
func matchRule(r Rule, resourceType, resourceID string) bool {
	if r.ResourceType != "*" && r.ResourceType != resourceType {
		return false
	}
	if r.ResourceID == "*" || r.ResourceID == "" {
		return true
	}
	if strings.HasSuffix(r.ResourceID, "*") {
		prefix := strings.TrimSuffix(r.ResourceID, "*")
		return strings.HasPrefix(resourceID, prefix)
	}
	return r.ResourceID == resourceID
}
