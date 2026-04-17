// Package diff provides pagination utilities for drift results.
//
// Use NewPager to create a pager with a fixed page size, then call
// Paginate to split a []drift.ResourceDiff slice into pages.
// PrintPage writes a human-readable summary of each page to the
// configured writer.
//
// Example:
//
//	p := diff.NewPager(20, os.Stdout)
//	pages := p.Paginate(results)
//	for _, pg := range pages {
//		p.PrintPage(pg)
//	}
package diff
