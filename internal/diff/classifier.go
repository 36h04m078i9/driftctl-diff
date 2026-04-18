package diff

import "github.com/owner/driftctl-diff/internal/drift"

// Severity levels for drift classification.
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// Classification holds a drift result with an assigned severity.
type Classification struct {
	Result   drift.ResourceDiff
	Severity Severity
	Reason   string
}

// ClassifierOptions controls how drift is classified.
type ClassifierOptions struct {
	// HighValueAttrs are attribute names that trigger high severity.
	HighValueAttrs []string
	// CriticalTypes are resource types that trigger critical severity.
	CriticalTypes []string
}

// Classifier assigns severity levels to drift results.
type Classifier struct {
	opts ClassifierOptions
}

// NewClassifier creates a Classifier with the given options.
func NewClassifier(opts ClassifierOptions) *Classifier {
	return &Classifier{opts: opts}
}

// Classify assigns a severity and reason to each ResourceDiff.
func (c *Classifier) Classify(results []drift.ResourceDiff) []Classification {
	out := make([]Classification, 0, len(results))
	for _, r := range results {
		sev, reason := c.classify(r)
		out = append(out, Classification{Result: r, Severity: sev, Reason: reason})
	}
	return out
}

func (c *Classifier) classify(r drift.ResourceDiff) (Severity, string) {
	for _, t := range c.opts.CriticalTypes {
		if t == r.ResourceType {
			return SeverityCritical, "critical resource type: " + t
		}
	}
	for _, ch := range r.Changes {
		for _, attr := range c.opts.HighValueAttrs {
			if ch.Attribute == attr {
				return SeverityHigh, "sensitive attribute changed: " + attr
			}
		}
	}
	if len(r.Changes) == 0 {
		return SeverityLow, "no changes detected"
	}
	if len(r.Changes) >= 5 {
		return SeverityMedium, "multiple attributes changed"
	}
	return SeverityLow, "minor drift"
}
