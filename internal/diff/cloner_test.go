package diff

import (
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeCloneResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceID:   "vpc-abc123",
			ResourceType: "aws_vpc",
			Changes: []drift.AttributeChange{
				{Attribute: "cidr_block", StateValue: "10.0.0.0/16", LiveValue: "10.1.0.0/16", Kind: drift.KindChanged},
			},
			Metadata: map[string]string{"env": "prod"},
		},
		{
			ResourceID:   "sg-xyz789",
			ResourceType: "aws_security_group",
			Changes:      nil,
			Metadata:     nil,
		},
	}
}

func TestCloner_ProducesIndependentCopy(t *testing.T) {
	orig := makeCloneResults()
	c := NewCloner(DefaultClonerOptions())
	got, err := c.Clone(orig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Mutate original — clone must not be affected.
	orig[0].ResourceID = "mutated"
	orig[0].Changes[0].Attribute = "mutated"

	if got[0].ResourceID == "mutated" {
		t.Errorf("clone ResourceID was mutated by original change")
	}
}

func TestCloner_MetadataIsCopied(t *testing.T) {
	orig := makeCloneResults()
	c := NewCloner(DefaultClonerOptions())
	got, _ := c.Clone(orig)

	orig[0].Metadata["env"] = "staging"

	if got[0].Metadata["env"] != "prod" {
		t.Errorf("expected cloned metadata to remain \"prod\", got %q", got[0].Metadata["env"])
	}
}

func TestCloner_MetadataSharedWhenDisabled(t *testing.T) {
	orig := makeCloneResults()
	c := NewCloner(ClonerOptions{CopyMetadata: false})
	got, _ := c.Clone(orig)

	orig[0].Metadata["env"] = "staging"

	if got[0].Metadata["env"] != "staging" {
		t.Errorf("expected shared metadata to reflect mutation, got %q", got[0].Metadata["env"])
	}
}

func TestCloner_NilChanges_Preserved(t *testing.T) {
	orig := makeCloneResults()
	c := NewCloner(DefaultClonerOptions())
	got, _ := c.Clone(orig)

	if got[1].Changes != nil {
		t.Errorf("expected nil Changes to remain nil after clone")
	}
}

func TestCloner_EmptyResourceID_ReturnsError(t *testing.T) {
	invalid := []drift.DriftResult{{ResourceID: "", ResourceType: "aws_vpc"}}
	c := NewCloner(DefaultClonerOptions())
	_, err := c.Clone(invalid)
	if err == nil {
		t.Fatal("expected error for empty ResourceID, got nil")
	}
}

func TestCloner_EmptyInput_ReturnsEmpty(t *testing.T) {
	c := NewCloner(DefaultClonerOptions())
	got, err := c.Clone([]drift.DriftResult{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(got))
	}
}
