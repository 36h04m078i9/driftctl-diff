package output

import (
	"fmt"
	"io"
	"os"

	"github.com/snyk/driftctl-diff/internal/drift"
)

// Writer handles writing formatted drift results to an output destination.
type Writer struct {
	formatter *Formatter
	dest      io.Writer
}

// NewWriter creates a Writer that writes to the given destination.
// If dest is nil, os.Stdout is used.
func NewWriter(f *Formatter, dest io.Writer) *Writer {
	if dest == nil {
		dest = os.Stdout
	}
	return &Writer{
		formatter: f,
		dest:      dest,
	}
}

// Write formats the drift results and writes them to the configured destination.
// It returns the number of bytes written and any error encountered.
func (w *Writer) Write(results []drift.ResourceDiff) (int, error) {
	if len(results) == 0 {
		n, err := fmt.Fprintln(w.dest, "No drift detected.")
		return n, err
	}

	formatted := w.formatter.Format(results)
	n, err := fmt.Fprint(w.dest, formatted)
	return n, err
}

// WriteTo satisfies io.WriterTo for a pre-formatted string, useful in pipelines.
func (w *Writer) WriteTo(results []drift.ResourceDiff, dest io.Writer) (int64, error) {
	if len(results) == 0 {
		n, err := fmt.Fprintln(dest, "No drift detected.")
		return int64(n), err
	}

	formatted := w.formatter.Format(results)
	n, err := fmt.Fprint(dest, formatted)
	return int64(n), err
}
