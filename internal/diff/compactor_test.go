package diff_test

import (
	"testing"

	"github.com/acme/driftctl-diff/internal/diff"
	"github.com/acme/driftctl-diff/internal/drift"
)

func makeCompactResults(resources []struct {
	id      string
	rtype   string
	changes []drift.Change
}) []drift.DriftResult {
	results := make([]drift.DriftResult, 0, len(resources))
	for _, r := range resources {
		results = append(results, drift.DriftResult{
			ResourceID:   r.id,
			ResourceType: r.rtype,
			Changes:      r.changes,
		})
	}
	return results
}

func TestCompactor_NoOptions_ReturnsAll(t *testing.T) {
	input := makeCompactResults([]struct {
		id      string
		rtype   string
		changes []drift.Change
	}{
		{"res-1", "aws_s3_bucket", []drift.Change{{Attribute: "acl", Got: "public", Want: "private"}}},
		{"res-2", "aws_instance", []drift.Change{{Attribute: "ami", Got: "ami-new", Want: "ami-old"}}},
	})

	c := diff.NewCompactor(diff.DefaultCompactOptions())
	out := c.Compact(input)

	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestCompactor_MergesAdjacentSameResource(t *testing.T) {
	changeA := drift.Change{Attribute: "acl", Got: "public", Want: "private", Kind: drift.KindChanged}
	changeB := drift.Change{Attribute: "tags", Got: "a", Want: "b", Kind: drift.KindChanged}

	input := []drift.DriftResult{
		{ResourceID: "res-1", ResourceType: "aws_s3_bucket", Changes: []drift.Change{changeA}},
		{ResourceID: "res-1", ResourceType: "aws_s3_bucket", Changes: []drift.Change{changeB}},
	}

	opts := diff.DefaultCompactOptions()
	opts.MergeAdjacentResources = true
	c := diff.NewCompactor(opts)
	out := c.Compact(input)

	if len(out) != 1 {
		t.Fatalf("expected 1 merged result, got %d", len(out))
	}
	if len(out[0].Changes) != 2 {
		t.Fatalf("expected 2 changes after merge, got %d", len(out[0].Changes))
	}
}

func TestCompactor_DeduplicatesChanges(t *testing.T) {
	change := drift.Change{Attribute: "acl", Got: "public", Want: "private", Kind: drift.KindChanged}

	input := []drift.DriftResult{
		{ResourceID: "res-1", ResourceType: "aws_s3_bucket", Changes: []drift.Change{change, change}},
	}

	opts := diff.DefaultCompactOptions()
	opts.DeduplicateChanges = true
	c := diff.NewCompactor(opts)
	out := c.Compact(input)

	if len(out[0].Changes) != 1 {
		t.Fatalf("expected 1 deduplicated change, got %d", len(out[0].Changes))
	}
}

func TestCompactor_EmptyInput_ReturnsEmpty(t *testing.T) {
	c := diff.NewCompactor(diff.DefaultCompactOptions())
	out := c.Compact(nil)

	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d", len(out))
	}
}

func TestCompactor_DropEmptyResults_RemovesNoChangeEntries(t *testing.T) {
	input := []drift.DriftResult{
		{ResourceID: "res-1", ResourceType: "aws_s3_bucket", Changes: []drift.Change{}},
		{ResourceID: "res-2", ResourceType: "aws_instance", Changes: []drift.Change{
			{Attribute: "ami", Got: "new", Want: "old", Kind: drift.KindChanged},
		}},
	}

	opts := diff.DefaultCompactOptions()
	opts.DropEmptyResults = true
	c := diff.NewCompactor(opts)
	out := c.Compact(input)

	if len(out) != 1 {
		t.Fatalf("expected 1 result after dropping empty, got %d", len(out))
	}
	if out[0].ResourceID != "res-2" {
		t.Errorf("expected res-2, got %s", out[0].ResourceID)
	}
}
