package state_test

import (
	"testing"

	"github.com/user/driftctl-diff/internal/state"
)

const validStateJSON = `{
  "version": 4,
  "resources": [
    {
      "type": "aws_s3_bucket",
      "name": "my_bucket",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "attributes": {"bucket": "my-bucket-name", "region": "us-east-1"}
    }
  ]
}`

func TestParse_ValidState(t *testing.T) {
	p := state.NewParser()
	s, err := p.Parse([]byte(validStateJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(s.Resources))
	}
	if s.Resources[0].Type != "aws_s3_bucket" {
		t.Errorf("expected type aws_s3_bucket, got %s", s.Resources[0].Type)
	}
}

func TestParse_UnsupportedVersion(t *testing.T) {
	p := state.NewParser()
	_, err := p.Parse([]byte(`{"version": 3, "resources": []}`))
	if err == nil {
		t.Fatal("expected error for unsupported version, got nil")
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	p := state.NewParser()
	_, err := p.Parse([]byte(`not json`))
	if err == nil {
		t JSON, got nil")
	}
}

func TestParse_EmptyResources(t *testing.T) {
	p := state.NewParser()
	s, err := p.Parse([]byte(`{"version": 4, "resources": []}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Resources) != 0 {
		t.Errorf("expected 0 resources, got %d", len(s.Resources))
	}
}
