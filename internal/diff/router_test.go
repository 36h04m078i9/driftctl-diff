package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeRouterResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceID:   "vpc-aaa",
			ResourceType: "aws_vpc",
			Metadata:     map[string]string{"env": "prod"},
			Changes: []drift.Change{
				{Attribute: "cidr", StateValue: "10.0.0.0/16", LiveValue: "10.1.0.0/16"},
			},
		},
		{
			ResourceID:   "sg-bbb",
			ResourceType: "aws_security_group",
			Metadata:     map[string]string{"env": "staging"},
			Changes: []drift.Change{
				{Attribute: "name", StateValue: "old", LiveValue: "new"},
			},
		},
		{
			ResourceID:   "s3-ccc",
			ResourceType: "aws_s3_bucket",
			Metadata:     map[string]string{"env": "dev"},
			Changes:      []drift.Change{},
		},
	}
}

func TestRouter_DispatchesToMatchingWriter(t *testing.T) {
	prod := &bytes.Buffer{}
	staging := &bytes.Buffer{}
	router := NewRouter(RouteOptions{
		Routes: map[string]io.Writer{"prod": prod, "staging": staging},
	})
	if err := router.Route(makeRouterResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(prod.String(), "vpc-aaa") {
		t.Errorf("expected prod writer to contain vpc-aaa, got: %s", prod.String())
	}
	if !strings.Contains(staging.String(), "sg-bbb") {
		t.Errorf("expected staging writer to contain sg-bbb, got: %s", staging.String())
	}
}

func TestRouter_UnmatchedGoesToDefault(t *testing.T) {
	defaultBuf := &bytes.Buffer{}
	router := NewRouter(RouteOptions{
		Routes:        map[string]io.Writer{},
		DefaultWriter: defaultBuf,
	})
	if err := router.Route(makeRouterResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(defaultBuf.String(), "vpc-aaa") {
		t.Errorf("expected default writer to receive unmatched result")
	}
}

func TestRouter_NilMetadata_GoesToDefault(t *testing.T) {
	defaultBuf := &bytes.Buffer{}
	router := NewRouter(RouteOptions{DefaultWriter: defaultBuf})
	results := []drift.DriftResult{
		{ResourceID: "rds-xyz", ResourceType: "aws_db_instance", Metadata: nil},
	}
	if err := router.Route(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(defaultBuf.String(), "rds-xyz") {
		t.Errorf("expected nil-metadata result in default writer")
	}
}

func TestRouter_CustomLabelKey(t *testing.T) {
	regionBuf := &bytes.Buffer{}
	router := NewRouter(RouteOptions{
		LabelKey: "region",
		Routes:   map[string]io.Writer{"us-east-1": regionBuf},
	})
	results := []drift.DriftResult{
		{ResourceID: "ec2-1", ResourceType: "aws_instance",
			Metadata: map[string]string{"region": "us-east-1"}},
	}
	if err := router.Route(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(regionBuf.String(), "ec2-1") {
		t.Errorf("expected result routed by region label")
	}
}

func TestRouter_EmptyResults_NoError(t *testing.T) {
	router := NewRouter(RouteOptions{})
	if err := router.Route([]drift.DriftResult{}); err != nil {
		t.Fatalf("expected no error for empty results, got: %v", err)
	}
}
