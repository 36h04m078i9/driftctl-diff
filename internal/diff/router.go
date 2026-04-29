package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/acme/driftctl-diff/internal/drift"
)

// RouteOptions controls how the Router dispatches results.
type RouteOptions struct {
	// Routes maps a label value to a writer. Results whose label matches a key
	// are written to the corresponding writer.
	Routes map[string]io.Writer
	// LabelKey is the metadata key used to determine the route (default: "env").
	LabelKey string
	// DefaultWriter receives results that do not match any route.
	DefaultWriter io.Writer
}

// Router dispatches drift results to different writers based on a metadata label.
type Router struct {
	opts RouteOptions
}

// NewRouter creates a Router with the given options.
func NewRouter(opts RouteOptions) *Router {
	if opts.LabelKey == "" {
		opts.LabelKey = "env"
	}
	if opts.DefaultWriter == nil {
		opts.DefaultWriter = os.Stdout
	}
	return &Router{opts: opts}
}

// Route dispatches each DriftResult to the appropriate writer.
// Results are written as a simple text summary line.
func (r *Router) Route(results []drift.DriftResult) error {
	for _, res := range results {
		w := r.writerFor(res)
		line := fmt.Sprintf("resource=%s type=%s changes=%d\n",
			res.ResourceID, res.ResourceType, len(res.Changes))
		if _, err := fmt.Fprint(w, line); err != nil {
			return fmt.Errorf("router: write failed for resource %s: %w", res.ResourceID, err)
		}
	}
	return nil
}

func (r *Router) writerFor(res drift.DriftResult) io.Writer {
	if res.Metadata != nil {
		if val, ok := res.Metadata[r.opts.LabelKey]; ok {
			if w, found := r.opts.Routes[val]; found {
				return w
			}
		}
	}
	return r.opts.DefaultWriter
}
