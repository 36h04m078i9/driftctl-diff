package output_test

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"

	"github.com/snyk/driftctl-diff/internal/drift"
	"github.com/snyk/driftctl-diff/internal/output"
)

func TestCSVFormatter_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewCSVFormatter(&buf)
	if err := f.Format(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := csv.NewReader(&buf)
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("invalid csv: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected header only, got %d rows", len(records))
	}
}

func TestCSVFormatter_WithChanges_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewCSVFormatter(&buf)
	_ = f.Format(sampleCSVChanges())
	line := buf.String()
	for _, h := range []string{"resource_type", "resource_id", "attribute", "kind"} {
		if !strings.Contains(line, h) {
			t.Errorf("missing header %q", h)
		}
	}
}

func TestCSVFormatter_WithChanges_ContainsValues(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewCSVFormatter(&buf)
	_ = f.Format(sampleCSVChanges())

	r := csv.NewReader(&buf)
	records, _ := r.ReadAll()
	// header + 1 change row
	if len(records) < 2 {
		t.Fatalf("expected at least 2 rows, got %d", len(records))
	}
	row := records[1]
	if row[0] != "aws_s3_bucket" {
		t.Errorf("expected aws_s3_bucket, got %s", row[0])
	}
	if row[1] != "my-bucket" {
		t.Errorf("expected my-bucket, got %s", row[1])
	}
}

func TestCSVFormatter_NilWriter_DefaultsToStdout(t *testing.T) {
	// Should not panic
	f := output.NewCSVFormatter(nil)
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}

func sampleCSVChanges() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "acl", Kind: drift.KindChanged, StateValue: "private", LiveValue: "public-read"},
			},
		},
	}
}
