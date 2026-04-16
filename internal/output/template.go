package output

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/owner/driftctl-diff/internal/drift"
)

const defaultTemplate = `{{- if not .Drifted}}No drift detected.
{{- else}}Drift detected in {{ len .Changes }} resource(s):
{{ range .Changes}}
Resource: {{ .ResourceType }}/{{ .ResourceID }}
{{ range .Attributes}}  {{ .Attribute }}: {{ .StateValue }} => {{ .LiveValue }}
{{ end}}{{- end}}`

// TemplateFormatter renders drift results using a Go text/template.
type TemplateFormatter struct {
	tmpl *template.Template
	out  io.Writer
}

type templateData struct {
	Drifted bool
	Changes []drift.ResourceDiff
}

// NewTemplateFormatter creates a TemplateFormatter with an optional custom
// template string. Pass an empty string to use the built-in default.
func NewTemplateFormatter(tmplStr string, out io.Writer) (*TemplateFormatter, error) {
	if tmplStr == "" {
		tmplStr = defaultTemplate
	}
	if out == nil {
		out = os.Stdout
	}
	t, err := template.New("drift").Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("template parse error: %w", err)
	}
	return &TemplateFormatter{tmpl: t, out: out}, nil
}

// Format executes the template against the provided changes and writes output.
func (f *TemplateFormatter) Format(changes []drift.ResourceDiff) error {
	var buf bytes.Buffer
	data := templateData{
		Drifted: len(changes) > 0,
		Changes: changes,
	}
	if err := f.tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}
	_, err := fmt.Fprint(f.out, buf.String())
	return err
}
