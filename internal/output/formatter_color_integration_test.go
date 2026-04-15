package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/driftctl-diff/internal/drift"
	"github.com/your-org/driftctl-diff/internal/output"
)

// sampleChanges returns a small set of drift changes for formatting tests.
func sampleChanges() []drift.Change {
	return []drift.Change{
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-0abc123",
			Attribute:    "instance_type",
			StateValue:   "t2.micro",
			LiveValue:    "t3.medium",
		},
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Attribute:    "tags.Env",
			StateValue:   "production",
			LiveValue:    "",
		},
	}
}

func TestFormatter_WithColorDisabled_NoEscapeCodes(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, false)
	f.Write(sampleChanges())

	result := buf.String()
	if strings.Contains(result, "\033[") {
		t.Errorf("expected no ANSI escape codes when color is disabled, got:\n%s", result)
	}
}

func TestFormatter_WithColorEnabled_ContainsEscapeCodes(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, true)
	f.Write(sampleChanges())

	result := buf.String()
	if !strings.Contains(result, "\033[") {
		t.Errorf("expected ANSI escape codes when color is enabled, got:\n%s", result)
	}
}

func TestFormatter_OutputContainsResourceInfo(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, false)
	f.Write(sampleChanges())

	result := buf.String()
	for _, want := range []string{"aws_instance", "i-0abc123", "instance_type", "t2.micro", "t3.medium"} {
		if !strings.Contains(result, want) {
			t.Errorf("expected output to contain %q\nfull output:\n%s", want, result)
		}
	}
}
