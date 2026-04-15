package reporter_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/acme/driftctl-diff/internal/drift"
	"github.com/acme/driftctl-diff/internal/reporter"
	"github.com/acme/driftctl-diff/internal/summary"
)

func sampleChanges() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceID:   "aws_instance.web",
			ResourceType: "aws_instance",
			Changes: []drift.AttributeChange{
				{Attribute: "instance_type", StateValue: "t2.micro", LiveValue: "t3.small"},
			},
		},
	}
}

func TestBuild_SetsGeneratedAt(t *testing.T) {
	before := time.Now().UTC()
	report := reporter.Build(sampleChanges(), summary.Result{})
	after := time.Now().UTC()

	if report.GeneratedAt.Before(before) || report.GeneratedAt.After(after) {
		t.Errorf("GeneratedAt %v not in expected range [%v, %v]", report.GeneratedAt, before, after)
	}
}

func TestBuild_PropagatesChanges(t *testing.T) {
	changes := sampleChanges()
	report := reporter.Build(changes, summary.Result{})

	if len(report.Changes) != len(changes) {
		t.Errorf("expected %d changes, got %d", len(changes), len(report.Changes))
	}
}

func TestWriteJSON_ProducesValidJSON(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)
	report := reporter.Build(sampleChanges(), summary.Result{Total: 1, Drifted: 1})

	if err := r.WriteJSON(report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded reporter.Report
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}

func TestWriteJSON_ContainsSummaryFields(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)
	report := reporter.Build(nil, summary.Result{Total: 5, Drifted: 2, Clean: 3})

	_ = r.WriteJSON(report)

	out := buf.String()
	for _, want := range []string{`"total"`, `"drifted"`, `"clean"`} {
		if !bytes.Contains(buf.Bytes(), []byte(want)) {
			t.Errorf("output missing field %s; got:\n%s", want, out)
		}
	}
}

func TestNew_NilDest_DefaultsToStdout(t *testing.T) {
	// Just ensure no panic when dest is nil.
	r := reporter.New(nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
