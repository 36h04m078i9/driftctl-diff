package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestValidatorPrinter_Valid_PrintsPassMessage(t *testing.T) {
	var buf bytes.Buffer
	p := NewValidatorPrinter(&buf)
	p.Print(ValidationResult{Valid: true})
	if !strings.Contains(buf.String(), "passed") {
		t.Errorf("expected 'passed' in output, got: %s", buf.String())
	}
}

func TestValidatorPrinter_Invalid_PrintsErrors(t *testing.T) {
	var buf bytes.Buffer
	p := NewValidatorPrinter(&buf)
	p.Print(ValidationResult{Valid: false, Errors: []string{"result[0]: empty resource_id"}})
	out := buf.String()
	if !strings.Contains(out, "failed") {
		t.Errorf("expected 'failed' in output")
	}
	if !strings.Contains(out, "empty resource_id") {
		t.Errorf("expected error detail in output")
	}
}

func TestValidatorPrinter_Warnings_Printed(t *testing.T) {
	var buf bytes.Buffer
	p := NewValidatorPrinter(&buf)
	p.Print(ValidationResult{Valid: true, Warnings: []string{"too many changes"}})
	if !strings.Contains(buf.String(), "WARNING") {
		t.Errorf("expected WARNING in output")
	}
}

func TestValidatorPrinter_NilWriter_DefaultsToStdout(t *testing.T) {
	p := NewValidatorPrinter(nil)
	if p.w == nil {
		t.Error("expected non-nil writer")
	}
}
