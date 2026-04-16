package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/driftctl-diff/internal/drift"
	"github.com/user/driftctl-diff/internal/output"
)

func sampleXMLChanges() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceID:   "i-0abc123",
			ResourceType: "aws_instance",
			Changes: []drift.AttributeChange{
				{Attribute: "instance_type", Kind: drift.KindChanged, Want: "t3.micro", Got: "t3.small"},
			},
		},
	}
}

func TestXMLFormatter_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewXMLFormatter(&buf)
	if err := f.Format(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Drifted=\"false\"") {
		t.Errorf("expected Drifted=false, got:\n%s", out)
	}
	if strings.Contains(out, "<Resource") {
		t.Errorf("expected no Resource elements for empty drift")
	}
}

func TestXMLFormatter_WithChanges_DriftedTrue(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewXMLFormatter(&buf)
	if err := f.Format(sampleXMLChanges()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Drifted=\"true\"") {
		t.Errorf("expected Drifted=true, got:\n%s", out)
	}
}

func TestXMLFormatter_ContainsResourceInfo(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewXMLFormatter(&buf)
	_ = f.Format(sampleXMLChanges())
	out := buf.String()
	for _, want := range []string{"i-0abc123", "aws_instance", "instance_type", "t3.micro", "t3.small"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

func TestXMLFormatter_NilWriter_DefaultsToStdout(t *testing.T) {
	f := output.NewXMLFormatter(nil)
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}
