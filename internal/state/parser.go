package state

import (
	"encoding/json"
	"fmt"
	"os"
)

// Resource represents a single resource entry from Terraform state.
type Resource struct {
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Provider   string                 `json:"provider"`
	Attributes map[string]interface{} `json:"attributes"`
}

// State holds the parsed contents of a Terraform state file.
type State struct {
	Version   int        `json:"version"`
	Resources []Resource `json:"resources"`
}

// Parser reads and decodes Terraform state files.
type Parser struct{}

// NewParser creates a new Parser instance.
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile reads the file at path and returns a decoded State.
func (p *Parser) ParseFile(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading state file: %w", err)
	}
	return p.Parse(data)
}

// Parse decodes raw JSON bytes into a State struct.
func (p *Parser) Parse(data []byte) (*State, error) {
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("decoding state JSON: %w", err)
	}
	if s.Version < 4 {
		return nil, fmt.Errorf("unsupported state version %d (minimum: 4)", s.Version)
	}
	return &s, nil
}
