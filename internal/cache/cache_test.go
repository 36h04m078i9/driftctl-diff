package cache_test

import (
	"testing"

	"github.com/you/driftctl-diff/internal/cache"
)

func TestNew_IsEmpty(t *testing.T) {
	c := cache.New()
	if c.Len() != 0 {
		t.Fatalf("expected empty cache, got len %d", c.Len())
	}
}

func TestSetAndGet(t *testing.T) {
	c := cache.New()
	attrs := map[string]interface{}{"region": "us-east-1"}
	c.Set("aws_s3_bucket", "my-bucket", attrs)

	got, ok := c.Get("aws_s3_bucket", "my-bucket")
	if !ok {
		t.Fatal("expected entry to be present")
	}
	if got["region"] != "us-east-1" {
		t.Fatalf("expected us-east-1, got %v", got["region"])
	}
}

func TestGet_Missing(t *testing.T) {
	c := cache.New()
	_, ok := c.Get("aws_s3_bucket", "no-such-bucket")
	if ok {
		t.Fatal("expected miss, got hit")
	}
}

func TestDelete(t *testing.T) {
	c := cache.New()
	c.Set("aws_s3_bucket", "my-bucket", map[string]interface{}{})
	c.Delete("aws_s3_bucket", "my-bucket")
	_, ok := c.Get("aws_s3_bucket", "my-bucket")
	if ok {
		t.Fatal("expected entry to be deleted")
	}
}

func TestFlush(t *testing.T) {
	c := cache.New()
	c.Set("aws_s3_bucket", "b1", map[string]interface{}{})
	c.Set("aws_s3_bucket", "b2", map[string]interface{}{})
	c.Flush()
	if c.Len() != 0 {
		t.Fatalf("expected empty cache after flush, got %d", c.Len())
	}
}

func TestLen(t *testing.T) {
	c := cache.New()
	c.Set("aws_instance", "i-001", map[string]interface{}{})
	c.Set("aws_instance", "i-002", map[string]interface{}{})
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}

func TestSet_Overwrite(t *testing.T) {
	c := cache.New()
	c.Set("aws_instance", "i-001", map[string]interface{}{"state": "running"})
	c.Set("aws_instance", "i-001", map[string]interface{}{"state": "stopped"})
	got, _ := c.Get("aws_instance", "i-001")
	if got["state"] != "stopped" {
		t.Fatalf("expected stopped, got %v", got["state"])
	}
	if c.Len() != 1 {
		t.Fatalf("overwrite should not increase len, got %d", c.Len())
	}
}
