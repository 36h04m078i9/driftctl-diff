package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/owner/driftctl-diff/internal/drift"
)

func makeSummarizerResults() []drift.ResourceDiff {
	return []drift.ResourceDiff{
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-001",
			Changes: []drift.AttributeChange{
				{Attribute: "ami", StateValue: "ami-old", LiveValue: "ami-new"},
				{Attribute: "instance_type", StateValue: "t2.micro", LiveValue: "t3.micro"},
			},
		},
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-002",
			Changes: []drift.AttributeChange{
				{Attribute: "ami", StateValue: "ami-old", LiveValue: "ami-new"},
			},
		},
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "acl", StateValue: "private", LiveValue: "public-read"},
				{Attribute: "versioning", StateValue: "false", LiveValue: "true"},
				{Attribute: "region", StateValue: "us-east-1", LiveValue: "us-west-2"},
			},
		},
	}
}

func TestSummarizer_GroupsByType(t *testing.T) {
	s := NewSummarizer(SummaryOptions{})
	results := makeSummarizerResults()
	summaries := s.Summarize(results)
	if len(summaries) != 2 {
		t.Fatalf("expected 2 type summaries, got %d", len(summaries))
	}
}

func TestSummarizer_SortedByChangedAttrsDesc(t *testing.T) {
	s := NewSummarizer(SummaryOptions{})
	summaries := s.Summarize(makeSummarizerResults())
	// aws_s3_bucket has 3 changed attrs; aws_instance has 3 total but s3 still tops
	if summaries[0].ResourceType != "aws_s3_bucket" {
		t.Errorf("expected aws_s3_bucket first, got %s", summaries[0].ResourceType)
	}
}

func TestSummarizer_TopN_LimitsResults(t *testing.T) {
	s := NewSummarizer(SummaryOptions{TopN: 1})
	summaries := s.Summarize(makeSummarizerResults())
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary with TopN=1, got %d", len(summaries))
	}
}

func TestSummarizer_ResourceCountIsCorrect(t *testing.T) {
	s := NewSummarizer(SummaryOptions{})
	summaries := s.Summarize(makeSummarizerResults())
	for _, ts := range summaries {
		if ts.ResourceType == "aws_instance" && ts.ResourceCount != 2 {
			t.Errorf("expected aws_instance count 2, got %d", ts.ResourceCount)
		}
	}
}

func TestSummarizer_Print_NoDrift(t *testing.T) {
	s := NewSummarizer(SummaryOptions{})
	var buf bytes.Buffer
	s.Print(&buf, []TypeSummary{})
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestSummarizer_Print_ContainsHeaders(t *testing.T) {
	s := NewSummarizer(SummaryOptions{})
	summaries := s.Summarize(makeSummarizerResults())
	var buf bytes.Buffer
	s.Print(&buf, summaries)
	out := buf.String()
	if !strings.Contains(out, "RESOURCE TYPE") {
		t.Errorf("expected RESOURCE TYPE header, got: %s", out)
	}
	if !strings.Contains(out, "aws_s3_bucket") {
		t.Errorf("expected aws_s3_bucket in output, got: %s", out)
	}
}
