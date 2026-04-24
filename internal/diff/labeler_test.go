package diff_test

import (
	"fmt"
	"testing"

	"github.com/acme/driftctl-diff/internal/diff"
	"github.com/acme/driftctl-diff/internal/drift"
)

func makeLabelResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceID:   "aws_s3_bucket.logs",
			ResourceType: "aws_s3_bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "acl", StateValue: "private", LiveValue: "public-read", Kind: drift.KindChanged},
				{Attribute: "versioning", StateValue: "true", LiveValue: "", Kind: drift.KindMissing},
			},
		},
		{
			ResourceID:   "aws_iam_role.exec",
			ResourceType: "aws_iam_role",
			Changes:      []drift.AttributeChange{},
		},
	}
}

func TestLabeler_AttachesResourceTypeLabel(t *testing.T) {
	l := diff.NewLabeler()
	results := makeLabelResults()
	labeled := l.Label(results)

	if len(labeled) != 2 {
		t.Fatalf("expected 2 labeled results, got %d", len(labeled))
	}
	for _, lr := range labeled {
		found := false
		for _, lbl := range lr.Labels {
			if lbl.Key == "resource_type" {
				found = true
				if lbl.Value != lr.ResourceType {
					t.Errorf("expected resource_type %q, got %q", lr.ResourceType, lbl.Value)
				}
			}
		}
		if !found {
			t.Errorf("resource_type label missing for %s", lr.ResourceID)
		}
	}
}

func TestLabeler_ChangeCountLabel(t *testing.T) {
	l := diff.NewLabeler()
	labeled := l.Label(makeLabelResults())

	for _, lr := range labeled {
		for _, lbl := range lr.Labels {
			if lbl.Key == "change_count" {
				expected := fmt.Sprintf("%d", len(lr.Changes))
				if lbl.Value != expected {
					t.Errorf("expected change_count %s, got %s", expected, lbl.Value)
				}
			}
		}
	}
}

func TestLabeler_DriftKindLabel_PresentWhenChanges(t *testing.T) {
	l := diff.NewLabeler()
	labeled := l.Label(makeLabelResults())

	first := labeled[0]
	for _, lbl := range first.Labels {
		if lbl.Key == "drift_kind" {
			if lbl.Value == "" {
				t.Error("expected non-empty drift_kind label")
			}
			return
		}
	}
	t.Error("drift_kind label not found")
}

func TestLabeler_DriftKindLabel_AbsentWhenNoChanges(t *testing.T) {
	l := diff.NewLabeler()
	labeled := l.Label(makeLabelResults())

	second := labeled[1]
	for _, lbl := range second.Labels {
		if lbl.Key == "drift_kind" {
			t.Errorf("expected no drift_kind label for result with no changes, got %q", lbl.Value)
		}
	}
}

func TestLabeler_EmptyInput_ReturnsEmpty(t *testing.T) {
	l := diff.NewLabeler()
	labeled := l.Label([]drift.DriftResult{})
	if len(labeled) != 0 {
		t.Errorf("expected empty slice, got %d", len(labeled))
	}
}

// findLabel is a helper that returns the value of a label by key, and whether it was found.
func findLabel(labels []diff.Label, key string) (string, bool) {
	for _, lbl := range labels {
		if lbl.Key == key {
			return lbl.Value, true
		}
	}
	return "", false
}
