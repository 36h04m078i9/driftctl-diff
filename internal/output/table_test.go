package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
	"github.com/acme/driftctl-diff/internal/output"
)

func TestTableFormatter_NoDrift_PrintsNoDriftMessage(t *testing.T) {
	c := output.NewColorizer(false)
	f := output.NewTableFormatter(c)
	var buf bytes.Buffer

	if err := f.Write(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %q", buf.String())
	}
}

func TestTableFormatter_WithChanges_ContainsHeaders(t *testing.T) {
	c := output.NewColorizer(false)
	f := output.NewTableFormatter(c)
	var buf bytes.Buffer

	results := []drift.ResourceDiff{
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-abc123",
			Changes: []drift.Change{
				{Attribute: "instance_type", StateValue: "t2.micro", LiveValue: "t3.small"},
			},
		},
	}

	if err := f.Write(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"RESOURCE TYPE", "RESOURCE ID", "ATTRIBUTE", "CHANGE"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected header %q in output, got: %q", want, out)
		}
	}
}

func TestTableFormatter_WithChanges_ContainsResourceInfo(t *testing.T) {
	c := output.NewColorizer(false)
	f := output.NewTableFormatter(c)
	var buf bytes.Buffer

	results := []drift.ResourceDiff{
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.Change{
				{Attribute: "acl", StateValue: "private", LiveValue: "public-read"},
				{Attribute: "region", StateValue: "", LiveValue: "us-east-1"},
			},
		},
	}

	if err := f.Write(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"aws_s3_bucket", "my-bucket", "acl", "private -> public-read", "(missing) -> us-east-1"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %q", want, out)
		}
	}
}
