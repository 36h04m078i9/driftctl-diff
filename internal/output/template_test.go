package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/drift"
)

func TestTemplateFormatter_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f, err := NewTemplateFormatter("", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := f.Format(nil); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %q", buf.String())
	}
}

func TestTemplateFormatter_WithChanges_ContainsResourceInfo(t *testing.T) {
	var buf bytes.Buffer
	f, err := NewTemplateFormatter("", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	changes := []drift.ResourceDiff{
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-abc123",
			Attributes: []drift.AttributeDiff{
				{Attribute: "instance_type", StateValue: "t2.micro", LiveValue: "t3.small"},
			},
		},
	}
	if err := f.Format(changes); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"aws_instance", "i-abc123", "instance_type", "t2.micro", "t3.small"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %q", want, out)
		}
	}
}

func TestTemplateFormatter_CustomTemplate(t *testing.T) {
	var buf bytes.Buffer
	custom := `CUSTOM:{{range .Changes}}{{.ResourceID}}{{end}}`
	f, err := NewTemplateFormatter(custom, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	changes := []drift.ResourceDiff{{ResourceType: "aws_s3_bucket", ResourceID: "my-bucket"}}
	if err := f.Format(changes); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	if !strings.Contains(buf.String(), "CUSTOM:my-bucket") {
		t.Errorf("custom template not applied, got: %q", buf.String())
	}
}

func TestTemplateFormatter_InvalidTemplate_ReturnsError(t *testing.T) {
	_, err := NewTemplateFormatter("{{.Unclosed", nil)
	if err == nil {
		t.Error("expected parse error for invalid template")
	}
}

func TestTemplateFormatter_NilWriter_DefaultsToStdout(t *testing.T) {
	f, err := NewTemplateFormatter("", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.out == nil {
		t.Error("expected non-nil writer")
	}
}
