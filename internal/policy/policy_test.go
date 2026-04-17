package policy_test

import (
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
	"github.com/acme/driftctl-diff/internal/policy"
)

func sampleChanges() []drift.Change {
	return []drift.Change{
		{ResourceType: "aws_s3_bucket", ResourceID: "my-bucket", Attribute: "acl", WantValue: "private", GotValue: "public"},
		{ResourceType: "aws_instance", ResourceID: "i-123", Attribute: "instance_type", WantValue: "t3.micro", GotValue: "t3.large"},
	}
}

func TestEvaluate_NoRules_NoViolations(t *testing.T) {
	e := policy.New(nil)
	v := e.Evaluate(sampleChanges())
	if len(v) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(v))
	}
}

func TestEvaluate_MatchingType(t *testing.T) {
	rules := []policy.Rule{
		{ResourceType: "aws_s3_bucket", Severity: policy.SeverityHigh},
	}
	e := policy.New(rules)
	v := e.Evaluate(sampleChanges())
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Severity != policy.SeverityHigh {
		t.Errorf("expected high severity")
	}
}

func TestEvaluate_WildcardType(t *testing.T) {
	rules := []policy.Rule{
		{ResourceType: "*", Severity: policy.SeverityLow},
	}
	e := policy.New(rules)
	v := e.Evaluate(sampleChanges())
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(v))
	}
}

func TestEvaluate_AttributeFilter(t *testing.T) {
	rules := []policy.Rule{
		{ResourceType: "aws_instance", Attribute: "instance_type", Severity: policy.SeverityMedium},
	}
	e := policy.New(rules)
	v := e.Evaluate(sampleChanges())
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Change.ResourceID != "i-123" {
		t.Errorf("unexpected resource id")
	}
}

func TestMaxSeverity_Empty(t *testing.T) {
	if s := policy.MaxSeverity(nil); s != "" {
		t.Errorf("expected empty, got %s", s)
	}
}

func TestMaxSeverity_ReturnsHighest(t *testing.T) {
	v := []policy.Violation{
		{Severity: policy.SeverityLow},
		{Severity: policy.SeverityHigh},
		{Severity: policy.SeverityMedium},
	}
	if s := policy.MaxSeverity(v); s != policy.SeverityHigh {
		t.Errorf("expected high, got %s", s)
	}
}
