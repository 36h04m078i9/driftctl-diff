package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeScorerResults() []drift.Result {
	return []drift.Result{
		{
			ResourceType: "aws_instance",
			ResourceID:   "i-001",
			Changes: []drift.Change{
				{Kind: drift.KindChanged, Attribute: "ami"},
				{Kind: drift.KindMissing, Attribute: "tags"},
			},
		},
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.Change{
				{Kind: drift.KindChanged, Attribute: "acl"},
			},
		},
	}
}

func TestScore_SortedDescending(t *testing.T) {
	s := NewScorer(DefaultScorerOptions())
	scores := s.Score(makeScorerResults())
	if len(scores) != 2 {
		t.Fatalf("expected 2 scores, got %d", len(scores))
	}
	if scores[0].Total <= scores[1].Total {
		t.Errorf("expected descending order, got %d <= %d", scores[0].Total, scores[1].Total)
	}
}

func TestScore_TotalsAreCorrect(t *testing.T) {
	opts := ScorerOptions{WeightChanged: 1, WeightMissing: 2}
	s := NewScorer(opts)
	scores := s.Score(makeScorerResults())
	// aws_instance: 1 changed + 1 missing = 1*1 + 1*2 = 3
	var inst Score
	for _, sc := range scores {
		if sc.ResourceID == "i-001" {
			inst = sc
		}
	}
	if inst.Total != 3 {
		t.Errorf("expected total 3, got %d", inst.Total)
	}
}

func TestScore_Empty(t *testing.T) {
	s := NewScorer(DefaultScorerOptions())
	scores := s.Score(nil)
	if len(scores) != 0 {
		t.Errorf("expected empty slice")
	}
}

func TestScorerPrinter_NoScores(t *testing.T) {
	var buf bytes.Buffer
	p := NewScorerPrinter(&buf)
	p.Print(nil)
	if !strings.Contains(buf.String(), "No drift scores") {
		t.Errorf("expected no-scores message, got: %s", buf.String())
	}
}

func TestScorerPrinter_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	p := NewScorerPrinter(&buf)
	s := NewScorer(DefaultScorerOptions())
	p.Print(s.Score(makeScorerResults()))
	out := buf.String()
	for _, hdr := range []string{"RESOURCE TYPE", "RESOURCE ID", "SCORE"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q in output: %s", hdr, out)
		}
	}
}

func TestScorerPrinter_ContainsResourceInfo(t *testing.T) {
	var buf bytes.Buffer
	p := NewScorerPrinter(&buf)
	s := NewScorer(DefaultScorerOptions())
	p.Print(s.Score(makeScorerResults()))
	out := buf.String()
	if !strings.Contains(out, "aws_instance") {
		t.Errorf("expected resource type in output: %s", out)
	}
	if !strings.Contains(out, "i-001") {
		t.Errorf("expected resource ID in output: %s", out)
	}
}
