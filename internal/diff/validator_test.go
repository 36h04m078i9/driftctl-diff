package diff

import (
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeValidatorResults() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "acl", StateValue: "private", LiveValue: "public-read", Kind: drift.KindChanged},
			},
		},
	}
}

func TestValidator_ValidResults_NoErrors(t *testing.T) {
	v := NewValidator(100)
	vr, err := v.Validate(makeValidatorResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !vr.Valid {
		t.Errorf("expected valid, got errors: %v", vr.Errors)
	}
}

func TestValidator_NilResults_ReturnsError(t *testing.T) {
	v := NewValidator(100)
	_, err := v.Validate(nil)
	if err == nil {
		t.Fatal("expected error for nil results")
	}
}

func TestValidator_MissingResourceID_Invalid(t *testing.T) {
	v := NewValidator(100)
	results := []drift.ResourceDiff{{ResourceType: "aws_s3_bucket", ResourceID: ""}}
	vr, err := v.Validate(results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vr.Valid {
		t.Error("expected invalid result")
	}
	if len(vr.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(vr.Errors))
	}
}

func TestValidator_MissingResourceType_Invalid(t *testing.T) {
	v := NewValidator(100)
	results := []drift.ResourceDiff{{ResourceType: "", ResourceID: "abc"}}
	vr, _ := v.Validate(results)
	if vr.Valid {
		t.Error("expected invalid")
	}
}

func TestValidator_TooManyChanges_AddsWarning(t *testing.T) {
	v := NewValidator(2)
	changes := make([]drift.AttributeChange, 5)
	results := []drift.ResourceDiff{{ResourceType: "aws_instance", ResourceID: "i-123", Changes: changes}}
	vr, _ := v.Validate(results)
	if !vr.Valid {
		t.Error("should still be valid")
	}
	if len(vr.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(vr.Warnings))
	}
}

func TestValidator_EmptyResults_Valid(t *testing.T) {
	v := NewValidator(0)
	vr, err := v.Validate([]drift.ResourceDiff{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !vr.Valid {
		t.Error("empty results should be valid")
	}
}
