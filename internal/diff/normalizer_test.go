package diff_test

import (
	"testing"

	"github.com/acme/driftctl-diff/internal/diff"
	"github.com/acme/driftctl-diff/internal/drift"
)

func makeNormResults(changes ...drift.AttributeChange) []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceID:   "res-1",
			ResourceType: "aws_instance",
			Changes:      changes,
		},
	}
}

func TestNormalizer_TrimSpace_RemovesCosmetic(t *testing.T) {
	n := diff.NewNormalizer(diff.NormalizeOptions{TrimSpace: true})
	input := makeNormResults(drift.AttributeChange{
		Attribute: "tags",
		Got:       "  production  ",
		Want:      "production",
		Kind:      drift.KindChanged,
	})
	out := n.Normalize(input)
	if len(out) != 0 {
		t.Fatalf("expected cosmetic diff to be removed, got %d results", len(out))
	}
}

func TestNormalizer_TrimSpace_KeepsRealDiff(t *testing.T) {
	n := diff.NewNormalizer(diff.NormalizeOptions{TrimSpace: true})
	input := makeNormResults(drift.AttributeChange{
		Attribute: "ami",
		Got:       "ami-111",
		Want:      "ami-222",
		Kind:      drift.KindChanged,
	})
	out := n.Normalize(input)
	if len(out) != 1 || len(out[0].Changes) != 1 {
		t.Fatal("expected real diff to be preserved")
	}
}

func TestNormalizer_LowerCase_Deduplicates(t *testing.T) {
	n := diff.NewNormalizer(diff.NormalizeOptions{LowerCase: true})
	input := makeNormResults(drift.AttributeChange{
		Attribute: "region",
		Got:       "US-EAST-1",
		Want:      "us-east-1",
		Kind:      drift.KindChanged,
	})
	out := n.Normalize(input)
	if len(out) != 0 {
		t.Fatalf("expected case-only diff to be removed, got %d results", len(out))
	}
}

func TestNormalizer_StripQuotes_RemovesQuotes(t *testing.T) {
	n := diff.NewNormalizer(diff.NormalizeOptions{StripQuotes: true})
	input := makeNormResults(drift.AttributeChange{
		Attribute: "name",
		Got:       `"my-bucket"`,
		Want:      "my-bucket",
		Kind:      drift.KindChanged,
	})
	out := n.Normalize(input)
	if len(out) != 0 {
		t.Fatalf("expected quote-only diff to be removed, got %d results", len(out))
	}
}

func TestNormalizer_EmptyInput_ReturnsEmpty(t *testing.T) {
	n := diff.NewNormalizer(diff.DefaultNormalizeOptions())
	out := n.Normalize(nil)
	if len(out) != 0 {
		t.Fatal("expected empty output for nil input")
	}
}

func TestNormalizer_AllChangesCosmetic_DropsResult(t *testing.T) {
	n := diff.NewNormalizer(diff.NormalizeOptions{TrimSpace: true})
	input := makeNormResults(
		drift.AttributeChange{Attribute: "a", Got: " x ", Want: "x", Kind: drift.KindChanged},
		drift.AttributeChange{Attribute: "b", Got: " y ", Want: "y", Kind: drift.KindChanged},
	)
	out := n.Normalize(input)
	if len(out) != 0 {
		t.Fatal("expected entire result to be dropped when all changes are cosmetic")
	}
}
