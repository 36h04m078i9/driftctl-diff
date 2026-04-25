package diff_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/diff"
	"github.com/acme/driftctl-diff/internal/drift"
)

func makeArchiverResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceID:   "aws_instance.web",
			ResourceType: "aws_instance",
			Changes: []drift.AttributeChange{
				{Attribute: "instance_type", StateValue: "t2.micro", LiveValue: "t3.micro"},
			},
		},
	}
}

func TestArchiver_SaveCreatesFile(t *testing.T) {
	dir := t.TempDir()
	a := diff.NewArchiver(dir)
	path, err := a.Save(makeArchiverResults(), "test-label")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist at %s", path)
	}
}

func TestArchiver_LoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	a := diff.NewArchiver(dir)
	original := makeArchiverResults()
	path, err := a.Save(original, "round-trip")
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	entry, err := a.Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if entry.Label != "round-trip" {
		t.Errorf("expected label 'round-trip', got %q", entry.Label)
	}
	if len(entry.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(entry.Results))
	}
	if entry.Results[0].ResourceID != "aws_instance.web" {
		t.Errorf("unexpected resource ID: %s", entry.Results[0].ResourceID)
	}
}

func TestArchiver_List_ReturnsAllFiles(t *testing.T) {
	dir := t.TempDir()
	a := diff.NewArchiver(dir)
	for i := 0; i < 3; i++ {
		if _, err := a.Save(makeArchiverResults(), ""); err != nil {
			t.Fatalf("save: %v", err)
		}
	}
	paths, err := a.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(paths) != 3 {
		t.Errorf("expected 3 archives, got %d", len(paths))
	}
}

func TestArchiver_Load_MissingFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	a := diff.NewArchiver(dir)
	_, err := a.Load(filepath.Join(dir, "nonexistent.json"))
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestArchiverPrinter_EmptyDir_PrintsNoArchivesMessage(t *testing.T) {
	dir := t.TempDir()
	a := diff.NewArchiver(dir)
	var buf strings.Builder
	p := diff.NewArchiverPrinter(a, &buf)
	if err := p.Print(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No archives found") {
		t.Errorf("expected 'No archives found' message, got: %s", buf.String())
	}
}

func TestArchiverPrinter_WithEntries_ContainsLabel(t *testing.T) {
	dir := t.TempDir()
	a := diff.NewArchiver(dir)
	if _, err := a.Save(makeArchiverResults(), "my-label"); err != nil {
		t.Fatalf("save: %v", err)
	}
	var buf strings.Builder
	p := diff.NewArchiverPrinter(a, &buf)
	if err := p.Print(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "my-label") {
		t.Errorf("expected label in output, got: %s", buf.String())
	}
}
