// Package output provides formatters that render drift results in various
// human- and machine-readable formats.
//
// Available formatters:
//
//   - Formatter      – coloured, line-oriented text (default)
//   - TableFormatter – ASCII table
//   - JSONFormatter  – JSON array of resource diffs
//   - MarkdownFormatter – GitHub-flavoured Markdown table
//   - CSVFormatter   – RFC 4180 CSV with one row per attribute change
//
// All formatters accept an io.Writer; passing nil defaults to os.Stdout.
package output
