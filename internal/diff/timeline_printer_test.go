package diff

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/acme/driftctl-diff/internal/drift"
)

func TestTimelinePrinter_Empty(t *testing.T) {
	var buf bytes.Buffer
	p := NewTimelinePrinter(&buf)
	p.Print(NewTimeline())
	if !strings.Contains(buf.String(), "No timeline") {
		t.Errorf("expected no-entries message, got: %s", buf.String())
	}
}

func TestTimelinePrinter_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	p := NewTimelinePrinter(&buf)
	tl := NewTimeline()
	tl.Add(time.Now(), []drift.ResourceDiff{})
	p.Print(tl)
	out := buf.String()
	for _, h := range []string{"Captured At", "Drifted Resources", "Total Changes"} {
		if !strings.Contains(out, h) {
			t.Errorf("missing header %q", h)
		}
	}
}

func TestTimelinePrinter_CountsDriftedResources(t *testing.T) {
	var buf bytes.Buffer
	p := NewTimelinePrinter(&buf)
	tl := NewTimeline()
	results := []drift.ResourceDiff{
		{ResourceID: "i-1", Changes: []drift.AttributeChange{{Attribute: "ami"}},},
		{ResourceID: "i-2", Changes: []drift.AttributeChange{}},
	}
	tl.Add(time.Now(), results)
	p.Print(tl)
	out := buf.String()
	if !strings.Contains(out, "1") {
		t.Errorf("expected drifted count 1 in output: %s", out)
	}
}

func TestTimelinePrinter_NilWriter_DefaultsToStdout(t *testing.T) {
	p := NewTimelinePrinter(nil)
	if p.w == nil {
		t.Error("expected non-nil writer")
	}
}
