package output

import (
	"strings"
	"testing"
)

func TestColorizer_Disabled(t *testing.T) {
	c := NewColorizer(false)

	cases := []struct {
		name string
		fn   func(string) string
		input string
	}{
		{"Red", c.Red, "hello"},
		{"Green", c.Green, "world"},
		{"Yellow", c.Yellow, "warn"},
		{"Cyan", c.Cyan, "info"},
		{"Bold", c.Bold, "strong"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.fn(tc.input)
			if got != tc.input {
				t.Errorf("expected %q, got %q", tc.input, got)
			}
		})
	}
}

func TestColorizer_Enabled(t *testing.T) {
	c := NewColorizer(true)

	cases := []struct {
		name  string
		fn    func(string) string
		code  string
	}{
		{"Red", c.Red, colorRed},
		{"Green", c.Green, colorGreen},
		{"Yellow", c.Yellow, colorYellow},
		{"Cyan", c.Cyan, colorCyan},
		{"Bold", c.Bold, colorBold},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.fn("text")
			if !strings.HasPrefix(got, tc.code) {
				t.Errorf("expected output to start with ANSI code %q, got %q", tc.code, got)
			}
			if !strings.HasSuffix(got, colorReset) {
				t.Errorf("expected output to end with reset code, got %q", got)
			}
			if !strings.Contains(got, "text") {
				t.Errorf("expected output to contain original text, got %q", got)
			}
		})
	}
}

func TestColorizer_EnabledDoesNotMutateInput(t *testing.T) {
	c := NewColorizer(true)
	input := "unchanged"
	_ = c.Red(input)
	if input != "unchanged" {
		t.Errorf("input string was mutated")
	}
}
