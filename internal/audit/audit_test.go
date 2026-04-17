package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/acme/driftctl-diff/internal/audit"
)

func TestRecord_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewWithWriter(&buf)
	e := audit.Entry{
		Timestamp:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		StateFile:  "terraform.tfstate",
		Drifted:    2,
		Total:      5,
		Attributes: 3,
	}
	if err := l.Record(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := strings.TrimSpace(buf.String())
	var got audit.Entry
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.StateFile != "terraform.tfstate" {
		t.Errorf("state_file: got %q", got.StateFile)
	}
	if got.Drifted != 2 {
		t.Errorf("drifted: got %d", got.Drifted)
	}
	if got.Total != 5 {
		t.Errorf("total: got %d", got.Total)
	}
}

func TestRecord_SetsTimestampWhenZero(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewWithWriter(&buf)
	if err := l.Record(audit.Entry{StateFile: "x.tfstate"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got audit.Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewWithWriter(&buf)
	for i := 0; i < 3; i++ {
		if err := l.Record(audit.Entry{StateFile: "s.tfstate", Total: i}); err != nil {
			t.Fatalf("record %d: %v", i, err)
		}
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestNew_EmptyPath_UsesStdout(t *testing.T) {
	l, err := audit.New("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Error("expected non-nil logger")
	}
}
