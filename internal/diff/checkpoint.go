package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/owner/driftctl-diff/internal/drift"
)

// Checkpoint records a named snapshot of drift results for later comparison.
type Checkpoint struct {
	Name      string              `json:"name"`
	CreatedAt time.Time           `json:"created_at"`
	Results   []drift.DriftResult `json:"results"`
}

// CheckpointStore persists and retrieves named checkpoints on disk.
type CheckpointStore struct {
	dir string
}

// NewCheckpointStore returns a CheckpointStore rooted at dir.
func NewCheckpointStore(dir string) *CheckpointStore {
	return &CheckpointStore{dir: dir}
}

// Save writes results as a named checkpoint to disk.
func (s *CheckpointStore) Save(name string, results []drift.DriftResult) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return fmt.Errorf("checkpoint: mkdir %s: %w", s.dir, err)
	}
	cp := Checkpoint{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Results:   results,
	}
	data, err := json.Marshal(cp)
	if err != nil {
		return fmt.Errorf("checkpoint: marshal: %w", err)
	}
	path := filepath.Join(s.dir, name+".json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("checkpoint: write %s: %w", path, err)
	}
	return nil
}

// Load reads a named checkpoint from disk.
func (s *CheckpointStore) Load(name string) (*Checkpoint, error) {
	path := filepath.Join(s.dir, name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("checkpoint: read %s: %w", path, err)
	}
	var cp Checkpoint
	if err := json.Unmarshal(data, &cp); err != nil {
		return nil, fmt.Errorf("checkpoint: unmarshal: %w", err)
	}
	return &cp, nil
}

// List returns the names of all stored checkpoints.
func (s *CheckpointStore) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("checkpoint: readdir: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}
