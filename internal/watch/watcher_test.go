package watch_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/driftctl-diff/internal/watch"
)

// stubRunner satisfies the interface used by Watcher via duck-typing through
// a thin wrapper so we don't import the real runner in tests.
type stubRunner struct {
	calls int
}

func (s *stubRunner) Detect(_ context.Context) ([]interface{}, error) {
	s.calls++
	return nil, nil
}

func TestWatcher_EmitsResultsOnTick(t *testing.T) {
	// Use a very short interval so the test finishes quickly.
	const interval = 20 * time.Millisecond

	w := watch.NewWithDetector(func(_ context.Context) ([]interface{}, error) {
		return nil, nil
	}, interval)

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()

	go w.Start(ctx)

	var count int
	for range w.Results() {
		count++
	}

	if count < 2 {
		t.Fatalf("expected at least 2 results, got %d", count)
	}
}

func TestWatcher_StopsOnContextCancel(t *testing.T) {
	w := watch.NewWithDetector(func(_ context.Context) ([]interface{}, error) {
		return nil, nil
	}, 10*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	done := make(chan struct{})
	go func() {
		w.Start(ctx)
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(200 * time.Millisecond):
		t.Fatal("watcher did not stop after context cancellation")
	}
}
