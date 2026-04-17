package policy

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// File is the top-level structure of a policy YAML file.
type File struct {
	Rules []Rule `yaml:"rules"`
}

// LoadFile reads and parses a policy file from the given path.
// If path is empty, an empty rule set is returned.
func LoadFile(path string) ([]Rule, error) {
	if path == "" {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("policy: read file: %w", err)
	}
	var f File
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("policy: parse yaml: %w", err)
	}
	for i, r := range f.Rules {
		if r.ResourceType == "" {
			return nil, fmt.Errorf("policy: rule %d missing resource_type", i)
		}
		switch r.Severity {
		case SeverityLow, SeverityMedium, SeverityHigh:
		default:
			return nil, fmt.Errorf("policy: rule %d has invalid severity %q", i, r.Severity)
		}
	}
	return f.Rules, nil
}
