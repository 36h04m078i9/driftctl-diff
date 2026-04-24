package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeWeigherResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.Change{
				{Attribute: "acl", Kind: drift.Changed, Got: "public", Want: "private"},
			},
		},
		{
			ResourceType: "aws_iam_user",
			ResourceID:   "alice",
			Changes: []drift.Change{
				{Attribute: "password", Kind: drift.Changed, Got: "x", Want: "y"},
				{Attribute: "email", Kind: drift.Changed, Got: "a", Want: "b"},
			},
		},
		{
			ResourceType: "aws_security_group",
			ResourceID:   "sg-001",
			Changes: []drift.Change{
				{Attribute: "port", Kind: drift.Missing},
			},
		},
	}
}

func TestWeigher_SortedDescending(t *testing.T) {
	w := NewWeigher(DefaultWeighOptions())
	weighed := w.Weigh(makeWeigherResults())

	for i := 1; i < len(weighed); i++ {
		if weighed[i].Weight > weighed[i-1].Weight {
			t.Errorf("results not sorted: index %d weight %.2f > index %d weight %.2f",
				i, weighed[i].Weight, i-1, weighed[i-1].Weight)
		}
	}
}

func TestWeigher_SensitiveAttributeBoostsWeight(t *testing.T) {
	w := NewWeigher(DefaultWeighOptions())
	weighed := w.Weigh(makeWeigherResults())

	// aws_iam_user has a "password" change — its weight must exceed aws_s3_bucket
	// which only has a plain attribute change.
	var iamWeight, s3Weight float64
	for _, wr := range weighed {
		switch wr.Result.ResourceType {
		case "aws_iam_user":
			iamWeight = wr.Weight
		case "aws_s3_bucket":
			s3Weight = wr.Weight
		}
	}
	if iamWeight <= s3Weight {
		t.Errorf("expected iam weight %.2f > s3 weight %.2f", iamWeight, s3Weight)
	}
}

func TestWeigher_MissingKindAppliesMultiplier(t *testing.T) {
	opts := DefaultWeighOptions()
	w := NewWeigher(opts)

	results := []drift.DriftResult{
		{ResourceType: "aws_security_group", ResourceID: "sg-001",
			Changes: []drift.Change{{Attribute: "port", Kind: drift.Missing}}},
		{ResourceType: "aws_s3_bucket", ResourceID: "bucket",
			Changes: []drift.Change{{Attribute: "acl", Kind: drift.Changed}}},
	}
	weighed := w.Weigh(results)

	if weighed[0].Weight <= weighed[1].Weight {
		t.Errorf("missing-kind result should outweigh changed-kind result")
	}
}

func TestWeigherPrinter_NoDrift_PrintsMessage(t *testing.T) {
	var buf bytes.Buffer
	p := NewWeigherPrinter(&buf)
	p.Print(nil)
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected 'No drift' message, got: %s", buf.String())
	}
}

func TestWeigherPrinter_ContainsResourceType(t *testing.T) {
	var buf bytes.Buffer
	p := NewWeigherPrinter(&buf)
	w := NewWeigher(DefaultWeighOptions())
	p.Print(w.Weigh(makeWeigherResults()))
	output := buf.String()
	if !strings.Contains(output, "aws_iam_user") {
		t.Errorf("expected resource type in output, got: %s", output)
	}
}

func TestWeigherPrinter_NilWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic
	p := NewWeigherPrinter(nil)
	if p.w == nil {
		t.Error("expected non-nil writer")
	}
}
