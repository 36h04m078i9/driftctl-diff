package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/drift"
	"github.com/owner/driftctl-diff/internal/output"
)

func sampleJUnitChanges() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-abc123",
			Changes: []drift.AttributeChange{
				{Attribute: "instance_type", Expected: "t2.micro", Actual: "t2.small"},
			},
		},
	}
}

func TestJUnitFormatter_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewJUnitFormatter(&buf)
	if err := f.Format([]drift.ResourceDiff{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "testsuites") {
		t.Error("expected testsuites element")
	}
	if strings.Contains(out, "failure") {
		t.Error("expected no failure elements for no drift")
	}
}

func TestJUnitFormatter_WithChanges_ContainsFailure(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewJUnitFormatter(&buf)
	if err := f.Format(sampleJUnitChanges()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "failure") {
		t.Error("expected failure element for drifted resource")
	}
	if !strings.Contains(out, "i-abc123") {
		t.Error("expected resource ID in output")
	}
}

func TestJUnitFormatter_WithChanges_FailureMessage(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewJUnitFormatter(&buf)
	_ = f.Format(sampleJUnitChanges())
	out := buf.String()
	if !strings.Contains(out, "instance_type") {
		t.Error("expected attribute name in failure text")
	}
	if !strings.Contains(out, "t2.micro") {
		t.Error("expected expected value in failure text")
	}
}

func TestJUnitFormatter_NilWriter_DefaultsToStdout(t *testing.T) {
	f := output.NewJUnitFormatter(nil)
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}

func TestJUnitFormatter_ContainsXMLHeader(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewJUnitFormatter(&buf)
	_ = f.Format([]drift.ResourceDiff{})
	if !strings.HasPrefix(buf.String(), "<?xml") {
		t.Error("expected XML declaration at start of output")
	}
}
