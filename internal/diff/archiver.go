package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/acme/driftctl-diff/internal/drift"
)

// ArchiveEntry holds a timestamped snapshot of drift results.
type ArchiveEntry struct {
	ArchivedAt time.Time           `json:"archived_at"`
	Label      string              `json:"label,omitempty"`
	Results    []drift.DriftResult `json:"results"`
}

// Archiver persists and retrieves drift result archives on disk.
type Archiver struct {
	dir string
}

// NewArchiver creates an Archiver that stores archives under dir.
func NewArchiver(dir string) *Archiver {
	return &Archiver{dir: dir}
}

// Save writes results to a JSON file named by timestamp and optional label.
func (a *Archiver) Save(results []drift.DriftResult, label string) (string, error) {
	if err := os.MkdirAll(a.dir, 0o755); err != nil {
		return "", fmt.Errorf("archiver: create dir: %w", err)
	}
	entry := ArchiveEntry{
		ArchivedAt: time.Now().UTC(),
		Label:      label,
		Results:    results,
	}
	filename := fmt.Sprintf("%d.json", entry.ArchivedAt.UnixNano())
	path := filepath.Join(a.dir, filename)
	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("archiver: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entry); err != nil {
		return "", fmt.Errorf("archiver: encode: %w", err)
	}
	return path, nil
}

// Load reads an archive entry from the given file path.
func (a *Archiver) Load(path string) (ArchiveEntry, error) {
	var entry ArchiveEntry
	f, err := os.Open(path)
	if err != nil {
		return entry, fmt.Errorf("archiver: open: %w", err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&entry); err != nil {
		return entry, fmt.Errorf("archiver: decode: %w", err)
	}
	return entry, nil
}

// List returns all archive file paths in the directory, sorted by name.
func (a *Archiver) List() ([]string, error) {
	entries, err := filepath.Glob(filepath.Join(a.dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("archiver: list: %w", err)
	}
	return entries, nil
}
