package diff

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// CheckpointPrinter renders a list of checkpoints in a tabular format.
type CheckpointPrinter struct {
	store *CheckpointStore
	w     io.Writer
}

// NewCheckpointPrinter returns a CheckpointPrinter writing to w.
// If w is nil it defaults to os.Stdout.
func NewCheckpointPrinter(store *CheckpointStore, w io.Writer) *CheckpointPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &CheckpointPrinter{store: store, w: w}
}

// Print lists all checkpoints with their creation time and result count.
func (p *CheckpointPrinter) Print() error {
	names, err := p.store.List()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Fprintln(p.w, "no checkpoints found")
		return nil
	}
	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NAME\tCREATED AT\tRESOURCES")
	for _, name := range names {
		cp, err := p.store.Load(name)
		if err != nil {
			fmt.Fprintf(tw, "%s\t(error loading)\t-\n", name)
			continue
		}
		fmt.Fprintf(tw, "%s\t%s\t%d\n",
			cp.Name,
			cp.CreatedAt.Format("2006-01-02 15:04:05"),
			len(cp.Results),
		)
	}
	return tw.Flush()
}
