package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makePivotResults() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceID:   "bucket-1",
			ResourceType: "aws_s3_bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "region", StateValue: "us-east-1", LiveValue: "eu-west-1"},
				{Attribute: "versioning", StateValue: "true", LiveValue: "false"},
			},
		},
		{
			ResourceID:   "bucket-2",
			ResourceType: "aws_s3_bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "region", StateValue: "us-east-1", LiveValue: "ap-south-1"},
			},
		},
	}
}

func TestPivot_AttributesAreSorted(t *testing.T) {
	p := NewPivot()
	entries := p.Compute(makePivotResults())
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Attribute != "region" {
		t.Errorf("expected first attribute 'region', got %q", entries[0].Attribute)
	}
	if entries[1].Attribute != "versioning" {
		t.Errorf("expected second attribute 'versioning', got %q", entries[1].Attribute)
	}
}

func TestPivot_ResourcesGroupedUnderAttribute(t *testing.T) {
	p := NewPivot()
	entries := p.Compute(makePivotResults())
	regionEntry := entries[0]
	if len(regionEntry.Resources) != 2 {
		t.Errorf("expected 2 resources under 'region', got %d", len(regionEntry.Resources))
	}
}

func TestPivot_EmptyResults_ReturnsEmpty(t *testing.T) {
	p := NewPivot()
	entries := p.Compute(nil)
	if len(entries) != 0 {
		t.Errorf("expected empty pivot, got %d entries", len(entries))
	}
}

func TestPivotPrinter_NoDrift_PrintsMessage(t *testing.T) {
	var buf bytes.Buffer
	pp := NewPivotPrinter(&buf)
	pp.Print(nil)
	if !strings.Contains(buf.String(), "no attribute-level drift") {
		t.Errorf("expected no-drift message, got %q", buf.String())
	}
}

func TestPivotPrinter_WithEntries_ContainsAttribute(t *testing.T) {
	p := NewPivot()
	entries := p.Compute(makePivotResults())
	var buf bytes.Buffer
	pp := NewPivotPrinter(&buf)
	pp.Print(entries)
	out := buf.String()
	if !strings.Contains(out, "region") {
		t.Errorf("expected 'region' in output, got %q", out)
	}
	if !strings.Contains(out, "bucket-1") {
		t.Errorf("expected 'bucket-1' in output, got %q", out)
	}
}

func TestPivotPrinter_NilWriter_DefaultsToStdout(t *testing.T) {
	// Just ensure no panic.
	pp := NewPivotPrinter(nil)
	pp.Print(nil)
}
