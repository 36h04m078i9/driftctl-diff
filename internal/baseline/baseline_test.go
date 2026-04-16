package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/snyk/driftctl-diff/internal/baseline"
	"github.com/snyk/driftctl-diff/internal/drift"
)

func sampleDiffs() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-123",
			Changes: []drift.AttributeChange{
				{Attribute: "instance_type", StateValue: "t2.micro", LiveValue: "t3.micro"},
				{Attribute: "tags", StateValue: "a", LiveValue: "b"},
			},
		},
	}
}

func TestNew_IsEmpty(t *testing.T) {
	b := baseline.New()
	if len(b.Entries) != 0 {
		t.Fatalf("expected empty baseline, got %d entries", len(b.Entries))
	}
}

func TestAdd_AndContains(t *testing.T) {
	b := baseline.New()
	b.Add("aws_instance", "i-123", "instance_type")
	if !b.Contains("aws_instance", "i-123", "instance_type") {
		t.Fatal("expected Contains to return true")
	}
	if b.Contains("aws_instance", "i-123", "tags") {
		t.Fatal("expected Contains to return false for un-added attribute")
	}
}

func TestFilter_RemovesAcknowledged(t *testing.T) {
	b := baseline.New()
	b.Add("aws_instance", "i-123", "instance_type")

	results := b.Filter(sampleDiffs())
	if len(results) != 1 {
		t.Fatalf("expected 1 resource diff, got %d", len(results))
	}
	if len(results[0].Changes) != 1 {
		t.Fatalf("expected 1 remaining change, got %d", len(results[0].Changes))
	}
	if results[0].Changes[0].Attribute != "tags" {
		t.Fatalf("unexpected attribute %s", results[0].Changes[0].Attribute)
	}
}

func TestFilter_AllAcknowledged_EmptyResult(t *testing.T) {
	b := baseline.New()
	b.Add("aws_instance", "i-123", "instance_type")
	b.Add("aws_instance", "i-123", "tags")

	results := b.Filter(sampleDiffs())
	if len(results) != 0 {
		t.Fatalf("expected empty results, got %d", len(results))
	}
}

func TestSaveTo_AndLoadFrom_RoundTrip(t *testing.T) {
	b := baseline.New()
	b.Add("aws_s3_bucket", "my-bucket", "acl")

	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	if err := b.SaveTo(path); err != nil {
		t.Fatalf("SaveTo error: %v", err)
	}

	loaded, err := baseline.LoadFrom(path)
	if err != nil {
		t.Fatalf("LoadFrom error: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].ResourceID != "my-bucket" {
		t.Fatalf("unexpected resource id %s", loaded.Entries[0].ResourceID)
	}
}

func TestLoadFrom_MissingFile(t *testing.T) {
	_, err := baseline.LoadFrom("/nonexistent/baseline.json")
	if !os.IsNotExist(err) {
		t.Fatalf("expected not-exist error, got %v", err)
	}
}
