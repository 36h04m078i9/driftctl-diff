package diff_test

import (
	"testing"

	"github.com/acme/driftctl-diff/internal/diff"
	"github.com/acme/driftctl-diff/internal/drift"
)

func makeTransformerResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceID:   "bucket-1",
			ResourceType: "aws_s3_bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "acl", StateValue: "private", LiveValue: "public-read", Kind: drift.KindChanged},
				{Attribute: "region", StateValue: "", LiveValue: "", Kind: drift.KindChanged},
			},
		},
	}
}

func TestTransformer_NoOptions_ReturnsUnchanged(t *testing.T) {
	results := makeTransformerResults()
	tr := diff.NewTransformer(diff.DefaultTransformOptions())
	out := tr.Transform(results)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].ResourceID != "bucket-1" {
		t.Errorf("unexpected resource ID: %s", out[0].ResourceID)
	}
}

func TestTransformer_PrefixResourceIDs(t *testing.T) {
	results := makeTransformerResults()
	opts := diff.DefaultTransformOptions()
	opts.PrefixResourceIDs = "prod/"
	out := diff.NewTransformer(opts).Transform(results)
	if out[0].ResourceID != "prod/bucket-1" {
		t.Errorf("expected prefixed ID, got %s", out[0].ResourceID)
	}
}

func TestTransformer_SuffixResourceTypes(t *testing.T) {
	results := makeTransformerResults()
	opts := diff.DefaultTransformOptions()
	opts.SuffixResourceTypes = "_v2"
	out := diff.NewTransformer(opts).Transform(results)
	if out[0].ResourceType != "aws_s3_bucket_v2" {
		t.Errorf("unexpected type: %s", out[0].ResourceType)
	}
}

func TestTransformer_UpperCaseAttributes(t *testing.T) {
	results := makeTransformerResults()
	opts := diff.DefaultTransformOptions()
	opts.UpperCaseAttributes = true
	out := diff.NewTransformer(opts).Transform(results)
	for _, ch := range out[0].Changes {
		if ch.Attribute != "ACL" && ch.Attribute != "REGION" {
			t.Errorf("expected upper-case attribute, got %s", ch.Attribute)
		}
	}
}

func TestTransformer_DropEmptyValues(t *testing.T) {
	results := makeTransformerResults()
	opts := diff.DefaultTransformOptions()
	opts.DropEmptyValues = true
	out := diff.NewTransformer(opts).Transform(results)
	// "region" change has both empty values and should be dropped
	if len(out[0].Changes) != 1 {
		t.Errorf("expected 1 change after drop, got %d", len(out[0].Changes))
	}
	if out[0].Changes[0].Attribute != "acl" {
		t.Errorf("expected acl change to survive, got %s", out[0].Changes[0].Attribute)
	}
}

func TestTransformer_DoesNotMutateOriginal(t *testing.T) {
	results := makeTransformerResults()
	opts := diff.DefaultTransformOptions()
	opts.PrefixResourceIDs = "x/"
	diff.NewTransformer(opts).Transform(results)
	if results[0].ResourceID != "bucket-1" {
		t.Error("original results were mutated")
	}
}
