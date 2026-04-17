package metrics

import (
	"testing"
	"time"
)

func TestNew_ZeroCounters(t *testing.T) {
	c := New()
	snap := c.Snapshot()
	if snap.ResourcesTotal != 0 || snap.ResourcesDrifted != 0 {
		t.Fatalf("expected zero counters, got %+v", snap)
	}
}

func TestIncResources(t *testing.T) {
	c := New()
	c.IncResources(5)
	c.IncResources(3)
	if got := c.Snapshot().ResourcesTotal; got != 8 {
		t.Fatalf("expected 8, got %d", got)
	}
}

func TestIncDrifted(t *testing.T) {
	c := New()
	c.IncDrifted(2)
	if got := c.Snapshot().ResourcesDrifted; got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestIncAttributes(t *testing.T) {
	c := New()
	c.IncAttributes(10)
	if got := c.Snapshot().AttributesChecked; got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestIncFetchErrors(t *testing.T) {
	c := New()
	c.IncFetchErrors(1)
	if got := c.Snapshot().FetchErrors; got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestSnapshot_Duration(t *testing.T) {
	c := New()
	time.Sleep(5 * time.Millisecond)
	snap := c.Snapshot()
	if snap.Duration < 5*time.Millisecond {
		t.Fatalf("expected duration >= 5ms, got %s", snap.Duration)
	}
}
