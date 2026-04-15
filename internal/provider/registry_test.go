package provider_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/snyk/driftctl-diff/internal/provider"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	reg := provider.NewRegistry()
	f := &mockFetcher{attrs: map[string]interface{}{"key": "val"}}
	reg.Register("aws", f)

	got, err := reg.Get("aws")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected fetcher, got nil")
	}
}

func TestRegistry_GetUnknownProvider(t *testing.T) {
	reg := provider.NewRegistry()
	_, err := reg.Get("gcp")
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestRegistry_FetchAttributes_Delegates(t *testing.T) {
	reg := provider.NewRegistry()
	expected := map[string]interface{}{"bucket": "b", "region": "eu-west-1"}
	reg.Register("aws", &mockFetcher{attrs: expected})

	got, err := reg.FetchAttributes(context.Background(), "aws", "aws_s3_bucket", "b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["bucket"] != "b" {
		t.Errorf("expected bucket b, got %v", got["bucket"])
	}
}

func TestRegistry_FetchAttributes_UnknownProvider(t *testing.T) {
	reg := provider.NewRegistry()
	_, err := reg.FetchAttributes(context.Background(), "azure", "azurerm_resource_group", "rg")
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

// Ensure mockFetcher satisfies the interface (compile-time check).
var _ provider.ResourceFetcher = (*mockFetcher)(nil)

// Re-declare here to keep test file self-contained (aws_test.go is same package).
var _ = fmt.Errorf // suppress unused import
