package diff

import (
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeDriftForAnnotation() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-abc123",
			Changes: []drift.AttributeChange{
				{Attribute: "instance_type", Kind: drift.KindChanged, StateValue: "t2.micro", LiveValue: "t3.small"},
				{Attribute: "ami", Kind: drift.KindChanged, StateValue: "ami-old", LiveValue: "ami-new"},
				{Attribute: "key_name", Kind: drift.KindChanged, StateValue: "old-key", LiveValue: "new-key"},
			},
		},
	}
}

func TestAnnotate_ReturnsOneAnnotationPerChange(t *testing.T) {
	a := NewAnnotator()
	results := makeDriftForAnnotation()
	anns := a.Annotate(results)
	if len(anns) != 3 {
		t.Fatalf("expected 3 annotations, got %d", len(anns))
	}
}

func TestAnnotate_KnownAttribute_UsesDefaultNote(t *testing.T) {
	a := NewAnnotator()
	results := makeDriftForAnnotation()
	anns := a.Annotate(results)
	var found *Annotation
	for i := range anns {
		if anns[i].Attribute == "instance_type" {
			found = &anns[i]
			break
		}
	}
	if found == nil {
		t.Fatal("expected annotation for instance_type")
	}
	if !strings.Contains(found.Note, "downtime") {
		t.Errorf("expected downtime note, got: %s", found.Note)
	}
}

func TestAnnotate_UnknownAttribute_FallbackNote(t *testing.T) {
	a := NewAnnotator()
	results := makeDriftForAnnotation()
	anns := a.Annotate(results)
	var found *Annotation
	for i := range anns {
		if anns[i].Attribute == "key_name" {
			found = &anns[i]
			break
		}
	}
	if found == nil {
		t.Fatal("expected annotation for key_name")
	}
	if !strings.Contains(found.Note, "drifted") {
		t.Errorf("expected fallback note, got: %s", found.Note)
	}
}

func TestAnnotate_EmptyResults_NoAnnotations(t *testing.T) {
	a := NewAnnotator()
	anns := a.Annotate(nil)
	if len(anns) != 0 {
		t.Errorf("expected 0 annotations, got %d", len(anns))
	}
}

func TestAnnotate_SetsResourceFields(t *testing.T) {
	a := NewAnnotator()
	results := makeDriftForAnnotation()
	anns := a.Annotate(results)
	for _, ann := range anns {
		if ann.ResourceType != "aws_instance" {
			t.Errorf("unexpected resource type: %s", ann.ResourceType)
		}
		if ann.ResourceID != "i-abc123" {
			t.Errorf("unexpected resource id: %s", ann.ResourceID)
		}
	}
}
