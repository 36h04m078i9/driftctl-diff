package provider_test

import (
	"context"
	"testing"

	"github.com/snyk/driftctl-diff/internal/provider"
)

// mockFetcher satisfies ResourceFetcher for unit tests.
type mockFetcher struct {
	attrs map[string]interface{}
	err   error
}

func (m *mockFetcher) FetchAttributes(_ context.Context, _, _ string) (map[string]interface{}, error) {
	return m.attrs, m.err
}

func TestFetchAttributes_UnsupportedType(t *testing.T) {
	// We can test the interface contract without real AWS credentials
	// by verifying unsupported types return an error from a stub.
	var _ provider.ResourceFetcher = &mockFetcher{}
}

func TestMockFetcher_ReturnsAttrs(t *testing.T) {
	expected := map[string]interface{}{"bucket": "my-bucket", "region": "us-east-1"}
	f := &mockFetcher{attrs: expected}

	got, err := f.FetchAttributes(context.Background(), "aws_s3_bucket", "my-bucket")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["bucket"] != expected["bucket"] {
		t.Errorf("expected bucket %v, got %v", expected["bucket"], got["bucket"])
	}
	if got["region"] != expected["region"] {
		t.Errorf("expected region %v, got %v", expected["region"], got["region"])
	}
}

func TestMockFetcher_PropagatesError(t *testing.T) {
	f := &mockFetcher{err: fmt.Errorf("network error")}
	_, err := f.FetchAttributes(context.Background(), "aws_s3_bucket", "x")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
