package diff_test

import (
	"bytes"
	"strings"
	"testing"

	pagerpkg "github.com/snyk/driftctl-diff/internal/diff"
	"github.com/snyk/driftctl-diff/internal/drift"
)

func sampleDiffs(n int) []drift.ResourceDiff {
	out := make([]drift.ResourceDiff, n)
	for i := 0; i < n; i++ {
		out[i] = drift.ResourceDiff{
			ResourceID:   fmt.Sprintf("res-%d", i),
			ResourceType: "aws_instance",
		}
	}
	return out
}

func TestPaginate_Empty(t *testing.T) {
	p := pagerpkg.NewPager(5, nil)
	pages := p.Paginate([]drift.ResourceDiff{})
	if len(pages) != 0 {
		t.Fatalf("expected 0 pages, got %d", len(pages))
	}
}

func TestPaginate_SinglePage(t *testing.T) {
	p := pagerpkg.NewPager(10, nil)
	results := make([]drift.ResourceDiff, 3)
	pages := p.Paginate(results)
	if len(pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(pages))
	}
	if pages[0].Total != 1 {
		t.Errorf("expected total=1, got %d", pages[0].Total)
	}
}

func TestPaginate_MultiplePages(t *testing.T) {
	p := pagerpkg.NewPager(3, nil)
	results := make([]drift.ResourceDiff, 7)
	pages := p.Paginate(results)
	if len(pages) != 3 {
		t.Fatalf("expected 3 pages, got %d", len(pages))
	}
	if len(pages[2].Results) != 1 {
		t.Errorf("last page should have 1 result, got %d", len(pages[2].Results))
	}
}

func TestPaginate_PageNumbers(t *testing.T) {
	p := pagerpkg.NewPager(2, nil)
	results := make([]drift.ResourceDiff, 4)
	pages := p.Paginate(results)
	for i, pg := range pages {
		if pg.Number != i+1 {
			t.Errorf("page %d: expected Number=%d, got %d", i, i+1, pg.Number)
		}
	}
}

func TestPrintPage_ContainsResourceID(t *testing.T) {
	var buf bytes.Buffer
	p := pagerpkg.NewPager(5, &buf)
	pg := pagerpkg.Page{
		Number:  1,
		Total:   1,
		Results: []drift.ResourceDiff{{ResourceID: "my-instance", ResourceType: "aws_instance"}},
	}
	p.PrintPage(pg)
	if !strings.Contains(buf.String(), "my-instance") {
		t.Errorf("expected output to contain resource id, got: %s", buf.String())
	}
}

func TestNewPager_ZeroPageSize_Defaults(t *testing.T) {
	p := pagerpkg.NewPager(0, nil)
	results := make([]drift.ResourceDiff, 25)
	pages := p.Paginate(results)
	if len(pages) != 2 {
		t.Errorf("expected 2 pages with default size 20, got %d", len(pages))
	}
}
