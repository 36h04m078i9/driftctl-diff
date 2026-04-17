package policy_test

import (
	"os"
	"testing"

	"github.com/acme/driftctl-diff/internal/policy"
)

func writeTempPolicy(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "policy-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoadFile_EmptyPath(t *testing.T) {
	rules, err := policy.LoadFile("")
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 0 {
		t.Errorf("expected no rules")
	}
}

func TestLoadFile_ValidFile(t *testing.T) {
	content := `rules:
  - resource_type: aws_s3_bucket
    severity: high
  - resource_type: aws_instance
    attribute: instance_type
    severity: medium
`
	path := writeTempPolicy(t, content)
	rules, err := policy.LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Severity != policy.SeverityHigh {
		t.Errorf("expected high severity")
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	_, err := policy.LoadFile("/nonexistent/policy.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFile_InvalidSeverity(t *testing.T) {
	content := "rules:\n  - resource_type: aws_s3_bucket\n    severity: critical\n"
	path := writeTempPolicy(t, content)
	_, err := policy.LoadFile(path)
	if err == nil {
		t.Fatal("expected error for invalid severity")
	}
}

func TestLoadFile_MissingResourceType(t *testing.T) {
	content := "rules:\n  - severity: high\n"
	path := writeTempPolicy(t, content)
	_, err := policy.LoadFile(path)
	if err == nil {
		t.Fatal("expected error for missing resource_type")
	}
}
