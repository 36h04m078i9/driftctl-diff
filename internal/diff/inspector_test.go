package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeInspectorResults() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceID:   "aws_s3_bucket.logs",
			ResourceType: "aws_s3_bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "acl", StateValue: "private", LiveValue: "public-read", Kind: drift.KindChanged},
				{Attribute: "versioning", StateValue: "true", LiveValue: "false", Kind: drift.KindChanged},
			},
		},
		{
			ResourceID:   "aws_iam_role.worker",
			ResourceType: "aws_iam_role",
			Changes:      []drift.AttributeChange{},
		},
	}
}

func TestInspect_KnownResource_ReturnsResult(t *testing.T) {
	ins := NewInspector(nil)
	res, err := ins.Inspect(makeInspectorResults(), "aws_s3_bucket.logs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ResourceID != "aws_s3_bucket.logs" {
		t.Errorf("expected resource id aws_s3_bucket.logs, got %s", res.ResourceID)
	}
	if res.TotalChanges != 2 {
		t.Errorf("expected 2 changes, got %d", res.TotalChanges)
	}
}

func TestInspect_UnknownResource_ReturnsError(t *testing.T) {
	ins := NewInspector(nil)
	_, err := ins.Inspect(makeInspectorResults(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown resource")
	}
}

func TestInspect_LinesContainAttribute(t *testing.T) {
	ins := NewInspector(nil)
	res, _ := ins.Inspect(makeInspectorResults(), "aws_s3_bucket.logs")
	found := false
	for _, l := range res.Lines {
		if strings.Contains(l, "acl") {
			found = true
		}
	}
	if !found {
		t.Error("expected lines to contain attribute 'acl'")
	}
}

func TestPrint_WritesResourceID(t *testing.T) {
	var buf bytes.Buffer
	ins := NewInspector(&buf)
	res, _ := ins.Inspect(makeInspectorResults(), "aws_s3_bucket.logs")
	ins.Print(res)
	if !strings.Contains(buf.String(), "aws_s3_bucket.logs") {
		t.Error("expected output to contain resource id")
	}
}

func TestPrint_NilResult_PrintsMessage(t *testing.T) {
	var buf bytes.Buffer
	ins := NewInspector(&buf)
	ins.Print(nil)
	if !strings.Contains(buf.String(), "no inspection result") {
		t.Error("expected 'no inspection result' message")
	}
}

func TestPrint_NoChanges_PrintsNoChanges(t *testing.T) {
	var buf bytes.Buffer
	ins := NewInspector(&buf)
	res, _ := ins.Inspect(makeInspectorResults(), "aws_iam_role.worker")
	ins.Print(res)
	if !strings.Contains(buf.String(), "no changes") {
		t.Error("expected '(no changes)' in output")
	}
}
