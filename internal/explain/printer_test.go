package explain_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/explain"
)

func TestPrinter_NoExplanations_PrintsMessage(t *testing.T) {
	var buf bytes.Buffer
	p := explain.NewPrinter(&buf)
	p.Print(nil)
	if !strings.Contains(buf.String(), "No drift explanations") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestPrinter_WithExplanations_ContainsSeverity(t *testing.T) {
	var buf bytes.Buffer
	p := explain.NewPrinter(&buf)
	exps := []explain.Explanation{
		{ResourceID: "aws_instance.web", Attribute: "ami", Severity: explain.SeverityWarning, Message: "changed"},
	}
	p.Print(exps)
	out := buf.String()
	if !strings.Contains(out, "WARNING") {
		t.Errorf("expected WARNING in output, got: %s", out)
	}
}

func TestPrinter_WithExplanations_ContainsResourceID(t *testing.T) {
	var buf bytes.Buffer
	p := explain.NewPrinter(&buf)
	exps := []explain.Explanation{
		{ResourceID: "aws_s3_bucket.data", Attribute: "policy", Severity: explain.SeverityWarning, Message: "policy changed"},
	}
	p.Print(exps)
	if !strings.Contains(buf.String(), "aws_s3_bucket.data") {
		t.Errorf("expected resource ID in output")
	}
}

func TestPrinter_NilWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic.
	p := explain.NewPrinter(nil)
	p.Print(nil)
}
