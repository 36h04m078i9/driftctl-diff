package diff

import (
	"bytes"
	"testing"

	"github.com/snyk/driftctl-diff/internal/drift"
)

func makeScorerResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.AttributeChange{
				{Attribute: "acl", Kind: drift.KindChanged},
				{Attribute: "tags", Kind: drift.KindMissing},
			},
		},
		{
			ResourceType: "aws_iam_role",
			ResourceID:   "my-role",
			Changes: []drift.AttributeChange{
				{Attribute: "assume_role_policy", Kind: drift.KindAdded},
			},
		},
	}
}

func TestScore_SortedDescending(t *testing.T) {
	s := NewScorer()
	scores := s.Score(makeScorerResults())
	if len(scores) != 2 {
		t.Fatalf("expected 2 scores, got %d", len(scores))
	}
	if scores[0].Total <= scores[1].Total {
		t.Errorf("expected descending order, got %d <= %d", scores[0].Total, scores[1].Total)
	}
}

func TestScore_TotalsAreCorrect(t *testing.T) {
	s := NewScorer()
	scores := s.Score(makeScorerResults())
	// bucket: KindChanged(2) + KindMissing(3) = 5
	if scores[0].Total != 5 {
		t.Errorf("expected total 5, got %d", scores[0].Total)
	}
}

func TestScore_Empty(t *testing.T) {
	s := NewScorer()
	scores := s.Score(nil)
	if len(scores) != 0 {
		t.Errorf("expected empty scores")
	}
}

func TestScorerPrinter_NoScores(t *testing.T) {
	var buf bytes.Buffer
	p := NewScorerPrinter(&buf)
	p.Print(nil)
	if !bytes.Contains(buf.Bytes(), []byte("no drift scores")) {
		t.Errorf("expected no drift message, got: %s", buf.String())
	}
}

func TestScorerPrinter_ContainsResourceID(t *testing.T) {
	var buf bytes.Buffer
	p := NewScorerPrinter(&buf)
	s := NewScorer()
	p.Print(s.Score(makeScorerResults()))
	if !bytes.Contains(buf.Bytes(), []byte("my-bucket")) {
		t.Errorf("expected resource ID in output, got: %s", buf.String())
	}
}

func TestScorerPrinter_NilWriter_DefaultsToStdout(t *testing.T) {
	p := NewScorerPrinter(nil)
	if p.w == nil {
		t.Error("expected non-nil writer")
	}
}
