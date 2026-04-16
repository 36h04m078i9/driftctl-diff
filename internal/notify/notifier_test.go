package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/acme/driftctl-diff/internal/notify"
)

func sampleEvent() notify.Event {
	return notify.Event{
		DriftedCount: 2,
		GeneratedAt:  time.Now(),
	}
}

func TestWebhookNotifier_Success(t *testing.T) {
	var received notify.Event
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := notify.NewWebhookNotifier(srv.URL)
	if err := n.Notify(sampleEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.DriftedCount != 2 {
		t.Errorf("expected drifted_count 2, got %d", received.DriftedCount)
	}
}

func TestWebhookNotifier_Non2xx_ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n := notify.NewWebhookNotifier(srv.URL)
	if err := n.Notify(sampleEvent()); err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestWebhookNotifier_BadURL_ReturnsError(t *testing.T) {
	n := notify.NewWebhookNotifier("http://127.0.0.1:0/nowhere")
	if err := n.Notify(sampleEvent()); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
