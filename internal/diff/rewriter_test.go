package diff

import (
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeRewriterResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceID:   "bucket-1",
			ResourceType: "aws_s3_bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "region", Kind: drift.KindChanged, WantValue: "us-east-1", GotValue: "eu-west-1"},
				{Attribute: "name", Kind: drift.KindChanged, WantValue: "prod-bucket", GotValue: "staging-bucket"},
			},
		},
		{
			ResourceID:   "sg-42",
			ResourceType: "aws_security_group",
			Changes: []drift.AttributeChange{
				{Attribute: "description", Kind: drift.KindChanged, WantValue: "prod env", GotValue: "staging env"},
			},
		},
	}
}

func TestRewriter_NoRules_ReturnsUnchanged(t *testing.T) {
	r := NewRewriter(DefaultRewriteOptions())
	input := makeRewriterResults()
	out := r.Rewrite(input)
	if len(out) != len(input) {
		t.Fatalf("expected %d results, got %d", len(input), len(out))
	}
	if out[0].Changes[0].GotValue != "eu-west-1" {
		t.Errorf("expected eu-west-1, got %s", out[0].Changes[0].GotValue)
	}
}

func TestRewriter_WildcardRule_ReplacesAllMatchingValues(t *testing.T) {
	opts := RewriteOptions{
		Rules: []RewriteRule{
			{ResourceType: "*", Attribute: "*", Find: "staging", Replace: "production"},
		},
	}
	r := NewRewriter(opts)
	out := r.Rewrite(makeRewriterResults())
	if out[0].Changes[1].GotValue != "production-bucket" {
		t.Errorf("unexpected value: %s", out[0].Changes[1].GotValue)
	}
	if out[1].Changes[0].GotValue != "production env" {
		t.Errorf("unexpected value: %s", out[1].Changes[0].GotValue)
	}
}

func TestRewriter_TypeScopedRule_OnlyAffectsMatchingType(t *testing.T) {
	opts := RewriteOptions{
		Rules: []RewriteRule{
			{ResourceType: "aws_s3_bucket", Attribute: "*", Find: "staging", Replace: "REPLACED"},
		},
	}
	r := NewRewriter(opts)
	out := r.Rewrite(makeRewriterResults())
	// s3 bucket name should be rewritten
	if out[0].Changes[1].GotValue != "REPLACED-bucket" {
		t.Errorf("expected REPLACED-bucket, got %s", out[0].Changes[1].GotValue)
	}
	// security group description should NOT be rewritten
	if out[1].Changes[0].GotValue != "staging env" {
		t.Errorf("expected staging env unchanged, got %s", out[1].Changes[0].GotValue)
	}
}

func TestRewriter_AttributeScopedRule_OnlyAffectsMatchingAttribute(t *testing.T) {
	opts := RewriteOptions{
		Rules: []RewriteRule{
			{ResourceType: "*", Attribute: "region", Find: "eu-west-1", Replace: "us-east-1"},
		},
	}
	r := NewRewriter(opts)
	out := r.Rewrite(makeRewriterResults())
	if out[0].Changes[0].GotValue != "us-east-1" {
		t.Errorf("expected us-east-1, got %s", out[0].Changes[0].GotValue)
	}
	// name attribute should be untouched
	if out[0].Changes[1].GotValue != "staging-bucket" {
		t.Errorf("expected staging-bucket unchanged, got %s", out[0].Changes[1].GotValue)
	}
}

func TestRewriter_DoesNotMutateInput(t *testing.T) {
	opts := RewriteOptions{
		Rules: []RewriteRule{
			{ResourceType: "*", Attribute: "*", Find: "eu-west-1", Replace: "MUTATED"},
		},
	}
	r := NewRewriter(opts)
	input := makeRewriterResults()
	original := input[0].Changes[0].GotValue
	r.Rewrite(input)
	if input[0].Changes[0].GotValue != original {
		t.Errorf("input was mutated: got %s", input[0].Changes[0].GotValue)
	}
}

func TestRewriter_PrefixWildcard_MatchesResourceType(t *testing.T) {
	opts := RewriteOptions{
		Rules: []RewriteRule{
			{ResourceType: "aws_s3*", Attribute: "*", Find: "eu-west-1", Replace: "ap-southeast-1"},
		},
	}
	r := NewRewriter(opts)
	out := r.Rewrite(makeRewriterResults())
	if out[0].Changes[0].GotValue != "ap-southeast-1" {
		t.Errorf("expected ap-southeast-1, got %s", out[0].Changes[0].GotValue)
	}
	// security group should be unaffected
	if out[1].Changes[0].GotValue != "staging env" {
		t.Errorf("expected staging env, got %s", out[1].Changes[0].GotValue)
	}
}
