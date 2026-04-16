package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
	"github.com/acme/driftctl-diff/internal/output"
)

func TestDotFormatter_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewDotFormatter(&buf)
	if err := f.Format(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "digraph drift") {
		t.Error("expected DOT graph header")
	}
	if !strings.Contains(out, "No drift detected") {
		t.Error("expected no-drift node")
	}
}

func TestDotFormatter_WithChanges_ContainsResourceNode(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewDotFormatter(&buf)
	changes := []drift.ResourceDiff{
		{
			ResourceID:   "aws_instance.web",
			ResourceType: "aws_instance",
			Changes: []drift.AttributeDiff{
				{Attribute: "instance_type", Expected: "t2.micro", Actual: "t3.small"},
			},
		},
	}
	if err := f.Format(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "aws_instance.web") {
		t.Error("expected resource ID in output")
	}
	if !strings.Contains(out, "instance_type") {
		t.Error("expected attribute name in output")
	}
	if !strings.Contains(out, "->") {
		t.Error("expected edge in DOT graph")
	}
}

func TestDotFormatter_NilWriter_DefaultsToStdout(t *testing.T) {
	f := output.NewDotFormatter(nil)
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}

func TestDotFormatter_SanitizesNodeIDs(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewDotFormatter(&buf)
	changes := []drift.ResourceDiff{
		{
			ResourceID:   "aws-s3-bucket.my.bucket",
			ResourceType: "aws_s3_bucket",
			Changes:      []drift.AttributeDiff{},
		},
	}
	if err := f.Format(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "aws-s3-bucket.my.bucket ") {
		t.Error("node ID should be sanitized (no hyphens or dots)")
	}
}
