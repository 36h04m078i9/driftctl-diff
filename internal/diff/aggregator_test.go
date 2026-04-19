package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/diff"
	"github.com/owner/driftctl-diff/internal/drift"
)

func makeAggResults() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceID:   "bucket-1",
			ResourceType: "aws_s3_bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "acl", Kind: drift.KindChanged, WantValue: "private", GotValue: "public"},
				{Attribute: "tags", Kind: drift.KindChanged, WantValue: "env=prod", GotValue: ""},
			},
		},
		{
			ResourceID:   "sg-1",
			ResourceType: "aws_security_group",
			Changes: []drift.AttributeChange{
				{Attribute: "ingress", Kind: drift.KindMissing, WantValue: "0.0.0.0/0", GotValue: ""},
			},
		},
		{
			ResourceID:   "bucket-2",
			ResourceType: "aws_s3_bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "versioning", Kind: drift.KindChanged, WantValue: "true", GotValue: "false"},
			},
		},
	}
}

func TestAggregator_NoDrift_PrintsNoDriftMessage(t *testing.T) {
	var buf bytes.Buffer
	a := diff.NewAggregator(&buf)
	a.Print(nil)
	if !strings.Contains(buf.String(), "no drift") {
		t.Errorf("expected no drift message, got: %s", buf.String())
	}
}

func TestAggregator_GroupsByTypeAndKind(t *testing.T) {
	a := diff.NewAggregator(nil)
	agg := a.Aggregate(makeAggResults())

	// aws_s3_bucket + KindChanged should have count 3 (acl, tags, versioning)
	var s3Changed *diff.AggregateResult
	for i := range agg {
		if agg[i].ResourceType == "aws_s3_bucket" && agg[i].Kind == drift.KindChanged {
			s3Changed = &agg[i]
		}
	}
	if s3Changed == nil {
		t.Fatal("expected aws_s3_bucket/changed bucket")
	}
	if s3Changed.Count != 3 {
		t.Errorf("expected count 3, got %d", s3Changed.Count)
	}
}

func TestAggregator_SortedByResourceType(t *testing.T) {
	a := diff.NewAggregator(nil)
	agg := a.Aggregate(makeAggResults())
	for i := 1; i < len(agg); i++ {
		if agg[i].ResourceType < agg[i-1].ResourceType {
			t.Errorf("results not sorted: %s before %s", agg[i-1].ResourceType, agg[i].ResourceType)
		}
	}
}

func TestAggregator_Print_ContainsResourceType(t *testing.T) {
	var buf bytes.Buffer
	a := diff.NewAggregator(&buf)
	a.Print(makeAggResults())
	if !strings.Contains(buf.String(), "aws_s3_bucket") {
		t.Errorf("expected resource type in output, got: %s", buf.String())
	}
}

func TestAggregator_NilWriter_DefaultsToStdout(t *testing.T) {
	a := diff.NewAggregator(nil)
	if a == nil {
		t.Fatal("expected non-nil aggregator")
	}
}
