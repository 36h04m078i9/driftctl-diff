package diff

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ArchiverPrinter renders a list of archive entries to a writer.
type ArchiverPrinter struct {
	archiver *Archiver
	w        io.Writer
}

// NewArchiverPrinter creates an ArchiverPrinter. If w is nil, os.Stdout is used.
func NewArchiverPrinter(archiver *Archiver, w io.Writer) *ArchiverPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &ArchiverPrinter{archiver: archiver, w: w}
}

// Print lists all archives and prints a summary of each entry.
func (p *ArchiverPrinter) Print() error {
	paths, err := p.archiver.List()
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		fmt.Fprintln(p.w, "No archives found.")
		return nil
	}
	fmt.Fprintf(p.w, "%-30s %-20s %s\n", "File", "ArchivedAt", "Label")
	fmt.Fprintf(p.w, "%-30s %-20s %s\n", "----", "----------", "-----")
	for _, path := range paths {
		entry, err := p.archiver.Load(path)
		if err != nil {
			fmt.Fprintf(p.w, "%-30s (error reading entry)\n", filepath.Base(path))
			continue
		}
		label := entry.Label
		if label == "" {
			label = "-"
		}
		fmt.Fprintf(p.w, "%-30s %-20s %s\n",
			filepath.Base(path),
			entry.ArchivedAt.Format("2006-01-02T15:04:05"),
			label,
		)
	}
	return nil
}
