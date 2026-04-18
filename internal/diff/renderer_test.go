package diff

import (
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func sampleRenderResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "acl", StateValue: "private", LiveValue: "public-read", Kind: drift.KindChanged},
			},
		},
	}
}

func TestRenderer_NoDrift_PrintsMessage(t *testing.T) {
	var buf strings.Builder
	r := NewRenderer(RenderOptions{Color: false}, &buf)
	if err := r.Render(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected no-drift message, got: %q", buf.String())
	}
}

func TestRenderer_WithDrift_ContainsResourceHeader(t *testing.T) {
	var buf strings.Builder
	r := NewRenderer(RenderOptions{Color: false}, &buf)
	if err := r.Render(sampleRenderResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "aws_s3_bucket/my-bucket") {
		t.Errorf("expected resource header in output, got: %q", out)
	}
}

func TestRenderer_WithDrift_ContainsAttributeValues(t *testing.T) {
	var buf strings.Builder
	r := NewRenderer(RenderOptions{Color: false}, &buf)
	_ = r.Render(sampleRenderResults())
	out := buf.String()
	if !strings.Contains(out, "private") || !strings.Contains(out, "public-read") {
		t.Errorf("expected attribute values in output, got: %q", out)
	}
}

func TestRenderer_NilWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic when w is nil.
	r := NewRenderer(RenderOptions{}, nil)
	if r.w == nil {
		t.Error("expected non-nil writer")
	}
}

func TestRenderer_ColorEnabled_ContainsEscapeCodes(t *testing.T) {
	var buf strings.Builder
	r := NewRenderer(RenderOptions{Color: true}, &buf)
	_ = r.Render(sampleRenderResults())
	out := buf.String()
	if !strings.Contains(out, "\x1b[") {
		t.Errorf("expected ANSI escape codes with color enabled, got: %q", out)
	}
}
