package output

import (
	"encoding/xml"
	"io"
	"os"
	"time"

	"github.com/user/driftctl-diff/internal/drift"
)

// XMLFormatter writes drift results as XML.
type XMLFormatter struct {
	w io.Writer
}

type xmlReport struct {
	XMLName     xml.Name     `xml:"DriftReport"`
	GeneratedAt string       `xml:"GeneratedAt,attr"`
	Drifted     bool         `xml:"Drifted,attr"`
	Resources   []xmlResource `xml:"Resource"`
}

type xmlResource struct {
	ID      string       `xml:"id,attr"`
	Type    string       `xml:"type,attr"`
	Changes []xmlChange  `xml:"Change"`
}

type xmlChange struct {
	Attribute string `xml:"attribute,attr"`
	Kind      string `xml:"kind,attr"`
	Want      string `xml:"want,omitempty"`
	Got       string `xml:"got,omitempty"`
}

// NewXMLFormatter creates an XMLFormatter writing to w; if w is nil stdout is used.
func NewXMLFormatter(w io.Writer) *XMLFormatter {
	if w == nil {
		w = os.Stdout
	}
	return &XMLFormatter{w: w}
}

// Format encodes results as indented XML.
func (f *XMLFormatter) Format(results []drift.ResourceDiff) error {
	report := xmlReport{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Drifted:     len(results) > 0,
	}
	for _, r := range results {
		xr := xmlResource{ID: r.ResourceID, Type: r.ResourceType}
		for _, c := range r.Changes {
			xr.Changes = append(xr.Changes, xmlChange{
				Attribute: c.Attribute,
				Kind:      c.Kind.String(),
				Want:      c.Want,
				Got:       c.Got,
			})
		}
		report.Resources = append(report.Resources, xr)
	}
	enc := xml.NewEncoder(f.w)
	enc.Indent("", "  ")
	if err := enc.Encode(report); err != nil {
		return err
	}
	return enc.Flush()
}
