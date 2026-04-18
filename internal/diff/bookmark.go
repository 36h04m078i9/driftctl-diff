package diff

import (
	"fmt"
	"strings"
	"time"

	"github.com/acme/driftctl-diff/internal/drift"
)

// Bookmark represents a saved position within a set of drift results.
type Bookmark struct {
	Name      string
	CreatedAt time.Time
	Filter    BookmarkFilter
	Results   []drift.ResourceDiff
}

// BookmarkFilter holds the criteria used when the bookmark was created.
type BookmarkFilter struct {
	ResourceType string
	ResourceID   string
	Kind         string
}

// BookmarkStore holds named bookmarks in memory.
type BookmarkStore struct {
	bookmarks map[string]*Bookmark
}

// NewBookmarkStore returns an empty BookmarkStore.
func NewBookmarkStore() *BookmarkStore {
	return &BookmarkStore{bookmarks: make(map[string]*Bookmark)}
}

// Save stores a bookmark under the given name, overwriting any existing entry.
func (s *BookmarkStore) Save(name string, f BookmarkFilter, results []drift.ResourceDiff) *Bookmark {
	b := &Bookmark{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Filter:    f,
		Results:   results,
	}
	s.bookmarks[name] = b
	return b
}

// Get retrieves a bookmark by name.
func (s *BookmarkStore) Get(name string) (*Bookmark, error) {
	b, ok := s.bookmarks[name]
	if !ok {
		return nil, fmt.Errorf("bookmark %q not found", name)
	}
	return b, nil
}

// Delete removes a bookmark by name.
func (s *BookmarkStore) Delete(name string) {
	delete(s.bookmarks, name)
}

// List returns all bookmark names sorted alphabetically.
func (s *BookmarkStore) List() []string {
	names := make([]string, 0, len(s.bookmarks))
	for n := range s.bookmarks {
		names = append(names, n)
	}
	sortStrings(names)
	return names
}

func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && strings.ToLower(ss[j]) < strings.ToLower(ss[j-1]); j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
