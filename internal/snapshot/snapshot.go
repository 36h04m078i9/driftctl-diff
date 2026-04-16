// Package snapshot captures and persists live cloud resource attributes
// so they can be compared against Terraform state across runs.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot holds a point-in-time capture of live resource attributes.
type Snapshot struct {
	CapturedAt time.Time                       `json:"captured_at"`
	Resources  map[string]map[string]string    `json:"resources"` // resourceID -> attrs
}

// New returns an empty Snapshot stamped with the current time.
func New() *Snapshot {
	return &Snapshot{
		CapturedAt: time.Now().UTC(),
		Resources:  make(map[string]map[string]string),
	}
}

// Add stores attributes for a resource ID.
func (s *Snapshot) Add(resourceID string, attrs map[string]string) {
	s.Resources[resourceID] = attrs
}

// Get returns the attributes for a resource ID and whether it was found.
func (s *Snapshot) Get(resourceID string) (map[string]string, bool) {
	attrs, ok := s.Resources[resourceID]
	return attrs, ok
}

// SaveTo serialises the snapshot as JSON to the given file path.
func (s *Snapshot) SaveTo(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// LoadFrom deserialises a snapshot from the given file path.
func LoadFrom(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}
