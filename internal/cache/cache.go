// Package cache provides a simple in-memory cache for live cloud resource attributes,
// reducing redundant API calls during a single drift-detection run.
package cache

import "sync"

// Entry holds a cached set of attributes for a single cloud resource.
type Entry struct {
	Attributes map[string]interface{}
}

// Cache is a thread-safe, in-memory store keyed by resource type + resource ID.
type Cache struct {
	mu    sync.RWMutex
	store map[string]Entry
}

// New returns an initialised, empty Cache.
func New() *Cache {
	return &Cache{
		store: make(map[string]Entry),
	}
}

// key builds the composite lookup key.
func key(resourceType, resourceID string) string {
	return resourceType + "/" + resourceID
}

// Set stores attributes for the given resource type and ID.
func (c *Cache) Set(resourceType, resourceID string, attrs map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key(resourceType, resourceID)] = Entry{Attributes: attrs}
}

// Get retrieves attributes for the given resource type and ID.
// The second return value reports whether the entry was found.
func (c *Cache) Get(resourceType, resourceID string) (map[string]interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.store[key(resourceType, resourceID)]
	if !ok {
		return nil, false
	}
	return e.Attributes, true
}

// Delete removes an entry from the cache.
func (c *Cache) Delete(resourceType, resourceID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key(resourceType, resourceID))
}

// Len returns the number of entries currently held in the cache.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.store)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]Entry)
}
