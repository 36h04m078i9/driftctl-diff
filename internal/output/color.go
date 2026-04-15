package output

import "fmt"

// ANSI color codes for terminal output.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// ColorMode controls whether ANSI color codes are emitted.
type ColorMode int

const (
	// ColorAuto enables color when stdout is a TTY.
	ColorAuto ColorMode = iota
	// ColorAlways forces color output.
	ColorAlways
	// ColorNever disables color output.
	ColorNever
)

// Colorizer wraps strings with ANSI escape sequences.
type Colorizer struct {
	enabled bool
}

// NewColorizer creates a Colorizer. enabled should be derived from ColorMode
// and whether the output file descriptor is a terminal.
func NewColorizer(enabled bool) *Colorizer {
	return &Colorizer{enabled: enabled}
}

func (c *Colorizer) apply(code, s string) string {
	if !c.enabled {
		return s
	}
	return fmt.Sprintf("%s%s%s", code, s, colorReset)
}

// Red renders s in red (used for removals / missing attributes).
func (c *Colorizer) Red(s string) string { return c.apply(colorRed, s) }

// Green renders s in green (used for additions / live values).
func (c *Colorizer) Green(s string) string { return c.apply(colorGreen, s) }

// Yellow renders s in yellow (used for warnings).
func (c *Colorizer) Yellow(s string) string { return c.apply(colorYellow, s) }

// Cyan renders s in cyan (used for resource headers).
func (c *Colorizer) Cyan(s string) string { return c.apply(colorCyan, s) }

// Bold renders s in bold.
func (c *Colorizer) Bold(s string) string { return c.apply(colorBold, s) }
