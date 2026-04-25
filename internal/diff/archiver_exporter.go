package diff

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// ArchiverExporter exports a loaded ArchiveEntry to an output writer.
type ArchiverExporter struct {
	archiver *Archiver
	w        io.Writer
}

// NewArchiverExporter creates an ArchiverExporter. If w is nil, os.Stdout is used.
func NewArchiverExporter(archiver *Archiver, w io.Writer) *ArchiverExporter {
	if w == nil {
		w = os.Stdout
	}
	return &ArchiverExporter{archiver: archiver, w: w}
}

// ExportJSON writes the archive entry at path as indented JSON.
func (e *ArchiverExporter) ExportJSON(path string) error {
	entry, err := e.archiver.Load(path)
	if err != nil {
		return fmt.Errorf("archiver exporter: load: %w", err)
	}
	enc := json.NewEncoder(e.w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("archiver exporter: encode: %w", err)
	}
	return nil
}

// ExportText writes a human-readable summary of the archive entry.
func (e *ArchiverExporter) ExportText(path string) error {
	entry, err := e.archiver.Load(path)
	if err != nil {
		return fmt.Errorf("archiver exporter: load: %w", err)
	}
	fmt.Fprintf(e.w, "Archive: %s\n", path)
	fmt.Fprintf(e.w, "Archived At: %s\n", entry.ArchivedAt.Format("2006-01-02T15:04:05Z"))
	if entry.Label != "" {
		fmt.Fprintf(e.w, "Label: %s\n", entry.Label)
	}
	fmt.Fprintf(e.w, "Resources: %d\n", len(entry.Results))
	for _, r := range entry.Results {
		fmt.Fprintf(e.w, "  [%s] %s (%d change(s))\n", r.ResourceType, r.ResourceID, len(r.Changes))
	}
	return nil
}
