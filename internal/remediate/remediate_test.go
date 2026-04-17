package remediate_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/drift"
	"github.com/owner/driftctl-diff/internal/remediate"
)

func sampleDiffs() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-abc123",
			Changes: []drift.AttributeChange{
				{Attribute: "instance_type", StateValue: "t2.micro", LiveValue: "t3.small"},
			},
		},
	}
}

func TestSuggest_NoDrift_EmptySlice(t *testing.T) {
	results := remediate.Suggest([]drift.ResourceDiff{})
	if len(results) != 0 {
		t.Fatalf("expected 0 suggestions, got %d", len(results))
	}
}

func TestSuggest_WithDrift_ReturnsSuggestion(t *testing.T) {
	results := remediate.Suggest(sampleDiffs())
	if len(results) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(results))
	}
	if results[0].ResourceID != "i-abc123" {
		t.Errorf("unexpected resource id: %s", results[0].ResourceID)
	}
	if len(results[0].Lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(results[0].Lines))
	}
}

func TestSuggest_LineContainsAttribute(t *testing.T) {
	results := remediate.Suggest(sampleDiffs())
	if !strings.Contains(results[0].Lines[0], "instance_type") {
		t.Errorf("expected line to mention attribute, got: %s", results[0].Lines[0])
	}
}

func TestPrint_NoDrift_PrintsSyncMessage(t *testing.T) {
	var buf bytes.Buffer
	s := remediate.New(&buf)
	s.Print(nil)
	if !strings.Contains(buf.String(), "in sync") {
		t.Errorf("expected sync message, got: %s", buf.String())
	}
}

func TestPrint_WithSuggestions_ContainsResourceID(t *testing.T) {
	var buf bytes.Buffer
	s := remediate.New(&buf)
	s.Print(remediate.Suggest(sampleDiffs()))
	if !strings.Contains(buf.String(), "i-abc123") {
		t.Errorf("expected resource id in output, got: %s", buf.String())
	}
}

func TestNew_NilWriter_DefaultsToStdout(t *testing.T) {
	s := remediate.New(nil)
	if s == nil {
		t.Fatal("expected non-nil Suggester")
	}
}
