package output

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/owner/driftctl-diff/internal/drift"
)

// JUnitFormatter writes drift results in JUnit XML format.
type JUnitFormatter struct {
	w io.Writer
}

func NewJUnitFormatter(w io.Writer) *JUnitFormatter {
	if w == nil {
		w = os.Stdout
	}
	return &JUnitFormatter{w: w}
}

type junitSuites struct {
	XMLName xml.Name     `xml:"testsuites"`
	Suites  []junitSuite `xml:"testsuite"`
}

type junitSuite struct {
	Name      string       `xml:"name,attr"`
	Tests     int          `xml:"tests,attr"`
	Failures  int          `xml:"failures,attr"`
	Timestamp string       `xml:"timestamp,attr"`
	Cases     []junitCase  `xml:"testcase"`
}

type junitCase struct {
	Name      string        `xml:"name,attr"`
	Classname string        `xml:"classname,attr"`
	Failure   *junitFailure `xml:"failure,omitempty"`
}

type junitFailure struct {
	Message string `xml:"message,attr"`
	Text    string `xml:",chardata"`
}

func (f *JUnitFormatter) Format(results []drift.ResourceDiff) error {
	cases := make([]junitCase, 0, len(results))
	failures := 0
	for _, r := range results {
		c := junitCase{
			Name:      r.ResourceID,
			Classname: r.ResourceType,
		}
		if len(r.Changes) > 0 {
			failures++
			msg := fmt.Sprintf("%d attribute(s) drifted", len(r.Changes))
			text := ""
			for _, ch := range r.Changes {
				text += fmt.Sprintf("%s: expected %q got %q\n", ch.Attribute, ch.Expected, ch.Actual)
			}
			c.Failure = &junitFailure{Message: msg, Text: text}
		}
		cases = append(cases, c)
	}
	suites := junitSuites{
		Suites: []junitSuite{
			{
				Name:      "driftctl-diff",
				Tests:     len(cases),
				Failures:  failures,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Cases:     cases,
			},
		},
	}
	out, err := xml.MarshalIndent(suites, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(f.w, "%s%s\n", xml.Header, out)
	return err
}
