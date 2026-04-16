package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/acme/driftctl-diff/internal/drift"
)

// JSONFormatter renders drift results as a JSON document.
type JSONFormatter struct{}

// NewJSONFormatter creates a new JSONFormatter.
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

type jsonResource struct {
	ID      string            `json:"id"`
	Type    string            `json:"type"`
	Changes []jsonChange      `json:"changes"`
}

type jsonChange struct {
	Attribute string `json:"attribute"`
	Kind      string `json:"kind"`
	Want      string `json:"want,omitempty"`
	Got       string `json:"got,omitempty"`
}

type jsonOutput struct {
	Drifted   bool           `json:"drifted"`
	Resources []jsonResource `json:"resources"`
}

// Format writes drift results as JSON to w.
func (f *JSONFormatter) Format(w io.Writer, results []drift.ResourceDiff) error {
	out := jsonOutput{
		Drifted:   len(results) > 0,
		Resources: make([]jsonResource, 0, len(results)),
	}

	for _, r := range results {
		jr := jsonResource{
			ID:   r.ResourceID,
			Type: r.ResourceType,
		}
		for _, c := range r.Changes {
			jr.Changes = append(jr.Changes, jsonChange{
				Attribute: c.Attribute,
				Kind:      formatKind(c.Kind),
				Want:      c.Want,
				Got:       c.Got,
			})
		}
		out.Resources = append(out.Resources, jr)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("json formatter: %w", err)
	}
	return nil
}
