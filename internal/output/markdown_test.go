package output_test

import (
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
	"github.com/acme/driftctl-diff/internal/output"
)

func TestMarkdownFormatter_NoDrift(t *testing.T) {
	f := output.NewMarkdownFormatter()
	var buf strings.Builder
	if err := f.Format(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", got)
	}
}

func TestMarkdownFormatter_WithChanges_ContainsHeaders(t *testing.T) {
	f := output.NewMarkdownFormatter()
	changes := []drift.ResourceDiff{
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-abc123",
			Changes: []drift.AttributeChange{
				{Attribute: "instance_type", StateValue: "t2.micro", LiveValue: "t3.small"},
			},
		},
	}
	var buf strings.Builder
	if err := f.Format(&buf, changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	for _, want := range []string{"## Drift Report", "| Attribute |", "aws_instance", "i-abc123"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q in output, got:\n%s", want, got)
		}
	}
}

func TestMarkdownFormatter_WithChanges_ContainsAttributeRow(t *testing.T) {
	f := output.NewMarkdownFormatter()
	changes := []drift.ResourceDiff{
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "versioning", StateValue: "true", LiveValue: "false"},
			},
		},
	}
	var buf strings.Builder
	_ = f.Format(&buf, changes)
	got := buf.String()
	if !strings.Contains(got, "versioning") || !strings.Contains(got, "false") {
		t.Errorf("expected attribute row in output, got:\n%s", got)
	}
}

func TestMarkdownFormatter_DriftedCount(t *testing.T) {
	f := output.NewMarkdownFormatter()
	changes := []drift.ResourceDiff{
		{ResourceType: "aws_instance", ResourceID: "id-1", Changes: []drift.AttributeChange{{Attribute: "a", StateValue: "x", LiveValue: "y"}}},
		{ResourceType: "aws_instance", ResourceID: "id-2", Changes: []drift.AttributeChange{{Attribute: "b", StateValue: "1", LiveValue: "2"}}},
	}
	var buf strings.Builder
	_ = f.Format(&buf, changes)
	if !strings.Contains(buf.String(), "2 resource(s) drifted") {
		t.Errorf("expected count in output, got: %s", buf.String())
	}
}
