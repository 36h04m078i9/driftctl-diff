package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeCounterResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "versioning", Kind: drift.KindChanged},
				{Attribute: "tags", Kind: drift.KindAdded},
			},
		},
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "other-bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "acl", Kind: drift.KindDeleted},
			},
		},
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-1234",
			Changes: []drift.AttributeChange{
				{Attribute: "instance_type", Kind: drift.KindChanged},
			},
		},
	}
}

func TestCounter_EmptyInput_ReturnsNoDrift(t *testing.T) {
	c := NewCounter()
	c.Count(nil)
	var buf bytes.Buffer
	c.Print(&buf)
	if !strings.Contains(buf.String(), "no drift") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestCounter_TotalsByType(t *testing.T) {
	c := NewCounter()
	c.Count(makeCounterResults())
	entries := c.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 resource types, got %d", len(entries))
	}
	// sorted: aws_instance, aws_s3_bucket
	if entries[0].ResourceType != "aws_instance" {
		t.Errorf("unexpected first type: %s", entries[0].ResourceType)
	}
	if entries[1].Total != 2 {
		t.Errorf("expected 2 s3 resources, got %d", entries[1].Total)
	}
}

func TestCounter_ChangedCounts(t *testing.T) {
	c := NewCounter()
	c.Count(makeCounterResults())
	for _, e := range c.Entries() {
		if e.ResourceType == "aws_s3_bucket" {
			if e.Changed != 1 {
				t.Errorf("expected 1 changed, got %d", e.Changed)
			}
			if e.Added != 1 {
				t.Errorf("expected 1 added, got %d", e.Added)
			}
			if e.Deleted != 1 {
				t.Errorf("expected 1 deleted, got %d", e.Deleted)
			}
		}
	}
}

func TestCounter_Print_ContainsHeaders(t *testing.T) {
	c := NewCounter()
	c.Count(makeCounterResults())
	var buf bytes.Buffer
	c.Print(&buf)
	out := buf.String()
	for _, hdr := range []string{"RESOURCE TYPE", "TOTAL", "CHANGED", "ADDED", "DELETED"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q in output", hdr)
		}
	}
}

func TestCounter_Print_NilWriter_DefaultsToStdout(t *testing.T) {
	c := NewCounter()
	c.Count(makeCounterResults())
	// should not panic
	c.Print(nil)
}
