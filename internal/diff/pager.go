// Package diff provides utilities for paginating drift results.
package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/snyk/driftctl-diff/internal/drift"
)

// Page holds a slice of results for a single page.
type Page struct {
	Number  int
	Total   int
	Results []drift.ResourceDiff
}

// Pager splits a slice of ResourceDiff into fixed-size pages.
type Pager struct {
	pageSize int
	out      io.Writer
}

// NewPager creates a Pager with the given page size.
func NewPager(pageSize int, out io.Writer) *Pager {
	if pageSize <= 0 {
		pageSize = 20
	}
	if out == nil {
		out = os.Stdout
	}
	return &Pager{pageSize: pageSize, out: out}
}

// Paginate divides results into pages and returns them all.
func (p *Pager) Paginate(results []drift.ResourceDiff) []Page {
	if len(results) == 0 {
		return []Page{}
	}
	total := (len(results) + p.pageSize - 1) / p.pageSize
	pages := make([]Page, 0, total)
	for i := 0; i < total; i++ {
		start := i * p.pageSize
		end := start + p.pageSize
		if end > len(results) {
			end = len(results)
		}
		pages = append(pages, Page{
			Number:  i + 1,
			Total:   total,
			Results: results[start:end],
		})
	}
	return pages
}

// PrintPage writes a summary of the given page to the writer.
func (p *Pager) PrintPage(pg Page) {
	fmt.Fprintf(p.out, "Page %d/%d (%d results)\n", pg.Number, pg.Total, len(pg.Results))
	for _, r := range pg.Results {
		fmt.Fprintf(p.out, "  ~ %s (%s)\n", r.ResourceID, r.ResourceType)
	}
}
