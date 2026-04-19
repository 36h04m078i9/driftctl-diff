package diff

import (
	"testing"

	"github.com/driftctl-diff/internal/drift"
)

func makeFlatResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-001",
			Changes: []drift.AttributeChange{
				{Attribute: "ami", Kind: drift.KindChanged, WantValue: "ami-old", GotValue: "ami-new"},
				{Attribute: "tags", Kind: drift.KindMissing, WantValue: "env=prod", GotValue: ""},
			},
		},
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "acl", Kind: drift.KindChanged, WantValue: "private", GotValue: "public-read"},
			},
		},
	}
}

func TestFlattener_TotalCount(t *testing.T) {
	f := NewFlattener()
	flat := f.Flatten(makeFlatResults())
	if len(flat) != 3 {
		t.Fatalf("expected 3 flat changes, got %d", len(flat))
	}
}

func TestFlattener_PreservesResourceContext(t *testing.T) {
	f := NewFlattener()
	flat := f.Flatten(makeFlatResults())
	if flat[0].ResourceType != "aws_instance" {
		t.Errorf("expected aws_instance, got %s", flat[0].ResourceType)
	}
	if flat[0].ResourceID != "i-001" {
		t.Errorf("expected i-001, got %s", flat[0].ResourceID)
	}
}

func TestFlattener_EmptyResults(t *testing.T) {
	f := NewFlattener()
	flat := f.Flatten(nil)
	if len(flat) != 0 {
		t.Errorf("expected empty slice, got %d", len(flat))
	}
}

func TestFlattener_FlattenByKind_Changed(t *testing.T) {
	f := NewFlattener()
	flat := f.FlattenByKind(makeFlatResults(), drift.KindChanged)
	if len(flat) != 2 {
		t.Fatalf("expected 2 changed, got %d", len(flat))
	}
	for _, fc := range flat {
		if fc.Kind != drift.KindChanged {
			t.Errorf("expected KindChanged, got %v", fc.Kind)
		}
	}
}

func TestFlattener_FlattenByKind_Missing(t *testing.T) {
	f := NewFlattener()
	flat := f.FlattenByKind(makeFlatResults(), drift.KindMissing)
	if len(flat) != 1 {
		t.Fatalf("expected 1 missing, got %d", len(flat))
	}
	if flat[0].Attribute != "tags" {
		t.Errorf("expected tags attribute, got %s", flat[0].Attribute)
	}
}
