package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeBookmarkDiffs(n int) []drift.ResourceDiff {
	out := make([]drift.ResourceDiff, n)
	for i := range out {
		out[i] = drift.ResourceDiff{ResourceID: "id", ResourceType: "aws_s3_bucket"}
	}
	return out
}

func TestBookmarkStore_SaveAndGet(t *testing.T) {
	s := NewBookmarkStore()
	f := BookmarkFilter{ResourceType: "aws_s3_bucket"}
	b := s.Save("my-bm", f, makeBookmarkDiffs(3))
	if b.Name != "my-bm" {
		t.Fatalf("expected name my-bm, got %s", b.Name)
	}
	got, err := s.Get("my-bm")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(got.Results))
	}
}

func TestBookmarkStore_GetMissing(t *testing.T) {
	s := NewBookmarkStore()
	_, err := s.Get("nope")
	if err == nil {
		t.Fatal("expected error for missing bookmark")
	}
}

func TestBookmarkStore_Delete(t *testing.T) {
	s := NewBookmarkStore()
	s.Save("bm", BookmarkFilter{}, makeBookmarkDiffs(1))
	s.Delete("bm")
	if len(s.List()) != 0 {
		t.Fatal("expected empty list after delete")
	}
}

func TestBookmarkStore_List_Sorted(t *testing.T) {
	s := NewBookmarkStore()
	s.Save("zebra", BookmarkFilter{}, nil)
	s.Save("apple", BookmarkFilter{}, nil)
	s.Save("mango", BookmarkFilter{}, nil)
	names := s.List()
	if names[0] != "apple" || names[1] != "mango" || names[2] != "zebra" {
		t.Fatalf("unexpected order: %v", names)
	}
}

func TestBookmarkPrinter_Empty(t *testing.T) {
	var buf bytes.Buffer
	p := NewBookmarkPrinter(&buf)
	p.Print(NewBookmarkStore())
	if !strings.Contains(buf.String(), "no bookmarks") {
		t.Fatalf("expected no-bookmarks message, got: %s", buf.String())
	}
}

func TestBookmarkPrinter_ContainsName(t *testing.T) {
	var buf bytes.Buffer
	p := NewBookmarkPrinter(&buf)
	s := NewBookmarkStore()
	s.Save("prod-drift", BookmarkFilter{ResourceType: "aws_instance"}, makeBookmarkDiffs(2))
	p.Print(s)
	out := buf.String()
	if !strings.Contains(out, "prod-drift") {
		t.Fatalf("expected bookmark name in output: %s", out)
	}
	if !strings.Contains(out, "aws_instance") {
		t.Fatalf("expected resource type in output: %s", out)
	}
}
