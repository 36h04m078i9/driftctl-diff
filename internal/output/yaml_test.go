package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/snyk/driftctl-diff/internal/drift"
)

func TestYAMLFormatter_NoDrift(t *testing.T) {
	f := NewYAMLFormatter()
	var buf bytes.Buffer
	if err := f.Format(nil, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "drifted: false") {
		t.Errorf("expected drifted: false, got:\n%s", out)
	}
	if !strings.Contains(out, "total_drift: 0") {
		t.Errorf("expected total_drift: 0, got:\n%s", out)
	}
}

func TestYAMLFormatter_WithChanges_DriftedTrue(t *testing.T) {
	changes := []drift.ResourceDiff{
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-abc123",
			Attributes: []drift.AttributeDiff{
				{Name: "instance_type", Kind: drift.KindChanged, Want: "t2.micro", Got: "t3.small"},
			},
		},
	}
	f := NewYAMLFormatter()
	var buf bytes.Buffer
	if err := f.Format(changes, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "drifted: true") {
		t.Errorf("expected drifted: true, got:\n%s", out)
	}
	if !strings.Contains(out, "total_drift: 1") {
		t.Errorf("expected total_drift: 1, got:\n%s", out)
	}
}

func TestYAMLFormatter_ContainsResourceInfo(t *testing.T) {
	changes := []drift.ResourceDiff{
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Attributes: []drift.AttributeDiff{
				{Name: "versioning", Kind: drift.KindMissing, Want: "enabled", Got: ""},
			},
		},
	}
	f := NewYAMLFormatter()
	var buf bytes.Buffer
	_ = f.Format(changes, &buf)
	out := buf.String()
	for _, want := range []string{"aws_s3_bucket", "my-bucket", "versioning"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestYAMLFormatter_NilWriter_DefaultsToStdout(t *testing.T) {
	f := NewYAMLFormatter()
	// Should not panic with nil writer.
	if err := f.Format(nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
