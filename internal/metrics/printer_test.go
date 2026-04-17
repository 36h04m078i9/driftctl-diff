package metrics

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestPrinter_ContainsAllFields(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	c := Counters{
		ResourcesTotal:    20,
		ResourcesDrifted:  4,
		AttributesChecked: 80,
		FetchErrors:       1,
		Duration:          123 * time.Millisecond,
	}
	p.Print(c)
	out := buf.String()
	for _, want := range []string{"20", "4", "80", "1", "123ms"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestPrinter_NilWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic when w is nil.
	p := NewPrinter(nil)
	if p.w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestPrinter_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	p.Print(Counters{})
	if !strings.Contains(buf.String(), "0") {
		t.Errorf("expected zeros in output, got: %s", buf.String())
	}
}
