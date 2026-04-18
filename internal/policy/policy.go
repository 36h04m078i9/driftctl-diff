// Package policy evaluates drift results against user-defined severity rules.
package policy

import (
	"fmt"

	"github.com/acme/driftctl-diff/internal/drift"
)

// Severity represents how critical a drift finding is.
type Severity string

const (
	SeverityLow    Severity = "low"
	SeverityMedium Severity = "medium"
	SeverityHigh   Severity = "high"
)

// Rule maps a resource type and optional attribute to a severity.
type Rule struct {
	ResourceType string   `yaml:"resource_type"`
	Attribute    string   `yaml:"attribute,omitempty"`
	Severity     Severity `yaml:"severity"`
}

// Violation is a drift change that matched a policy rule.
type Violation struct {
	Change   drift.Change
	Severity Severity
	Rule     Rule
}

// Evaluator checks drift changes against a set of rules.
type Evaluator struct {
	rules []Rule
}

// New creates an Evaluator with the given rules.
func New(rules []Rule) *Evaluator {
	return &Evaluator{rules: rules}
}

// Evaluate returns all violations found in the provided changes.
func (e *Evaluator) Evaluate(changes []drift.Change) []Violation {
	var violations []Violation
	for _, c := range changes {
		for _, r := range e.rules {
			if r.ResourceType != "*" && r.ResourceType != c.ResourceType {
				continue
			}
			if r.Attribute != "" && r.Attribute != c.Attribute {
				continue
			}
			violations = append(violations, Violation{Change: c, Severity: r.Severity, Rule: r})
			break
		}
	}
	return violations
}

// MaxSeverity returns the highest severity among violations, or empty string if none.
func MaxSeverity(violations []Violation) Severity {
	rank := map[Severity]int{SeverityLow: 1, SeverityMedium: 2, SeverityHigh: 3}
	max := 0
	var result Severity
	for _, v := range violations {
		if r := rank[v.Severity]; r > max {
			max = r
			result = v.Severity
		}
	}
	return result
}

// FilterBySeverity returns only the violations that match the given severity.
func FilterBySeverity(violations []Violation, s Severity) []Violation {
	var out []Violation
	for _, v := range violations {
		if v.Severity == s {
			out = append(out, v)
		}
	}
	return out
}

// String formats a Violation for display.
func (v Violation) String() string {
	return fmt.Sprintf("[%s] %s.%s attr=%s", v.Severity, v.Change.ResourceType, v.Change.ResourceID, v.Change.Attribute)
}
