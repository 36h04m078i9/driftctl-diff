// Package audit records drift scan events to an append-only log file.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp  time.Time `json:"timestamp"`
	StateFile  string    `json:"state_file"`
	Drifted    int       `json:"drifted_resources"`
	Total      int       `json:"total_resources"`
	Attributes int       `json:"drifted_attributes"`
}

// Logger writes audit entries to an underlying writer.
type Logger struct {
	w io.Writer
}

// New returns a Logger that appends JSON lines to path.
// If path is empty, os.Stdout is used.
func New(path string) (*Logger, error) {
	if path == "" {
		return &Logger{w: os.Stdout}, nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("audit: open %s: %w", path, err)
	}
	return &Logger{w: f}, nil
}

// NewWithWriter returns a Logger backed by w (useful in tests).
func NewWithWriter(w io.Writer) *Logger {
	return &Logger{w: w}
}

// Record encodes e as a JSON line and writes it to the log.
func (l *Logger) Record(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}
