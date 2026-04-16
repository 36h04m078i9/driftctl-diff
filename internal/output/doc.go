// Package output provides formatters for rendering infrastructure drift results
// in a variety of formats: human-readable text, table, JSON, YAML, CSV,
// Markdown, XML, JUnit XML, and Go templates.
//
// Each formatter accepts a slice of drift.ResourceDiff values and writes the
// formatted output to an io.Writer. When a nil writer is supplied the
// formatter falls back to os.Stdout.
//
// Available formatters:
//   - NewFormatter        – plain-text diff
//   - NewTableFormatter   – ASCII table
//   - NewJSONFormatter    – JSON
//   - NewYAMLFormatter    – YAML
//   - NewCSVFormatter     – CSV
//   - NewMarkdownFormatter – Markdown table
//   - NewXMLFormatter     – XML
//   - NewJUnitFormatter   – JUnit XML (CI-friendly)
//   - NewTemplateFormatter – custom Go template
package output
