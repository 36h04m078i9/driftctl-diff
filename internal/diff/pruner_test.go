package diff

import (
	"testing"

	"github.com/driftctl-diff/internal/drift"
)

func makePrunerResults() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceID:   "bucket-1",
			ResourceType: "aws_s3_bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "acl", StateValue: "private", LiveValue: "public"},
			},
		},
		{
			ResourceID:   "sg-1",
			ResourceType: "aws_security_group",
			Changes:      []drift.AttributeChange{},
		},
		{
			ResourceID:   "role-1",
			ResourceType: "aws_iam_role",
			Changes: []drift.AttributeChange{
				{Attribute: "name", StateValue: "old", LiveValue: "new"},
				{Attribute: "path", StateValue: "/", LiveValue: "/svc/"},
			},
		},
	}
}

func TestPruner_RemoveUnchanged_DropsEmpty(t *testing.T) {
	p := NewPruner(DefaultPruneOptions())
	out := p.Prune(makePrunerResults())
	for _, r := range out {
		if len(r.Changes) == 0 {
			t.Errorf("expected no unchanged results, got %q", r.ResourceID)
		}
	}
}

func TestPruner_RemoveUnchanged_False_KeepsEmpty(t *testing.T) {
	opts := DefaultPruneOptions()
	opts.RemoveUnchanged = false
	p := NewPruner(opts)
	out := p.Prune(makePrunerResults())
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
}

func TestPruner_MinChanges_FiltersBelow(t *testing.T) {
	opts := DefaultPruneOptions()
	opts.RemoveUnchanged = false
	opts.MinChanges = 2
	p := NewPruner(opts)
	out := p.Prune(makePrunerResults())
	if len(out) != 1 {
		t.Fatalf("expected 1 result with >=2 changes, got %d", len(out))
	}
	if out[0].ResourceID != "role-1" {
		t.Errorf("expected role-1, got %s", out[0].ResourceID)
	}
}

func TestPruner_ExcludeTypes_DropsType(t *testing.T) {
	opts := DefaultPruneOptions()
	opts.ExcludeTypes = map[string]bool{"aws_s3_bucket": true}
	p := NewPruner(opts)
	out := p.Prune(makePrunerResults())
	for _, r := range out {
		if r.ResourceType == "aws_s3_bucket" {
			t.Errorf("expected aws_s3_bucket to be excluded")
		}
	}
}

func TestPruner_EmptyInput_ReturnsEmpty(t *testing.T) {
	p := NewPruner(DefaultPruneOptions())
	out := p.Prune([]drift.ResourceDiff{})
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d", len(out))
	}
}
