package diff

import (
	"fmt"
	"strings"

	"github.com/acme/driftctl-diff/internal/drift"
)

// Highlighter produces side-by-side or inline diff strings for attribute changes.
type Highlighter struct {
	colorize bool
}

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

// NewHighlighter returns a Highlighter. When colorize is false all ANSI codes
// are omitted so output is safe for pipes and log files.
func NewHighlighter(colorize bool) *Highlighter {
	return &Highlighter{colorize: colorize}
}

// Highlight formats a single AttributeChange as a two-line diff snippet:
//
//	- old_value
//	+ new_value
func (h *Highlighter) Highlight(c drift.AttributeChange) string {
	var sb strings.Builder
	old := fmt.Sprintf("- %s: %s", c.Attribute, c.OldValue)
	new_ := fmt.Sprintf("+ %s: %s", c.Attribute, c.NewValue)
	if h.colorize {
		old = colorRed + old + colorReset
		new_ = colorGreen + new_ + colorReset
	}
	sb.WriteString(old)
	sb.WriteByte('\n')
	sb.WriteString(new_)
	return sb.String()
}

// HighlightAll formats every change in a DriftResult as a labelled block.
func (h *Highlighter) HighlightAll(r drift.DriftResult) string {
	var sb strings.Builder
	header := fmt.Sprintf("# %s / %s", r.ResourceType, r.ResourceID)
	sb.WriteString(header)
	sb.WriteByte('\n')
	for _, ch := range r.Changes {
		sb.WriteString(h.Highlight(ch))
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n")
}
