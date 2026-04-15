package provider

import (
	"context"
	"fmt"
)

// Registry maps provider names to their ResourceFetcher implementations.
type Registry struct {
	fetchers map[string]ResourceFetcher
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{fetchers: make(map[string]ResourceFetcher)}
}

// Register adds a fetcher under the given provider name (e.g. "aws").
func (r *Registry) Register(name string, f ResourceFetcher) {
	r.fetchers[name] = f
}

// Get returns the fetcher registered under name, or an error if absent.
func (r *Registry) Get(name string) (ResourceFetcher, error) {
	f, ok := r.fetchers[name]
	if !ok {
		return nil, fmt.Errorf("no provider registered for %q", name)
	}
	return f, nil
}

// FetchAttributes is a convenience method that resolves the provider by name
// and delegates to its FetchAttributes.
func (r *Registry) FetchAttributes(ctx context.Context, providerName, resourceType, resourceID string) (map[string]interface{}, error) {
	f, err := r.Get(providerName)
	if err != nil {
		return nil, err
	}
	return f.FetchAttributes(ctx, resourceType, resourceID)
}
