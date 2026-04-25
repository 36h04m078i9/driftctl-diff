package diff

import (
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeCensorResults() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceID:   "aws_instance.web",
			ResourceType: "aws_instance",
			Changes: []drift.Change{
				{Attribute: "instance_type", Got: "t2.micro", Want: "t3.micro", Kind: drift.KindChanged},
				{Attribute: "db_password", Got: "hunter2", Want: "s3cr3t", Kind: drift.KindChanged},
				{Attribute: "api_token", Got: "tok_live_abc", Want: "", Kind: drift.KindMissing},
			},
		},
		{
			ResourceID:   "aws_s3_bucket.assets",
			ResourceType: "aws_s3_bucket",
			Changes: []drift.Change{
				{Attribute: "bucket", Got: "my-bucket", Want: "your-bucket", Kind: drift.KindChanged},
			},
		},
	}
}

func TestCensor_SensitiveAttributesAreRedacted(t *testing.T) {
	c := NewCensor(DefaultCensorOptions())
	out := c.Apply(makeCensorResults())

	changes := out[0].Changes
	for _, ch := range changes {
		if ch.Attribute == "db_password" || ch.Attribute == "api_token" {
			if ch.Got != "[REDACTED]" || ch.Want != "[REDACTED]" {
				t.Errorf("expected [REDACTED] for %s, got Got=%q Want=%q", ch.Attribute, ch.Got, ch.Want)
			}
		}
	}
}

func TestCensor_NonSensitiveAttributesUnchanged(t *testing.T) {
	c := NewCensor(DefaultCensorOptions())
	out := c.Apply(makeCensorResults())

	for _, ch := range out[0].Changes {
		if ch.Attribute == "instance_type" {
			if ch.Got != "t2.micro" {
				t.Errorf("expected original value, got %q", ch.Got)
			}
		}
	}
}

func TestCensor_DoesNotMutateInput(t *testing.T) {
	original := makeCensorResults()
	c := NewCensor(DefaultCensorOptions())
	c.Apply(original)

	for _, ch := range original[0].Changes {
		if ch.Attribute == "db_password" && ch.Got == "[REDACTED]" {
			t.Error("Apply must not mutate the input slice")
		}
	}
}

func TestCensor_EmptyResults_ReturnsEmpty(t *testing.T) {
	c := NewCensor(DefaultCensorOptions())
	out := c.Apply([]drift.ResourceDiff{})
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(out))
	}
}

func TestCensor_CustomReplacement(t *testing.T) {
	opts := CensorOptions{
		Patterns:    []string{"secret"},
		Replacement: "***",
	}
	c := NewCensor(opts)
	input := []drift.ResourceDiff{
		{
			ResourceID:   "res",
			ResourceType: "aws_rds",
			Changes: []drift.Change{
				{Attribute: "master_secret", Got: "abc", Want: "xyz", Kind: drift.KindChanged},
			},
		},
	}
	out := c.Apply(input)
	if out[0].Changes[0].Got != "***" {
		t.Errorf("expected *** replacement, got %q", out[0].Changes[0].Got)
	}
}
