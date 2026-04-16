package snapshot_test

import (
	"errors"
	"testing"

	"github.com/example/driftctl-diff/internal/provider"
	"github.com/example/driftctl-diff/internal/snapshot"
	"github.com/example/driftctl-diff/internal/state"
)

type mockFetcher struct {
	attrs map[string]string
	err   error
}

func (m *mockFetcher) FetchAttributes(resourceType, id string) (map[string]string, error) {
	return m.attrs, m.err
}

func makeRegistry(t *testing.T, fetcher provider.Fetcher) *provider.Registry {
	t.Helper()
	reg := provider.NewRegistry()
	reg.Register("aws_instance", fetcher)
	return reg
}

func TestCollect_Success(t *testing.T) {
	attrs := map[string]string{"ami": "ami-abc", "instance_type": "t3.micro"}
	reg := makeRegistry(t, &mockFetcher{attrs: attrs})
	c := snapshot.NewCollector(reg)

	resources := []state.Resource{
		{Type: "aws_instance", Name: "web", ID: "i-001"},
	}

	snap, err := c.Collect(resources)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, ok := snap.Get("aws_instance.web")
	if !ok {
		t.Fatal("expected resource in snapshot")
	}
	if got["ami"] != "ami-abc" {
		t.Errorf("unexpected ami: %s", got["ami"])
	}
}

func TestCollect_FetchError_PartialSnapshot(t *testing.T) {
	reg := makeRegistry(t, &mockFetcher{err: errors.New("api error")})
	c := snapshot.NewCollector(reg)

	resources := []state.Resource{
		{Type: "aws_instance", Name: "web", ID: "i-001"},
	}

	snap, err := c.Collect(resources)
	if err == nil {
		t.Error("expected error")
	}
	if snap == nil {
		t.Error("expected partial snapshot even on error")
	}
}

func TestCollect_EmptyResources(t *testing.T) {
	reg := provider.NewRegistry()
	c := snapshot.NewCollector(reg)

	snap, err := c.Collect([]state.Resource{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Resources) != 0 {
		t.Errorf("expected empty snapshot, got %d resources", len(snap.Resources))
	}
}
