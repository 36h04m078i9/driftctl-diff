// Package explain provides human-readable explanations for detected drift.
package explain

import (
	"fmt"
	"strings"

	"github.com/acme/driftctl-diff/internal/drift"
)

// Severity levels for drift explanations.
type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityCritical Severity = "critical"
)

// Explanation describes why a drift change matters.
type Explanation struct {
	ResourceID string
	Attribute  string
	Severity   Severity
	Message    string
}

// Explainer generates explanations for drift results.
type Explainer struct {
	rules []Rule
}

// Rule maps an attribute pattern to a severity and message template.
type Rule struct {
	AttributeContains string
	Severity          Severity
	MessageTemplate   string
}

// New returns an Explainer with default rules.
func New() *Explainer {
	return &Explainer{
		rules: defaultRules(),
	}
}

// Explain returns explanations for all changes in the given results.
func (e *Explainer) Explain(results []drift.Result) []Explanation {
	var out []Explanation
	for _, r := range results {
		for _, c := range r.Changes {
			out = append(out, e.explainChange(r.ResourceID, c))
		}
	}
	return out
}

func (e *Explainer) explainChange(resourceID string, c drift.Change) Explanation {
	for _, rule := range e.rules {
		if strings.Contains(c.Attribute, rule.AttributeContains) {
			return Explanation{
				ResourceID: resourceID,
				Attribute:  c.Attribute,
				Severity:   rule.Severity,
				Message:    fmt.Sprintf(rule.MessageTemplate, c.Attribute, c.StateValue, c.LiveValue),
			}
		}
	}
	return Explanation{
		ResourceID: resourceID,
		Attribute:  c.Attribute,
		Severity:   SeverityInfo,
		Message:    fmt.Sprintf("attribute %q changed from %q to %q", c.Attribute, c.StateValue, c.LiveValue),
	}
}

func defaultRules() []Rule {
	return []Rule{
		{AttributeContains: "password", Severity: SeverityCritical, MessageTemplate: "sensitive attribute %q changed (state: %q, live: %q) — possible secret rotation or breach"},
		{AttributeContains: "policy", Severity: SeverityWarning, MessageTemplate: "IAM/resource policy %q changed (state: %q, live: %q) — review for privilege escalation"},
		{AttributeContains: "port", Severity: SeverityWarning, MessageTemplate: "network attribute %q changed (state: %q, live: %q) — verify firewall rules"},
	}
}
