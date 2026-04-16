package output_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
	"github.com/acme/driftctl-diff/internal/output"
)

func TestJSONFormatter_NoDrift(t *testing.T) {
	f := output.NewJSONFormatter()
	var buf bytes.Buffer
	if err := f.Format(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["drifted"].(bool) {
		t.Error("expected drifted=false for empty results")
	}
}

func TestJSONFormatter_WithChanges_DriftedTrue(t *testing.T) {
	f := output.NewJSONFormatter()
	results := []drift.ResourceDiff{
		{
			ResourceID:   "i-abc123",
			ResourceType: "aws_instance",
			Changes: []drift.AttributeChange{
				{Attribute: "instance_type", Kind: drift.KindChanged, Want: "t2.micro", Got: "t3.small"},
			},
		},
	}
	var buf bytes.Buffer
	if err := f.Format(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !out["drifted"].(bool) {
		t.Error("expected drifted=true")
	}
	resources := out["resources"].([]interface{})
	if len(resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(resources))
	}
	res := resources[0].(map[string]interface{})
	if res["id"] != "i-abc123" {
		t.Errorf("unexpected id: %v", res["id"])
	}
}

func TestJSONFormatter_ContainsChangeFields(t *testing.T) {
	f := output.NewJSONFormatter()
	results := []drift.ResourceDiff{
		{
			ResourceID:   "sg-001",
			ResourceType: "aws_security_group",
			Changes: []drift.AttributeChange{
				{Attribute: "description", Kind: drift.KindMissing, Want: "managed", Got: ""},
			},
		},
	}
	var buf bytes.Buffer
	_ = f.Format(&buf, results)

	if !bytes.Contains(buf.Bytes(), []byte("description")) {
		t.Error("expected attribute name in output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("managed")) {
		t.Error("expected want value in output")
	}
}
