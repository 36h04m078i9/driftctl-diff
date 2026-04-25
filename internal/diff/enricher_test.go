package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeEnricherResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.Change{
				{Attribute: "acl", StateValue: "private", LiveValue: "public-read"},
				{Attribute: "versioning", StateValue: "true", LiveValue: "false"},
			},
		},
		{
			ResourceType: "aws_security_group",
			ResourceID:   "sg-abc123",
			Changes:      []drift.Change{},
		},
	}
}

func TestEnricher_ChangeCountAttached(t *testing.T) {
	e := NewEnricher(DefaultEnrichOptions())
	results := e.Enrich(makeEnricherResults())
	if results[0].ChangeCount != 2 {
		t.Fatalf("expected 2 changes, got %d", results[0].ChangeCount)
	}
	if results[1].ChangeCount != 0 {
		t.Fatalf("expected 0 changes, got %d", results[1].ChangeCount)
	}
}

func TestEnricher_ChangeCountMetadata(t *testing.T) {
	e := NewEnricher(DefaultEnrichOptions())
	results := e.Enrich(makeEnricherResults())
	if results[0].Metadata["change_count"] != "2" {
		t.Fatalf("expected metadata change_count=2, got %q", results[0].Metadata["change_count"])
	}
}

func TestEnricher_ResourceURL_WhenEnabled(t *testing.T) {
	opts := DefaultEnrichOptions()
	opts.AddResourceURL = true
	e := NewEnricher(opts)
	results := e.Enrich(makeEnricherResults())
	if !strings.Contains(results[0].ResourceURL, "aws-s3-bucket") {
		t.Fatalf("expected URL to contain resource type, got %q", results[0].ResourceURL)
	}
	if !strings.Contains(results[0].ResourceURL, "my-bucket") {
		t.Fatalf("expected URL to contain resource ID, got %q", results[0].ResourceURL)
	}
}

func TestEnricher_ResourceURL_WhenDisabled(t *testing.T) {
	opts := DefaultEnrichOptions()
	opts.AddResourceURL = false
	e := NewEnricher(opts)
	results := e.Enrich(makeEnricherResults())
	if results[0].ResourceURL != "" {
		t.Fatalf("expected empty URL, got %q", results[0].ResourceURL)
	}
}

func TestEnricher_EmptyInput_ReturnsEmpty(t *testing.T) {
	e := NewEnricher(DefaultEnrichOptions())
	results := e.Enrich(nil)
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestEnricherPrinter_NoDrift_PrintsMessage(t *testing.T) {
	var buf bytes.Buffer
	p := NewEnricherPrinter(&buf)
	p.Print(nil)
	if !strings.Contains(buf.String(), "No enriched") {
		t.Fatalf("expected no-drift message, got %q", buf.String())
	}
}

func TestEnricherPrinter_ContainsHeaders(t *testing.T) {
	e := NewEnricher(DefaultEnrichOptions())
	enriched := e.Enrich(makeEnricherResults())
	var buf bytes.Buffer
	p := NewEnricherPrinter(&buf)
	p.Print(enriched)
	output := buf.String()
	for _, hdr := range []string{"RESOURCE TYPE", "RESOURCE ID", "CHANGES"} {
		if !strings.Contains(output, hdr) {
			t.Fatalf("expected header %q in output, got:\n%s", hdr, output)
		}
	}
}

func TestEnricherPrinter_NilWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic.
	p := NewEnricherPrinter(nil)
	_ = p
}
