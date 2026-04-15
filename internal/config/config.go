package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level CLI configuration.
type Config struct {
	StatePath  string   `yaml:"state_path"`
	Provider   string   `yaml:"provider"`
	Region     string   `yaml:"region"`
	OutputFmt  string   `yaml:"output_format"`
	Color      bool     `yaml:"color"`
	Ignore     []string `yaml:"ignore"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		StatePath: "terraform.tfstate",
		Provider:  "aws",
		Region:    "us-east-1",
		OutputFmt: "text",
		Color:     true,
		Ignore:    []string{},
	}
}

// Load reads a YAML config file from path and merges it over defaults.
// If path is empty, defaults are returned unchanged.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}
	return cfg, nil
}
