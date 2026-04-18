package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func sampleDifferResult(changes ...drift.AttributeChange) drift.DriftResult {
	return drift.DriftResult{
		ResourceType: "aws_s3_bucket",
		ResourceID:   "my-bucket",
		Changes:      changes,
	}
}

func TestDiffer_NoDrift_PrintsNoDriftMessage(t *testing.T) {
	var buf bytes.Buffer
	d := NewDiffer(DefaultDiffOptions(), &buf)
	_ = d.Render(sampleDifferResult())
	if !strings.Contains(buf.String(), "no drift") {
		t.Fatalf("expected 'no drift', got: %s", buf.String())
	}
}

func TestDiffer_WithChanges_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	d := NewDiffer(DefaultDiffOptions(), &buf)
	ch := drift.AttributeChange{Attribute: "acl", StateValue: "private", LiveValue: "public-read", Kind: drift.KindChanged}
	_ = d.Render(sampleDifferResult(ch))
	out := buf.String()
	if !strings.Contains(out, "--- aws_s3_bucket/my-bucket") {
		t.Fatalf("missing state header, got: %s", out)
	}
	if !strings.Contains(out, "+++ aws_s3_bucket/my-bucket") {
		t.Fatalf("missing live header, got: %s", out)
	}
}

func TestDiffer_WithChanges_ContainsAttributeHunk(t *testing.T) {
	var buf bytes.Buffer
	d := NewDiffer(DefaultDiffOptions(), &buf)
	ch := drift.AttributeChange{Attribute: "acl", StateValue: "private", LiveValue: "public-read", Kind: drift.KindChanged}
	_ = d.Render(sampleDifferResult(ch))
	out := buf.String()
	if !strings.Contains(out, "@@ attribute: acl @@") {
		t.Fatalf("missing hunk header, got: %s", out)
	}
	if !strings.Contains(out, "-  private") {
		t.Fatalf("missing removed line, got: %s", out)
	}
	if !strings.Contains(out, "+  public-read") {
		t.Fatalf("missing added line, got: %s", out)
	}
}

func TestDiffer_ColorEnabled_ContainsEscapeCodes(t *testing.T) {
	var buf bytes.Buffer
	opts := DiffOptions{ContextLines: 3, Color: true}
	d := NewDiffer(opts, &buf)
	ch := drift.AttributeChange{Attribute: "tags", StateValue: "v1", LiveValue: "v2", Kind: drift.KindChanged}
	_ = d.Render(sampleDifferResult(ch))
	if !strings.Contains(buf.String(), "\033[") {
		t.Fatal("expected ANSI escape codes when color enabled")
	}
}

func TestDiffer_NilWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic.
	d := NewDiffer(DefaultDiffOptions(), nil)
	if d.out == nil {
		t.Fatal("expected non-nil writer")
	}
}
