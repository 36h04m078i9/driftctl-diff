package diff

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/drift"
)

func sampleExportDiffs() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-abc123",
			Changes: []drift.AttributeChange{
				{Attribute: "instance_type", StateValue: "t2.micro", LiveValue: "t3.small"},
			},
		},
	}
}

func TestExporter_TextFormat_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	e := NewExporter("text", &buf)
	if err := e.Export(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestExporter_TextFormat_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	e := NewExporter("text", &buf)
	if err := e.Export(sampleExportDiffs()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "aws_instance") {
		t.Errorf("expected resource type in output, got: %s", out)
	}
	if !strings.Contains(out, "instance_type") {
		t.Errorf("expected attribute in output, got: %s", out)
	}
}

func TestExporter_JSONFormat_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	e := NewExporter("json", &buf)
	if err := e.Export(sampleExportDiffs()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []drift.ResourceDiff
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 result, got %d", len(out))
	}
}

func TestExporter_UnsupportedFormat_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	e := NewExporter("xml", &buf)
	if err := e.Export(sampleExportDiffs()); err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExporter_NilWriter_DefaultsToStdout(t *testing.T) {
	e := NewExporter("text", nil)
	if e.dest == nil {
		t.Error("expected non-nil dest when nil passed")
	}
}
