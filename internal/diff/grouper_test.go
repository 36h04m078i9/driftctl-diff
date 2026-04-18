package diff_test

import (
	"testing"

	"github.com/driftctl/driftctl-diff/internal/diff"
	"github.com/driftctl/driftctl-diff/internal/drift"
)

func makeGroupResults() []drift.Result {
	return []drift.Result{
		{ResourceType: "aws_s3_bucket", ResourceID: "bucket-1", Changes: []drift.Change{{Attribute: "acl", Kind: drift.KindChanged}}} ,
		{ResourceType: "aws_instance", ResourceID: "inst-1", Changes: []drift.Change{{Attribute: "ami", Kind: drift.KindChanged}}},
		{ResourceType: "aws_s3_bucket", ResourceID: "bucket-2", Changes: []drift.Change{{Attribute: "region", Kind: drift.KindMissing}}},
		{ResourceType: "aws_instance", ResourceID: "inst-2", Changes: []drift.Change{}},
	}
}

func TestGroupByType_CorrectNumberOfGroups(t *testing.T) {
	g := diff.NewGrouper()
	groups := g.GroupByType(makeGroupResults())
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}

func TestGroupByType_PreservesOrder(t *testing.T) {
	g := diff.NewGrouper()
	groups := g.GroupByType(makeGroupResults())
	if groups[0].ResourceType != "aws_s3_bucket" {
		t.Errorf("expected first group aws_s3_bucket, got %s", groups[0].ResourceType)
	}
	if len(groups[0].Results) != 2 {
		t.Errorf("expected 2 results in s3 group, got %d", len(groups[0].Results))
	}
}

func TestGroupByType_Empty(t *testing.T) {
	g := diff.NewGrouper()
	groups := g.GroupByType(nil)
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(groups))
	}
}

func TestGroupByKind_SeparatesKinds(t *testing.T) {
	g := diff.NewGrouper()
	out := g.GroupByKind(makeGroupResults())
	if len(out[string(drift.KindChanged)]) != 2 {
		t.Errorf("expected 2 changed, got %d", len(out[string(drift.KindChanged)]))
	}
	if len(out[string(drift.KindMissing)]) != 1 {
		t.Errorf("expected 1 missing, got %d", len(out[string(drift.KindMissing)]))
	}
	if len(out["none"]) != 1 {
		t.Errorf("expected 1 none, got %d", len(out["none"]))
	}
}

func TestGroupByKind_Empty(t *testing.T) {
	g := diff.NewGrouper()
	out := g.GroupByKind([]drift.Result{})
	if len(out) != 0 {
		t.Errorf("expected empty map")
	}
}
