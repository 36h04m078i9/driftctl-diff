package explain_test

import (
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
	"github.com/acme/driftctl-diff/internal/explain"
)

func sampleResults(attr, state, live string) []drift.Result {
	return []drift.Result{
		{
			ResourceID:   "aws_instance.web",
			ResourceType: "aws_instance",
			Changes: []drift.Change{
				{Attribute: attr, StateValue: state, LiveValue: live},
			},
		},
	}
}

func TestExplain_DefaultSeverity_Info(t *testing.T) {
	e := explain.New()
	results := sampleResults("tags", "v1", "v2")
	exps := e.Explain(results)
	if len(exps) != 1 {
		t.Fatalf("expected 1 explanation, got %d", len(exps))
	}
	if exps[0].Severity != explain.SeverityInfo {
		t.Errorf("expected info severity, got %s", exps[0].Severity)
	}
}

func TestExplain_PasswordAttribute_CriticalSeverity(t *testing.T) {
	e := explain.New()
	results := sampleResults("db_password", "old", "new")
	exps := e.Explain(results)
	if exps[0].Severity != explain.SeverityCritical {
		t.Errorf("expected critical severity, got %s", exps[0].Severity)
	}
}

func TestExplain_PolicyAttribute_WarningSeverity(t *testing.T) {
	e := explain.New()
	results := sampleResults("iam_policy", "a", "b")
	exps := e.Explain(results)
	if exps[0].Severity != explain.SeverityWarning {
		t.Errorf("expected warning severity, got %s", exps[0].Severity)
	}
}

func TestExplain_PortAttribute_WarningSeverity(t *testing.T) {
	e := explain.New()
	results := sampleResults("ingress_port", "80", "443")
	exps := e.Explain(results)
	if exps[0].Severity != explain.SeverityWarning {
		t.Errorf("expected warning severity, got %s", exps[0].Severity)
	}
}

func TestExplain_NoChanges_EmptyExplanations(t *testing.T) {
	e := explain.New()
	results := []drift.Result{{ResourceID: "aws_instance.web", Changes: nil}}
	exps := e.Explain(results)
	if len(exps) != 0 {
		t.Errorf("expected 0 explanations, got %d", len(exps))
	}
}

func TestExplain_MessageContainsAttribute(t *testing.T) {
	e := explain.New()
	results := sampleResults("ami", "ami-old", "ami-new")
	exps := e.Explain(results)
	if len(exps) == 0 {
		t.Fatal("expected explanation")
	}
	if exps[0].Attribute != "ami" {
		t.Errorf("expected attribute ami, got %s", exps[0].Attribute)
	}
}
