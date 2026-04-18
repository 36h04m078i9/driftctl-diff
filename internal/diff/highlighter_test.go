package diff

import (
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func sampleChange(attr, old, new_ string) drift.AttributeChange {
	return drift.AttributeChange{Attribute: attr, OldValue: old, NewValue: new_}
}

func TestHighlight_NoColor_ContainsBothValues(t *testing.T) {
	h := NewHighlighter(false)
	out := h.Highlight(sampleChange("instance_type", "t2.micro", "t3.small"))
	if !strings.Contains(out, "t2.micro") {
		t.Errorf("expected old value in output, got: %s", out)
	}
	if !strings.Contains(out, "t3.small") {
		t.Errorf("expected new value in output, got: %s", out)
	}
}

func TestHighlight_NoColor_NoEscapeCodes(t *testing.T) {
	h := NewHighlighter(false)
	out := h.Highlight(sampleChange("ami", "ami-old", "ami-new"))
	if strings.Contains(out, "\033[") {
		t.Errorf("did not expect ANSI codes when color disabled, got: %q", out)
	}
}

func TestHighlight_Color_ContainsEscapeCodes(t *testing.T) {
	h := NewHighlighter(true)
	out := h.Highlight(sampleChange("ami", "ami-old", "ami-new"))
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI codes when color enabled, got: %q", out)
	}
}

func TestHighlight_PrefixMinusAndPlus(t *testing.T) {
	h := NewHighlighter(false)
	out := h.Highlight(sampleChange("size", "10", "20"))
	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "-") {
		t.Errorf("first line should start with '-', got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[1], "+") {
		t.Errorf("second line should start with '+', got: %s", lines[1])
	}
}

func TestHighlightAll_ContainsResourceHeader(t *testing.T) {
	h := NewHighlighter(false)
	r := drift.DriftResult{
		ResourceType: "aws_instance",
		ResourceID:   "i-1234",
		Changes:      []drift.AttributeChange{sampleChange("type", "t2.micro", "t3.small")},
	}
	out := h.HighlightAll(r)
	if !strings.Contains(out, "aws_instance") {
		t.Errorf("expected resource type in header, got: %s", out)
	}
	if !strings.Contains(out, "i-1234") {
		t.Errorf("expected resource id in header, got: %s", out)
	}
}

func TestHighlightAll_NoChanges_OnlyHeader(t *testing.T) {
	h := NewHighlighter(false)
	r := drift.DriftResult{ResourceType: "aws_s3_bucket", ResourceID: "my-bucket", Changes: nil}
	out := h.HighlightAll(r)
	lines := strings.Split(out, "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header line, got %d lines: %v", len(lines), lines)
	}
}
