package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/snyk/driftctl-diff/internal/drift"
	"github.com/snyk/driftctl-diff/internal/output"
)

func newTestWriter(t *testing.T, color bool) (*output.Writer, *bytes.Buffer) {
	t.Helper()
	buf := &bytes.Buffer{}
	colorizer := output.NewColorizer(color)
	f := output.NewFormatter(colorizer)
	w := output.NewWriter(f, buf)
	return w, buf
}

func TestWriter_NoDrift_PrintsNoDriftMessage(t *testing.T) {
	w, buf := newTestWriter(t, false)

	_, err := w.Write([]drift.ResourceDiff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "No drift detected") {
		t.Errorf("expected 'No drift detected' in output, got: %q", got)
	}
}

func TestWriter_WithDrift_ContainsResourceID(t *testing.T) {
	w, buf := newTestWriter(t, false)

	results := []drift.ResourceDiff{
		{
			ResourceID:   "aws_s3_bucket.my-bucket",
			ResourceType: "aws_s3_bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "versioning", StateValue: "disabled", LiveValue: "enabled"},
			},
		},
	}

	_, err := w.Write(results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "aws_s3_bucket.my-bucket") {
		t.Errorf("expected resource ID in output, got: %q", got)
	}
}

func TestWriter_NilDest_DefaultsToStdout(t *testing.T) {
	colorizer := output.NewColorizer(false)
	f := output.NewFormatter(colorizer)
	// Should not panic with nil dest
	w := output.NewWriter(f, nil)
	if w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestWriter_WriteTo_AlternativeDest(t *testing.T) {
	w, _ := newTestWriter(t, false)
	altBuf := &bytes.Buffer{}

	results := []drift.ResourceDiff{}
	_, err := w.WriteTo(results, altBuf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(altBuf.String(), "No drift detected") {
		t.Errorf("expected output in alternative dest, got: %q", altBuf.String())
	}
}
