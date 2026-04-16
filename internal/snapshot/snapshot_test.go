package snapshot_test

import (
	"os"
	"testing"
	"time"

	"github.com/example/driftctl-diff/internal/snapshot"
)

func TestNew_IsEmpty(t *testing.T) {
	s := snapshot.New()
	if s == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if len(s.Resources) != 0 {
		t.Errorf("expected empty resources, got %d", len(s.Resources))
	}
	if s.CapturedAt.IsZero() {
		t.Error("expected CapturedAt to be set")
	}
}

func TestAdd_AndGet(t *testing.T) {
	s := snapshot.New()
	attrs := map[string]string{"region": "us-east-1", "id": "i-123"}
	s.Add("aws_instance.web", attrs)

	got, ok := s.Get("aws_instance.web")
	if !ok {
		t.Fatal("expected resource to be found")
	}
	if got["region"] != "us-east-1" {
		t.Errorf("unexpected region: %s", got["region"])
	}
}

func TestGet_Missing(t *testing.T) {
	s := snapshot.New()
	_, ok := s.Get("nonexistent")
	if ok {
		t.Error("expected missing resource to return false")
	}
}

func TestSaveTo_AndLoadFrom_RoundTrip(t *testing.T) {
	s := snapshot.New()
	s.CapturedAt = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	s.Add("aws_s3_bucket.logs", map[string]string{"versioning": "enabled"})

	tmp, err := os.CreateTemp("", "snapshot-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	if err := s.SaveTo(tmp.Name()); err != nil {
		t.Fatalf("SaveTo: %v", err)
	}

	loaded, err := snapshot.LoadFrom(tmp.Name())
	if err != nil {
		t.Fatalf("LoadFrom: %v", err)
	}

	if !loaded.CapturedAt.Equal(s.CapturedAt) {
		t.Errorf("CapturedAt mismatch: got %v want %v", loaded.CapturedAt, s.CapturedAt)
	}
	attrs, ok := loaded.Get("aws_s3_bucket.logs")
	if !ok {
		t.Fatal("expected resource in loaded snapshot")
	}
	if attrs["versioning"] != "enabled" {
		t.Errorf("unexpected versioning value: %s", attrs["versioning"])
	}
}

func TestLoadFrom_MissingFile(t *testing.T) {
	_, err := snapshot.LoadFrom("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
