package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/diff"
	"github.com/owner/driftctl-diff/internal/drift"
)

func makeClassifierResults() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceType: "aws_iam_role",
			ResourceID:   "role-1",
			Changes: []drift.Change{
				{Attribute: "password", StateValue: "x", LiveValue: "y"},
			},
		},
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "bucket-1",
			Changes:      []drift.Change{},
		},
		{
			ResourceType: "aws_security_group",
			ResourceID:   "sg-1",
			Changes: []drift.Change{
				{Attribute: "a"}, {Attribute: "b"}, {Attribute: "c"},
				{Attribute: "d"}, {Attribute: "e"},
			},
		},
	}
}

func TestClassifier_HighSensitiveAttribute(t *testing.T) {
	c := diff.NewClassifier(diff.ClassifierOptions{
		HighValueAttrs: []string{"password"},
	})
	results := c.Classify(makeClassifierResults())
	if results[0].Severity != diff.SeverityHigh {
		t.Errorf("expected high, got %s", results[0].Severity)
	}
}

func TestClassifier_LowWhenNoChanges(t *testing.T) {
	c := diff.NewClassifier(diff.ClassifierOptions{})
	results := c.Classify(makeClassifierResults())
	if results[1].Severity != diff.SeverityLow {
		t.Errorf("expected low, got %s", results[1].Severity)
	}
}

func TestClassifier_MediumWhenManyChanges(t *testing.T) {
	c := diff.NewClassifier(diff.ClassifierOptions{})
	results := c.Classify(makeClassifierResults())
	if results[2].Severity != diff.SeverityMedium {
		t.Errorf("expected medium, got %s", results[2].Severity)
	}
}

func TestClassifier_CriticalResourceType(t *testing.T) {
	c := diff.NewClassifier(diff.ClassifierOptions{
		CriticalTypes: []string{"aws_iam_role"},
	})
	results := c.Classify(makeClassifierResults())
	if results[0].Severity != diff.SeverityCritical {
		t.Errorf("expected critical, got %s", results[0].Severity)
	}
}

func TestClassifierPrinter_NoDrift_PrintsMessage(t *testing.T) {
	var buf bytes.Buffer
	p := diff.NewClassifierPrinter(&buf)
	p.Print(nil)
	if !strings.Contains(buf.String(), "No classifications") {
		t.Error("expected no classifications message")
	}
}

func TestClassifierPrinter_ContainsResourceID(t *testing.T) {
	c := diff.NewClassifier(diff.ClassifierOptions{})
	classifications := c.Classify(makeClassifierResults())
	var buf bytes.Buffer
	p := diff.NewClassifierPrinter(&buf)
	p.Print(classifications)
	if !strings.Contains(buf.String(), "role-1") {
		t.Error("expected resource id in output")
	}
}
