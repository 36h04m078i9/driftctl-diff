package drift

import (
	"fmt"

	"github.com/user/driftctl-diff/internal/state"
)

// ChangeType classifies the kind of drift detected.
type ChangeType string

const (
	ChangeAdded   ChangeType = "added"   // exists in cloud, missing from state
	ChangeDeleted ChangeType = "deleted" // exists in state, missing from cloud
	ChangeModified ChangeType = "modified" // exists in both but attributes differ
)

// Change describes a single drifted resource.
type Change struct {
	ResourceType string
	ResourceName string
	ChangeType   ChangeType
	Diff         map[string]DiffEntry
}

// DiffEntry holds the state value vs the live value for a single attribute.
type DiffEntry struct {
	StateValue interface{}
	LiveValue  interface{}
}

// Results aggregates all detected drift changes.
type Results struct {
	Changes []Change
}

// HasDrift returns true when at least one change was detected.
func (r *Results) HasDrift() bool {
	return len(r.Changes) > 0
}

// Detector compares Terraform state against live cloud resources.
type Detector struct {
	region string
}

// NewDetector creates a Detector targeting the given AWS region.
func NewDetector(region string) *Detector {
	return &Detector{region: region}
}

// Detect runs drift detection for all resources in tfState.
// NOTE: Live cloud lookup is stubbed; replace with real AWS SDK calls.
func (d *Detector) Detect(tfState *state.State) (*Results, error) {
	results := &Results{}
	for _, res := range tfState.Resources {
		live, err := d.fetchLive(res)
		if err != nil {
			return nil, fmt.Errorf("fetching live resource %s.%s: %w", res.Type, res.Name, err)
		}
		if live == nil {
			results.Changes = append(results.Changes, Change{
				ResourceType: res.Type,
				ResourceName: res.Name,
				ChangeType:   ChangeDeleted,
			})
			continue
		}
		if diff := diffAttributes(res.Attributes, live); len(diff) > 0 {
			results.Changes = append(results.Changes, Change{
				ResourceType: res.Type,
				ResourceName: res.Name,
				ChangeType:   ChangeModified,
				Diff:         diff,
			})
		}
	}
	return results, nil
}

// fetchLive is a stub that should be replaced with real AWS SDK lookups.
func (d *Detector) fetchLive(res state.Resource) (map[string]interface{}, error) {
	// Returning the same attributes simulates no drift for now.
	return res.Attributes, nil
}

func diffAttributes(stateAttrs, liveAttrs map[string]interface{}) map[string]DiffEntry {
	diff := make(map[string]DiffEntry)
	for k, sv := range stateAttrs {
		lv, ok := liveAttrs[k]
		if !ok || fmt.Sprintf("%v", sv) != fmt.Sprintf("%v", lv) {
			diff[k] = DiffEntry{StateValue: sv, LiveValue: lv}
		}
	}
	return diff
}
