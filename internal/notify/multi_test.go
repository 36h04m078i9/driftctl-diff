package notify_test

import (
	"errors"
	"testing"

	"github.com/acme/driftctl-diff/internal/notify"
)

type stubNotifier struct {
	called bool
	err    error
}

func (s *stubNotifier) Notify(_ notify.Event) error {
	s.called = true
	return s.err
}

// sampleEvent returns a minimal Event for use in tests.
func sampleEvent() notify.Event {
	return notify.Event{}
}

func TestMultiNotifier_CallsAll(t *testing.T) {
	a, b := &stubNotifier{}, &stubNotifier{}
	m := notify.NewMultiNotifier(a, b)
	if err := m.Notify(sampleEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a.called || !b.called {
		t.Error("expected both notifiers to be called")
	}
}

func TestMultiNotifier_CollectsErrors(t *testing.T) {
	a := &stubNotifier{err: errors.New("fail a")}
	b := &stubNotifier{err: errors.New("fail b")}
	m := notify.NewMultiNotifier(a, b)
	if err := m.Notify(sampleEvent()); err == nil {
		t.Fatal("expected combined error")
	}
}

func TestMultiNotifier_PartialError_StillCallsRest(t *testing.T) {
	a := &stubNotifier{err: errors.New("fail")}
	b := &stubNotifier{}
	m := notify.NewMultiNotifier(a, b)
	m.Notify(sampleEvent()) //nolint:errcheck
	if !b.called {
		t.Error("second notifier should still be called after first fails")
	}
}

func TestMultiNotifier_Empty_NoError(t *testing.T) {
	m := notify.NewMultiNotifier()
	if err := m.Notify(sampleEvent()); err != nil {
		t.Fatalf("expected no error from empty notifier, got: %v", err)
	}
}
