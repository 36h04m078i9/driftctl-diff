// Package notify provides drift notification hooks (stdout, webhook, etc.).
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/acme/driftctl-diff/internal/drift"
)

// Event is the payload sent to a notifier.
type Event struct {
	DriftedCount int                `json:"drifted_count"`
	Changes      []drift.DriftResult `json:"changes"`
	GeneratedAt  time.Time          `json:"generated_at"`
}

// Notifier sends drift events to an external destination.
type Notifier interface {
	Notify(e Event) error
}

// WebhookNotifier posts drift events as JSON to a URL.
type WebhookNotifier struct {
	URL    string
	Client *http.Client
}

// NewWebhookNotifier creates a WebhookNotifier with a default HTTP client.
func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{
		URL:    url,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify marshals the event and POSTs it to the configured URL.
func (w *WebhookNotifier) Notify(e Event) error {
	body, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("notify: marshal: %w", err)
	}
	resp, err := w.Client.Post(w.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notify: unexpected status %d", resp.StatusCode)
	}
	return nil
}
