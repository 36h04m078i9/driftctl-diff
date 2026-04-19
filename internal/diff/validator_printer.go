package diff

import (
	"fmt"
	"io"
	"os"
)

// ValidatorPrinter renders a ValidationResult in a human-readable format.
type ValidatorPrinter struct {
	w io.Writer
}

// NewValidatorPrinter creates a ValidatorPrinter writing to w, defaulting to stdout.
func NewValidatorPrinter(w io.Writer) *ValidatorPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &ValidatorPrinter{w: w}
}

// Print writes the validation result to the underlying writer.
func (p *ValidatorPrinter) Print(vr ValidationResult) {
	if vr.Valid {
		fmt.Fprintln(p.w, "validation passed: all drift results are structurally sound")
	} else {
		fmt.Fprintln(p.w, "validation failed:")
		for _, e := range vr.Errors {
			fmt.Fprintf(p.w, "  ERROR   %s\n", e)
		}
	}
	for _, w := range vr.Warnings {
		fmt.Fprintf(p.w, "  WARNING %s\n", w)
	}
}
