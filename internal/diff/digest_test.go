package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/acme/driftctl-diff/internal/diff"
	"github.com/acme/driftctl-diff/internal/drift"
)

func makeDigestResults() []drift.DriftResult {
	return []drift.DriftResult{
		{
			ResourceType: "aws_s3_bucket",
			ResourceID:   "my-bucket",
			Changes: []drift.Change{
				{Attribute: "acl", Expected: "private", Actual: "public-read", Kind: drift.KindChanged},
			},
		},
	}
}

func TestDigest_EmptyResults_IsStable(t *testing.T) {
	d := diff.NewDigest(diff.DefaultDigestOptions())
	h1 := d.Compute(nil)
	h2 := d.Compute([]drift.DriftResult{})
	if h1 != h2 {
		t.Fatalf("expected same digest for empty input, got %q vs %q", h1, h2)
	}
}

func TestDigest_SameInputProducesSameHash(t *testing.T) {
	d := diff.NewDigest(diff.DefaultDigestOptions())
	results := makeDigestResults()
	if d.Compute(results) != d.Compute(results) {
		t.Fatal("digest is not deterministic")
	}
}

func TestDigest_DifferentInputProducesDifferentHash(t *testing.T) {
	d := diff.NewDigest(diff.DefaultDigestOptions())
	a := makeDigestResults()
	b := makeDigestResults()
	b[0].Changes[0].Actual = "private" // no longer drifted
	if d.Compute(a) == d.Compute(b) {
		t.Fatal("expected different digests for different inputs")
	}
}

func TestDigest_OrderIndependent(t *testing.T) {
	d := diff.NewDigest(diff.DefaultDigestOptions())
	r1 := makeDigestResults()
	r2 := []drift.DriftResult{r1[0]}
	r1 = append(r1, drift.DriftResult{
		ResourceType: "aws_instance",
		ResourceID:   "i-abc",
		Changes: []drift.Change{
			{Attribute: "instance_type", Expected: "t2.micro", Actual: "t3.micro", Kind: drift.KindChanged},
		},
	})
	r2 = append([]drift.DriftResult{r1[1]}, r2...)
	if d.Compute(r1) != d.Compute(r2) {
		t.Fatal("digest should be order-independent")
	}
}

func TestDigest_IncludeKindFalse_IgnoresKind(t *testing.T) {
	opts := diff.DigestOptions{IncludeKind: false}
	d := diff.NewDigest(opts)
	a := makeDigestResults()
	b := makeDigestResults()
	b[0].Changes[0].Kind = drift.KindAdded
	if d.Compute(a) != d.Compute(b) {
		t.Fatal("expected same digest when IncludeKind is false")
	}
}

func TestDigestPrinter_ContainsHash(t *testing.T) {
	d := diff.NewDigest(diff.DefaultDigestOptions())
	var buf bytes.Buffer
	p := diff.NewDigestPrinter(d, &buf)
	if err := p.Print(makeDigestResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "drift-digest:") {
		t.Errorf("output missing prefix: %q", buf.String())
	}
}

func TestDigestPrinter_NilWriter_DefaultsToStdout(t *testing.T) {
	d := diff.NewDigest(diff.DefaultDigestOptions())
	// Should not panic when writer is nil.
	p := diff.NewDigestPrinter(d, nil)
	if p == nil {
		t.Fatal("expected non-nil printer")
	}
}
