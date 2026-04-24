package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/diff"
	"github.com/acme/driftctl-diff/internal/drift"
)

func TestTransformerPrinter_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	p := diff.NewTransformerPrinter(&buf)
	p.Print(makeTransformerResults(), makeTransformerResults())
	out := buf.String()
	for _, hdr := range []string{"FIELD", "BEFORE", "AFTER"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestTransformerPrinter_CountsResources(t *testing.T) {
	var buf bytes.Buffer
	p := diff.NewTransformerPrinter(&buf)
	before := makeTransformerResults()
	after := []drift.DriftResult{}
	p.Print(before, after)
	out := buf.String()
	if !strings.Contains(out, "1") {
		t.Error("expected before count 1 in output")
	}
}

func TestTransformerPrinter_NilWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic
	p := diff.NewTransformerPrinter(nil)
	p.Print(nil, nil)
}

func TestTransformerPrinter_ZeroResults_ShowsZero(t *testing.T) {
	var buf bytes.Buffer
	p := diff.NewTransformerPrinter(&buf)
	p.Print(nil, nil)
	out := buf.String()
	if !strings.Contains(out, "0") {
		t.Error("expected zero counts in output")
	}
}
