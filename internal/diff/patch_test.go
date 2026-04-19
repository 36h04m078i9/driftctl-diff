package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/drift"
)

func samplePatchResults() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-abc123",
			Changes: []drift.AttributeChange{
				{Attribute: "instance_type", TerraformValue: "t2.micro", LiveValue: "t3.small", Kind: drift.KindChanged},
			},
		},
	}
}

func TestPatcher_NoDrift_PrintsNoDriftMessage(t *testing.T) {
	var buf bytes.Buffer
	p := NewPatcher(DefaultPatchOptions(), &buf)
	_ = p.Generate(nil)
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestPatcher_WithChanges_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	p := NewPatcher(DefaultPatchOptions(), &buf)
	_ = p.Generate(samplePatchResults())
	out := buf.String()
	if !strings.Contains(out, "--- terraform") {
		t.Errorf("expected terraform header line, got: %s", out)
	}
	if !strings.Contains(out, "+++ live") {
		t.Errorf("expected live header line, got: %s", out)
	}
}

func TestPatcher_WithChanges_ContainsHunk(t *testing.T) {
	var buf bytes.Buffer
	p := NewPatcher(DefaultPatchOptions(), &buf)
	_ = p.Generate(samplePatchResults())
	out := buf.String()
	if !strings.Contains(out, "-t2.micro") {
		t.Errorf("expected removed line, got: %s", out)
	}
	if !strings.Contains(out, "+t3.small") {
		t.Errorf("expected added line, got: %s", out)
	}
}

func TestPatcher_NoHeader_OmitsHeaderLines(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultPatchOptions()
	opts.Header = false
	p := NewPatcher(opts, &buf)
	_ = p.Generate(samplePatchResults())
	out := buf.String()
	if strings.Contains(out, "---") {
		t.Errorf("expected no header, got: %s", out)
	}
}

func TestPatcher_NilWriter_DefaultsToStdout(t *testing.T) {
	p := NewPatcher(DefaultPatchOptions(), nil)
	if p.out == nil {
		t.Error("expected non-nil writer")
	}
}
