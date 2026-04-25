package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func TestRewriterPrinter_NoDrift_PrintsMessage(t *testing.T) {
	var buf bytes.Buffer
	p := NewRewriterPrinter(&buf)
	p.Print(nil)
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestRewriterPrinter_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	p := NewRewriterPrinter(&buf)
	p.Print(makeRewriterResults())
	out := buf.String()
	for _, header := range []string{"RESOURCE TYPE", "RESOURCE ID", "ATTRIBUTE", "WANT", "GOT"} {
		if !strings.Contains(out, header) {
			t.Errorf("missing header %q in output", header)
		}
	}
}

func TestRewriterPrinter_ContainsResourceInfo(t *testing.T) {
	var buf bytes.Buffer
	p := NewRewriterPrinter(&buf)
	p.Print(makeRewriterResults())
	out := buf.String()
	if !strings.Contains(out, "aws_s3_bucket") {
		t.Errorf("expected resource type in output")
	}
	if !strings.Contains(out, "bucket-1") {
		t.Errorf("expected resource ID in output")
	}
	if !strings.Contains(out, "eu-west-1") {
		t.Errorf("expected attribute value in output")
	}
}

func TestRewriterPrinter_NilWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic.
	p := NewRewriterPrinter(nil)
	p.Print([]drift.DriftResult{})
}
