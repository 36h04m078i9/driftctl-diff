package diff

import (
	"fmt"
	"strings"

	"github.com/acme/driftctl-diff/internal/drift"
)

// Annotation holds a human-readable note attached to a single attribute change.
type Annotation struct {
	ResourceType string
	ResourceID   string
	Attribute    string
	Kind         drift.ChangeKind
	Note         string
}

// Annotator enriches drift results with contextual notes.
type Annotator struct {
	rules []annotationRule
}

type annotationRule struct {
	attribute string
	kind      drift.ChangeKind
	note      string
}

// NewAnnotator returns an Annotator pre-loaded with default rules.
func NewAnnotator() *Annotator {
	return &Annotator{
		rules: []annotationRule{
			{attribute: "tags", kind: drift.KindChanged, note: "Tag drift is common after manual console edits."},
			{attribute: "instance_type", kind: drift.KindChanged, note: "Instance type change may cause downtime."},
			{attribute: "ami", kind: drift.KindChanged, note: "AMI change requires instance replacement."},
			{attribute: "security_groups", kind: drift.KindChanged, note: "Security group drift may expose or restrict traffic unexpectedly."},
		},
	}
}

// Annotate returns annotations for the given drift results.
func (a *Annotator) Annotate(results []drift.ResourceDiff) []Annotation {
	var out []Annotation
	for _, r := range results {
		for _, ch := range r.Changes {
			note := a.matchNote(ch.Attribute, ch.Kind)
			if note == "" {
				note = fmt.Sprintf("Attribute %q has drifted from state.", ch.Attribute)
			}
			out = append(out, Annotation{
				ResourceType: r.ResourceType,
				ResourceID:   r.ResourceID,
				Attribute:    ch.Attribute,
				Kind:         ch.Kind,
				Note:         note,
			})
		}
	}
	return out
}

func (a *Annotator) matchNote(attr string, kind drift.ChangeKind) string {
	for _, r := range a.rules {
		if strings.EqualFold(r.attribute, attr) && r.kind == kind {
			return r.note
		}
	}
	return ""
}
