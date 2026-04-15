package cache

import "fmt"

// AttributeFetcher is the interface satisfied by provider implementations
// (mirrors the contract in internal/provider).
type AttributeFetcher interface {
	FetchAttributes(resourceType, resourceID string) (map[string]interface{}, error)
}

// CachedProvider wraps any AttributeFetcher and transparently caches results
// so that repeated lookups for the same resource do not trigger additional
// cloud API calls within the same run.
type CachedProvider struct {
	inner AttributeFetcher
	cache *Cache
}

// NewCachedProvider returns a CachedProvider backed by the supplied fetcher.
func NewCachedProvider(inner AttributeFetcher) *CachedProvider {
	return &CachedProvider{
		inner: inner,
		cache: New(),
	}
}

// FetchAttributes returns cached attributes when available, otherwise delegates
// to the inner fetcher, stores the result, and returns it.
func (cp *CachedProvider) FetchAttributes(resourceType, resourceID string) (map[string]interface{}, error) {
	if attrs, ok := cp.cache.Get(resourceType, resourceID); ok {
		return attrs, nil
	}

	attrs, err := cp.inner.FetchAttributes(resourceType, resourceID)
	if err != nil {
		return nil, fmt.Errorf("cached_provider: fetch %s/%s: %w", resourceType, resourceID, err)
	}

	cp.cache.Set(resourceType, resourceID, attrs)
	return attrs, nil
}

// CacheLen exposes the current number of cached entries (useful for diagnostics
// and testing).
func (cp *CachedProvider) CacheLen() int {
	return cp.cache.Len()
}

// Invalidate removes a single entry from the cache, forcing the next call to
// FetchAttributes to hit the live provider.
func (cp *CachedProvider) Invalidate(resourceType, resourceID string) {
	cp.cache.Delete(resourceType, resourceID)
}
