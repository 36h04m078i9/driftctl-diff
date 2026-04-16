// Package baseline provides functionality to save and load a known-good
// drift baseline so that acknowledged drift can be suppressed in future runs.
package baseline

import (
	"encoding/json"
	"os"
	"time"

	"github.com/snyk/driftctl-diff/internal/drift"
)

// Entry records a single acknowledged drift item.
type Entry struct {
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	Attribute    string `json:"attribute"`
	Acknowledged time.Time `json:"acknowledged"`
}

// Baseline holds a set of acknowledged drift entries.
type Baseline struct {
	Entries []Entry `json:"entries"`
}

// New returns an empty Baseline.
func New() *Baseline {
	return &Baseline{}
}

// Add appends an entry derived from a DriftResult change.
func (b *Baseline) Add(resourceType, resourceID, attribute string) {
	b.Entries = append(b.Entries, Entry{
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Attribute:    attribute,
		Acknowledged: time.Now().UTC(),
	})
}

// Contains reports whether the given change is already acknowledged.
func (b *Baseline) Contains(resourceType, resourceID, attribute string) bool {
	for _, e := range b.Entries {
		if e.ResourceType == resourceType && e.ResourceID == resourceID && e.Attribute == attribute {
			return true
		}
	}
	return false
}

// Filter returns only the changes not covered by this baseline.
func (b *Baseline) Filter(results []drift.ResourceDiff) []drift.ResourceDiff {
	var out []drift.ResourceDiff
	for _, r := range results {
		var changes []drift.AttributeChange
		for _, c := range r.Changes {
			if !b.Contains(r.ResourceType, r.ResourceID, c.Attribute) {
				changes = append(changes, c)
			}
		}
		if len(changes) > 0 {
			out = append(out, drift.ResourceDiff{
				ResourceType: r.ResourceType,
				ResourceID:   r.ResourceID,
				Changes:      changes,
			})
		}
	}
	return out
}

// SaveTo writes the baseline to a JSON file at path.
func (b *Baseline) SaveTo(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(b)
}

// LoadFrom reads a baseline from a JSON file at path.
func LoadFrom(path string) (*Baseline, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, err
	}
	return &b, nil
}
