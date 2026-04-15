package filter_test

import (
	"testing"

	"github.com/yourorg/driftctl-diff/internal/filter"
)

func TestAllow_NoRules_AllowsEverything(t *testing.T) {
	f := filter.New()
	if !f.Allow("aws_s3_bucket", "my-bucket") {
		t.Error("expected resource to be allowed when no rules are set")
	}
}

func TestAllow_IncludeRule_AllowsMatchingType(t *testing.T) {
	f := filter.New()
	f.AddInclude("aws_s3_bucket", "*")
	if !f.Allow("aws_s3_bucket", "any-bucket") {
		t.Error("expected aws_s3_bucket to be allowed")
	}
	if f.Allow("aws_instance", "i-123") {
		t.Error("expected aws_instance to be denied when not in include list")
	}
}

func TestAllow_ExcludeRule_DeniesMatchingResource(t *testing.T) {
	f := filter.New()
	f.AddExclude("aws_s3_bucket", "ignored-bucket")
	if f.Allow("aws_s3_bucket", "ignored-bucket") {
		t.Error("expected excluded resource to be denied")
	}
	if !f.Allow("aws_s3_bucket", "other-bucket") {
		t.Error("expected non-excluded resource to be allowed")
	}
}

func TestAllow_ExcludeTakesPrecedenceOverInclude(t *testing.T) {
	f := filter.New()
	f.AddInclude("aws_s3_bucket", "*")
	f.AddExclude("aws_s3_bucket", "sensitive-bucket")
	if f.Allow("aws_s3_bucket", "sensitive-bucket") {
		t.Error("exclude should take precedence over include")
	}
	if !f.Allow("aws_s3_bucket", "normal-bucket") {
		t.Error("non-excluded bucket should still be allowed")
	}
}

func TestAllow_WildcardPrefix_MatchesPrefix(t *testing.T) {
	f := filter.New()
	f.AddInclude("aws_iam_role", "prod-*")
	if !f.Allow("aws_iam_role", "prod-admin") {
		t.Error("expected prod-admin to match prod-* wildcard")
	}
	if f.Allow("aws_iam_role", "dev-admin") {
		t.Error("expected dev-admin to not match prod-* wildcard")
	}
}

func TestAllow_ExactIDMatch(t *testing.T) {
	f := filter.New()
	f.AddInclude("aws_instance", "i-0abc123")
	if !f.Allow("aws_instance", "i-0abc123") {
		t.Error("expected exact ID match to be allowed")
	}
	if f.Allow("aws_instance", "i-0abc124") {
		t.Error("expected different ID to be denied")
	}
}

func TestAllow_WildcardType(t *testing.T) {
	f := filter.New()
	f.AddExclude("*", "skip-me")
	if f.Allow("aws_s3_bucket", "skip-me") {
		t.Error("expected wildcard type exclude to apply to any resource type")
	}
	if !f.Allow("aws_s3_bucket", "keep-me") {
		t.Error("expected non-matching ID to be allowed")
	}
}
