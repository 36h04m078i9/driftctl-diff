package output

import (
	"fmt"
	"io"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/snyk/driftctl-diff/internal/drift"
)

// NewYAMLFormatter returns a Formatter that writes drift results as YAML.
func NewYAMLFormatter() Formatter {
	return &yamlFormatter{}
}

type yamlFormatter struct{}

type yamlOutput struct {
	GeneratedAt string          `yaml:"generated_at"`
	Drifted     bool            `yaml:"drifted"`
	TotalDrift  int             `yaml:"total_drift"`
	Changes     []yamlChange    `yaml:"changes,omitempty"`
}

type yamlChange struct {
	ResourceType string            `yaml:"resource_type"`
	ResourceID   string            `yaml:"resource_id"`
	Attributes   []yamlAttribute   `yaml:"attributes,omitempty"`
}

type yamlAttribute struct {
	Name     string `yaml:"name"`
	Kind     string `yaml:"kind"`
	Want     string `yaml:"want,omitempty"`
	Got      string `yaml:"got,omitempty"`
}

func (f *yamlFormatter) Format(changes []drift.ResourceDiff, w io.Writer) error {
	if w == nil {
		w = os.Stdout
	}

	out := yamlOutput{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Drifted:     len(changes) > 0,
		TotalDrift:  len(changes),
	}

	for _, c := range changes {
		yc := yamlChange{
			ResourceType: c.ResourceType,
			ResourceID:   c.ResourceID,
		}
		for _, a := range c.Attributes {
			yc.Attributes = append(yc.Attributes, yamlAttribute{
				Name: a.Name,
				Kind: fmt.Sprintf("%v", a.Kind),
				Want: a.Want,
				Got:  a.Got,
			})
		}
		out.Changes = append(out.Changes, yc)
	}

	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	return enc.Encode(out)
}
