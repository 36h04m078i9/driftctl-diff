package diff

import (
	"testing"

	"github.com/driftctl-diff/internal/drift"
)

func makeTaggerResults() []drift.DriftResult {
	return []drift.DriftResult{
		{ResourceID: "bucket-1", ResourceType: "aws_s3_bucket"},
		{ResourceID: "sg-2", ResourceType: "aws_security_group"},
	}
}

func TestTagger_AttachesEnvTag(t *testing.T) {
	tagger := NewTagger(TaggerOptions{EnvTag: "staging"})
	tagged := tagger.Tag(makeTaggerResults())
	if len(tagged) != 2 {
		t.Fatalf("expected 2 results, got %d", len(tagged))
	}
	if tagged[0].Tags[0].Key != "env" || tagged[0].Tags[0].Value != "staging" {
		t.Errorf("expected env=staging, got %+v", tagged[0].Tags)
	}
}

func TestTagger_AttachesRegionTag(t *testing.T) {
	tagger := NewTagger(TaggerOptions{RegionTag: "eu-west-1"})
	tagged := tagger.Tag(makeTaggerResults())
	found := false
	for _, tag := range tagged[0].Tags {
		if tag.Key == "region" && tag.Value == "eu-west-1" {
			found = true
		}
	}
	if !found {
		t.Error("region tag not found")
	}
}

func TestTagger_AttachesCustomTags(t *testing.T) {
	tagger := NewTagger(TaggerOptions{CustomTags: map[string]string{"team": "platform"}})
	tagged := tagger.Tag(makeTaggerResults())
	found := false
	for _, tag := range tagged[0].Tags {
		if tag.Key == "team" && tag.Value == "platform" {
			found = true
		}
	}
	if !found {
		t.Error("custom tag not found")
	}
}

func TestTagger_NoOptions_NoTags(t *testing.T) {
	tagger := NewTagger(TaggerOptions{})
	tagged := tagger.Tag(makeTaggerResults())
	for _, r := range tagged {
		if len(r.Tags) != 0 {
			t.Errorf("expected no tags, got %v", r.Tags)
		}
	}
}

func TestFilterByTag_ReturnsMatching(t *testing.T) {
	tagger := NewTagger(TaggerOptions{EnvTag: "prod"})
	tagged := tagger.Tag(makeTaggerResults())
	filtered := FilterByTag(tagged, "env", "prod")
	if len(filtered) != 2 {
		t.Errorf("expected 2, got %d", len(filtered))
	}
}

func TestFilterByTag_NoMatch_ReturnsEmpty(t *testing.T) {
	tagger := NewTagger(TaggerOptions{EnvTag: "prod"})
	tagged := tagger.Tag(makeTaggerResults())
	filtered := FilterByTag(tagged, "env", "staging")
	if len(filtered) != 0 {
		t.Errorf("expected 0, got %d", len(filtered))
	}
}
