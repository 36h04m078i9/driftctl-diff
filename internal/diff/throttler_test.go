package diff

import (
	"testing"
	"time"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeThrottlerResults(n int) []drift.DriftResult {
	out := make([]drift.DriftResult, n)
	for i := range out {
		out[i] = drift.DriftResult{
			ResourceID:   fmt.Sprintf("res-%d", i),
			ResourceType: "aws_instance",
		}
	}
	return out
}

func TestThrottler_NoLimit_ReturnsAll(t *testing.T) {
	th := NewThrottler(DefaultThrottleOptions())
	input := makeThrottlerResults(5)
	out, err := th.Throttle(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(input) {
		t.Errorf("expected %d results, got %d", len(input), len(out))
	}
}

func TestThrottler_EmptyInput_ReturnsEmpty(t *testing.T) {
	th := NewThrottler(ThrottleOptions{MaxPerSecond: 10})
	out, err := th.Throttle(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d items", len(out))
	}
}

func TestThrottler_SleepsAfterBurst(t *testing.T) {
	var slept []time.Duration
	th := NewThrottler(ThrottleOptions{MaxPerSecond: 2, BurstSize: 2})
	th.sleep = func(d time.Duration) { slept = append(slept, d) }

	input := makeThrottlerResults(5)
	out, err := th.Throttle(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 5 {
		t.Errorf("expected 5 results, got %d", len(out))
	}
	if len(slept) == 0 {
		t.Error("expected at least one sleep call")
	}
}

func TestThrottler_PreservesOrder(t *testing.T) {
	th := NewThrottler(ThrottleOptions{MaxPerSecond: 10, BurstSize: 3})
	th.sleep = func(time.Duration) {}

	input := makeThrottlerResults(6)
	out, err := th.Throttle(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, r := range out {
		if r.ResourceID != input[i].ResourceID {
			t.Errorf("order mismatch at index %d: got %s, want %s", i, r.ResourceID, input[i].ResourceID)
		}
	}
}

func TestThrottler_SleepInterval_MatchesRate(t *testing.T) {
	var slept []time.Duration
	th := NewThrottler(ThrottleOptions{MaxPerSecond: 4, BurstSize: 1})
	th.sleep = func(d time.Duration) { slept = append(slept, d) }

	input := makeThrottlerResults(3)
	_, _ = th.Throttle(input)

	expected := time.Second / 4
	for _, d := range slept {
		if d != expected {
			t.Errorf("expected sleep %v, got %v", expected, d)
		}
	}
}
