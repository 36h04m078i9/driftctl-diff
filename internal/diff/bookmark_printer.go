package diff

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// BookmarkPrinter renders a BookmarkStore's contents to a writer.
type BookmarkPrinter struct {
	w io.Writer
}

// NewBookmarkPrinter returns a BookmarkPrinter writing to w.
// If w is nil it defaults to os.Stdout.
func NewBookmarkPrinter(w io.Writer) *BookmarkPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &BookmarkPrinter{w: w}
}

// Print writes a human-readable table of bookmarks to the writer.
func (p *BookmarkPrinter) Print(store *BookmarkStore) {
	names := store.List()
	if len(names) == 0 {
		fmt.Fprintln(p.w, "no bookmarks saved")
		return
	}
	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NAME\tCREATED\tRESOURCES\tTYPE FILTER\tID FILTER")
	for _, name := range names {
		b, _ := store.Get(name)
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\t%s\n",
			b.Name,
			b.CreatedAt.Format("2006-01-02 15:04:05"),
			len(b.Results),
			orDash(b.Filter.ResourceType),
			orDash(b.Filter.ResourceID),
		)
	}
	tw.Flush()
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
