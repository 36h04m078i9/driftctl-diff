// Package notify implements pluggable drift notification hooks.
//
// Notifiers receive an Event containing the drift results and can forward
// them to external systems such as webhooks, Slack, or PagerDuty.
//
// Usage:
//
//	wh := notify.NewWebhookNotifier("https://hooks.example.com/drift")
//	err := wh.Notify(notify.Event{
//		DriftedCount: len(changes),
//		Changes:      changes,
//		GeneratedAt:  time.Now(),
//	})
//
// Multiple notifiers can be composed with NewMultiNotifier.
package notify
