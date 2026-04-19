package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/drift"
)

func TestPatchExporter_NilWriter_DefaultsToStdout(t *testing.T) {
	pe := NewPatchExporter(DefaultPatchOptions(), nil)
	if pe.out == nil {
		t.Error("expected non-nil writer")
	}
}

func TestPatchExporter_NoDrift_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	pe := NewPatchExporter(DefaultPatchOptions(), &buf)
	_ = pe.Export(nil)
	out := buf.String()
	if !strings.Contains(out, "driftctl-diff patch") {
		t.Errorf("expected patch header, got: %s", out)
	}
}

func TestPatchExporter_ResourceCountInHeader(t *testing.T) {
	var buf bytes.Buffer
	pe := NewPatchExporter(DefaultPatchOptions(), &buf)
	results := []drift.ResourceDiff{
		{ResourceType: "aws_s3_bucket", ResourceID: "my-bucket", Changes: []drift.AttributeChange{
			{Attribute: "acl", TerraformValue: "private", LiveValue: "public-read", Kind: drift.KindChanged},
		}},
	}
	_ = pe.Export(results)
	out := buf.String()
	if !strings.Contains(out, "resources: 1") {
		t.Errorf("expected resource count in header, got: %s", out)
	}
}

func TestPatchExporter_WithDrift_ContainsHunk(t *testing.T) {
	var buf bytes.Buffer
	pe := NewPatchExporter(DefaultPatchOptions(), &buf)
	results := []drift.ResourceDiff{
		{ResourceType: "aws_s3_bucket", ResourceID: "my-bucket", Changes: []drift.AttributeChange{
			{Attribute: "acl", TerraformValue: "private", LiveValue: "public-read", Kind: drift.KindChanged},
		}},
	}
	_ = pe.Export(results)
	out := buf.String()
	if !strings.Contains(out, "-private") || !strings.Contains(out, "+public-read") {
		t.Errorf("expected hunk values, got: %s", out)
	}
}
