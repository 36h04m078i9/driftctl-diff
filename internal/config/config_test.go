package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/driftctl-diff/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.StatePath != "terraform.tfstate" {
		t.Errorf("expected default state_path, got %q", cfg.StatePath)
	}
	if cfg.Provider != "aws" {
		t.Errorf("expected default provider aws, got %q", cfg.Provider)
	}
	if !cfg.Color {
		t.Error("expected color to be enabled by default")
	}
}

func TestLoad_EmptyPath(t *testing.T) {
	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}

func TestLoad_ValidFile(t *testing.T) {
	content := []byte("state_path: custom.tfstate\nprovider: gcp\nregion: eu-west-1\ncolor: false\n")
	tmp := filepath.Join(t.TempDir(), "driftctl.yaml")
	if err := os.WriteFile(tmp, content, 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	cfg, err := config.Load(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.StatePath != "custom.tfstate" {
		t.Errorf("expected custom.tfstate, got %q", cfg.StatePath)
	}
	if cfg.Provider != "gcp" {
		t.Errorf("expected gcp, got %q", cfg.Provider)
	}
	if cfg.Color {
		t.Error("expected color disabled")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/driftctl.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.yaml")
	if err := os.WriteFile(tmp, []byte(": : invalid: yaml:"), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	_, err := config.Load(tmp)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}
