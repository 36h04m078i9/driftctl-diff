package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/drift"
)

func makeLinkerResults(specs []struct {
	id  string
	attr string
	val string
}) []drift.ResourceDiff {
	var out []drift.ResourceDiff
	for _, s := range specs {
		out = append(out, drift.ResourceDiff{
			ResourceID:   s.id,
			ResourceType: "aws_instance",
			Changes: []drift.AttributeChange{
				{Attribute: s.attr, LiveValue: s.val, StateValue: "old"},
			},
		})
	}
	return out
}

func TestLinker_NoAttribute_ReturnsNoLinks(t *testing.T) {
	results := makeLinkerResults([]struct{ id, attr, val string }{
		{"r1", "ami", "ami-abc"},
		{"r2", "ami", "ami-abc"},
	})
	l := NewLinker(LinkerOptions{LinkByAttribute: "", MaxLinks: 5})
	lr := l.Link(results)
	if len(lr.Links) != 0 {
		t.Fatalf("expected 0 links, got %d", len(lr.Links))
	}
}

func TestLinker_SharedAttribute_ProducesLink(t *testing.T) {
	results := makeLinkerResults([]struct{ id, attr, val string }{
		{"r1", "subnet_id", "subnet-123"},
		{"r2", "subnet_id", "subnet-123"},
	})
	l := NewLinker(LinkerOptions{LinkByAttribute: "subnet_id", MaxLinks: 10})
	lr := l.Link(results)
	if len(lr.Links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(lr.Links))
	}
	if lr.Links[0].SharedAttribute != "subnet_id" {
		t.Errorf("unexpected shared attribute: %s", lr.Links[0].SharedAttribute)
	}
}

func TestLinker_NoSharedValues_NoLinks(t *testing.T) {
	results := makeLinkerResults([]struct{ id, attr, val string }{
		{"r1", "subnet_id", "subnet-aaa"},
		{"r2", "subnet_id", "subnet-bbb"},
	})
	l := NewLinker(DefaultLinkerOptions())
	l.opts.LinkByAttribute = "subnet_id"
	lr := l.Link(results)
	if len(lr.Links) != 0 {
		t.Fatalf("expected 0 links, got %d", len(lr.Links))
	}
}

func TestLinker_EmptyResults_ReturnsEmpty(t *testing.T) {
	l := NewLinker(LinkerOptions{LinkByAttribute: "vpc_id", MaxLinks: 5})
	lr := l.Link(nil)
	if len(lr.Links) != 0 {
		t.Fatalf("expected 0 links for empty input")
	}
}

func TestLinkerPrinter_NoDrift_PrintsNoLinks(t *testing.T) {
	var buf bytes.Buffer
	p := NewLinkerPrinter(&buf)
	p.Print(LinkerResult{})
	if !strings.Contains(buf.String(), "No links found.") {
		t.Errorf("expected 'No links found.' in output, got: %s", buf.String())
	}
}

func TestLinkerPrinter_WithLinks_ContainsIDs(t *testing.T) {
	var buf bytes.Buffer
	p := NewLinkerPrinter(&buf)
	p.Print(LinkerResult{
		Links: []Link{{SourceID: "r1", TargetID: "r2", SharedAttribute: "vpc_id"}},
	})
	out := buf.String()
	if !strings.Contains(out, "r1") || !strings.Contains(out, "r2") {
		t.Errorf("expected resource IDs in output, got: %s", out)
	}
	if !strings.Contains(out, "vpc_id") {
		t.Errorf("expected shared attribute in output, got: %s", out)
	}
}

func TestLinkerPrinter_NilWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic
	p := NewLinkerPrinter(nil)
	if p == nil {
		t.Fatal("expected non-nil printer")
	}
}
