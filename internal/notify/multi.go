package notify

import "fmt"

// MultiNotifier fans out an Event to multiple Notifiers.
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier creates a MultiNotifier from the provided notifiers.
func NewMultiNotifier(nn ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: nn}
}

// Notify calls every registered Notifier and collects errors.
func (m *MultiNotifier) Notify(e Event) error {
	var errs []error
	for _, n := range m.notifiers {
		if err := n.Notify(e); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("multi-notify errors: %v", errs)
}
