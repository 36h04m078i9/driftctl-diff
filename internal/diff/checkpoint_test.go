package diff

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/drift"
)

func makeCheckpointResults(n int) []drift.DriftResult {
	results := make([]drift.DriftResult, n)
	for i := range results {
		results[i] = drift.DriftResult{
			ResourceID:   "res-" + string(rune('a'+i)),
			ResourceType: "aws_instance",
		}
	}
	return results
}

func TestCheckpointStore_SaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	store := NewCheckpointStore(dir)
	results := makeCheckpointResults(3)

	if err := store.Save("snap1", results); err != nil {
		t.Fatalf("Save: %v", err)
	}
	cp, err := store.Load("snap1")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cp.Name != "snap1" {
		t.Errorf("name = %q, want snap1", cp.Name)
	}
	if len(cp.Results) != 3 {
		t.Errorf("results len = %d, want 3", len(cp.Results))
	}
}

func TestCheckpointStore_Load_MissingFile_ReturnsError(t *testing.T) {
	store := NewCheckpointStore(t.TempDir())
	_, err := store.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCheckpointStore_List_ReturnsNames(t *testing.T) {
	dir := t.TempDir()
	store := NewCheckpointStore(dir)
	_ = store.Save("alpha", makeCheckpointResults(1))
	_ = store.Save("beta", makeCheckpointResults(2))

	names, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("len = %d, want 2", len(names))
	}
}

func TestCheckpointStore_List_EmptyDir_ReturnsNil(t *testing.T) {
	store := NewCheckpointStore(t.TempDir())
	names, err := store.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty, got %v", names)
	}
}

func TestCheckpointStore_List_MissingDir_ReturnsNil(t *testing.T) {
	store := NewCheckpointStore(filepath.Join(os.TempDir(), "no-such-dir-xyz"))
	names, err := store.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if names != nil {
		t.Errorf("expected nil, got %v", names)
	}
}

func TestCheckpointPrinter_NoDrift_PrintsNoCheckpointsMessage(t *testing.T) {
	store := NewCheckpointStore(t.TempDir())
	var buf bytes.Buffer
	p := NewCheckpointPrinter(store, &buf)
	if err := p.Print(); err != nil {
		t.Fatalf("Print: %v", err)
	}
	if !strings.Contains(buf.String(), "no checkpoints") {
		t.Errorf("output = %q, want 'no checkpoints'...", buf.String())
	}
}

func TestCheckpointPrinter_WithCheckpoints_ContainsName(t *testing.T) {
	dir := t.TempDir()
	store := NewCheckpointStore(dir)
	_ = store.Save("mysnap", makeCheckpointResults(2))
	var buf bytes.Buffer
	p := NewCheckpointPrinter(store, &buf)
	if err := p.Print(); err != nil {
		t.Fatalf("Print: %v", err)
	}
	if !strings.Contains(buf.String(), "mysnap") {
		t.Errorf("output = %q, want 'mysnap'", buf.String())
	}
}

func TestCheckpointPrinter_NilWriter_DefaultsToStdout(t *testing.T) {
	store := NewCheckpointStore(t.TempDir())
	p := NewCheckpointPrinter(store, nil)
	if p.w == nil {
		t.Error("expected non-nil writer")
	}
}
