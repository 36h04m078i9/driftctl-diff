package output

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/acme/driftctl-diff/internal/drift"
)

// DotFormatter renders drift results as a Graphviz DOT graph.
type DotFormatter struct {
	w io.Writer
}

// NewDotFormatter creates a DotFormatter writing to w (defaults to stdout).
func NewDotFormatter(w io.Writer) *DotFormatter {
	if w == nil {
		w = os.Stdout
	}
	return &DotFormatter{w: w}
}

// Format writes a DOT graph of the drift results.
func (f *DotFormatter) Format(results []drift.ResourceDiff) error {
	fmt.Fprintln(f.w, "digraph drift {")
	fmt.Fprintln(f.w, "\tgraph [rankdir=LR];")
	fmt.Fprintln(f.w, "\tnode [shape=box];")

	if len(results) == 0 {
		fmt.Fprintln(f.w, "\tno_drift [label=\"No drift detected\" shape=ellipse];")
		fmt.Fprintln(f.w, "}")
		return nil
	}

	for _, r := range results {
		nodeID := sanitize(r.ResourceID)
		label := fmt.Sprintf("%s\\n(%s)", r.ResourceID, r.ResourceType)
		fmt.Fprintf(f.w, "\t%s [label=%q]\n", nodeID, label)

		for _, ch := range r.Changes {
			attID := sanitize(r.ResourceID + "_" + ch.Attribute)
			attLabel := fmt.Sprintf("%s\\nwant: %v\\ngot: %v", ch.Attribute, ch.Expected, ch.Actual)
			fmt.Fprintf(f.w, "\t%s [label=%q shape=ellipse style=filled fillcolor=lightyellow]\n", attID, attLabel)
			fmt.Fprintf(f.w, "\t%s -> %s\n", nodeID, attID)
		}
	}

	fmt.Fprintln(f.w, "}")
	return nil
}

func sanitize(s string) string {
	replacer := strings.NewReplacer("-", "_", ".", "_", "/", "_", ":", "_")
	return replacer.Replace(s)
}
