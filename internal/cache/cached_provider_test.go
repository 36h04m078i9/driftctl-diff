package cache_test

import (
	"errors"
	"testing"

	"github.com/you/driftctl-diff/internal/cache"
)

// stubFetcher records how many times it has been called.
type stubFetcher struct {
	calls  int
	attrs  map[string]interface{}
	err    error
}

func (s *stubFetcher) FetchAttributes(_, _ string) (map[string]interface{}, error) {
	s.calls++
	return s.attrs, s.err
}

func TestCachedProvider_HitsInnerOnFirstCall(t *testing.T) {
	stub := &stubFetcher{attrs: map[string]interface{}{"key": "val"}}
	cp := cache.NewCachedProvider(stub)

	attrs, err := cp.FetchAttributes("aws_s3_bucket", "my-bucket")
	if err != nil {
		t.Fatal(err)
	}
	if attrs["key"] != "val" {
		t.Fatalf("unexpected attrs: %v", attrs)
	}
	if stub.calls != 1 {
		t.Fatalf("expected 1 call, got %d", stub.calls)
	}
}

func TestCachedProvider_ReturnsCachedOnSecondCall(t *testing.T) {
	stub := &stubFetcher{attrs: map[string]interface{}{"key": "val"}}
	cp := cache.NewCachedProvider(stub)

	cp.FetchAttributes("aws_s3_bucket", "my-bucket") //nolint:errcheck
	cp.FetchAttributes("aws_s3_bucket", "my-bucket") //nolint:errcheck

	if stub.calls != 1 {
		t.Fatalf("expected 1 inner call, got %d", stub.calls)
	}
	if cp.CacheLen() != 1 {
		t.Fatalf("expected cache len 1, got %d", cp.CacheLen())
	}
}

func TestCachedProvider_PropagatesError(t *testing.T) {
	stub := &stubFetcher{err: errors.New("api error")}
	cp := cache.NewCachedProvider(stub)

	_, err := cp.FetchAttributes("aws_s3_bucket", "bad-bucket")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// Nothing should be cached on error.
	if cp.CacheLen() != 0 {
		t.Fatalf("expected empty cache after error, got %d", cp.CacheLen())
	}
}

func TestCachedProvider_Invalidate(t *testing.T) {
	stub := &stubFetcher{attrs: map[string]interface{}{"key": "val"}}
	cp := cache.NewCachedProvider(stub)

	cp.FetchAttributes("aws_s3_bucket", "my-bucket") //nolint:errcheck
	cp.Invalidate("aws_s3_bucket", "my-bucket")
	cp.FetchAttributes("aws_s3_bucket", "my-bucket") //nolint:errcheck

	if stub.calls != 2 {
		t.Fatalf("expected 2 inner calls after invalidation, got %d", stub.calls)
	}
}
