package runner_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/acme/driftctl-diff/internal/config"
	"github.com/acme/driftctl-diff/internal/runner"
)

func writeTempState(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "terraform.tfstate")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp state: %v", err)
	}
	return p
}

const minimalState = `{
  "version": 4,
  "terraform_version": "1.5.0",
  "resources": []
}`

func TestRunner_EmptyState_NoError(t *testing.T) {
	statePath := writeTempState(t, minimalState)
	cfg := config.DefaultConfig()
	var buf bytes.Buffer

	r := runner.New(cfg, &buf)
	if err := r.Run(statePath); err != nil {
		t.Fatalf("expected no error for empty state, got: %v", err)
	}
}

func TestRunner_EmptyState_OutputContainsSummary(t *testing.T) {
	statePath := writeTempState(t, minimalState)
	cfg := config.DefaultConfig()
	var buf bytes.Buffer

	r := runner.New(cfg, &buf)
	_ = r.Run(statePath)

	out := buf.String()
	if len(out) == 0 {
		t.Error("expected non-empty output, got empty string")
	}
}

func TestRunner_InvalidStatePath_ReturnsError(t *testing.T) {
	cfg := config.DefaultConfig()
	var buf bytes.Buffer

	r := runner.New(cfg, &buf)
	if err := r.Run("/nonexistent/path/terraform.tfstate"); err == nil {
		t.Fatal("expected error for missing state file, got nil")
	}
}

func TestRunner_NilWriter_DefaultsToStdout(t *testing.T) {
	statePath := writeTempState(t, minimalState)
	cfg := config.DefaultConfig()

	// Passing nil should fall back to os.Stdout without panicking.
	r := runner.New(cfg, nil)
	if err := r.Run(statePath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
