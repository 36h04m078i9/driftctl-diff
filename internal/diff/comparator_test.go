package diff

import (
	"testing"

	"github.com/acme/driftctl-diff/internal/drift"
)

func makeResourceDiff(rtype, rid string) drift.ResourceDiff {
	return drift.ResourceDiff{
		ResourceType: rtype,
		ResourceID:   rid,
		Changes: []drift.AttributeChange{
			{Attribute: "tags", Got: "a", Want: "b"},
		},
	}
}

func TestCompare_AllNew(t *testing.T) {
	c := NewComparator(CompareOptions{})
	current := []drift.ResourceDiff{makeResourceDiff("aws_s3_bucket", "my-bucket")}
	res := c.Compare(nil, current)
	if len(res.NewDrift) != 1 {
		t.Fatalf("expected 1 new drift, got %d", len(res.NewDrift))
	}
	if len(res.ResolvedDrift) != 0 {
		t.Fatalf("expected 0 resolved, got %d", len(res.ResolvedDrift))
	}
}

func TestCompare_AllResolved(t *testing.T) {
	c := NewComparator(CompareOptions{})
	baseline := []drift.ResourceDiff{makeResourceDiff("aws_s3_bucket", "my-bucket")}
	res := c.Compare(baseline, nil)
	if len(res.ResolvedDrift) != 1 {
		t.Fatalf("expected 1 resolved drift, got %d", len(res.ResolvedDrift))
	}
	if len(res.NewDrift) != 0 {
		t.Fatalf("expected 0 new, got %d", len(res.NewDrift))
	}
}

func TestCompare_Persisted(t *testing.T) {
	c := NewComparator(CompareOptions{})
	d := makeResourceDiff("aws_instance", "i-123")
	res := c.Compare([]drift.ResourceDiff{d}, []drift.ResourceDiff{d})
	if len(res.Persisted) != 1 {
		t.Fatalf("expected 1 persisted, got %d", len(res.Persisted))
	}
	if len(res.NewDrift) != 0 || len(res.ResolvedDrift) != 0 {
		t.Fatal("expected no new or resolved drift")
	}
}

func TestCompare_IgnoreAdded(t *testing.T) {
	c := NewComparator(CompareOptions{IgnoreAdded: true})
	current := []drift.ResourceDiff{makeResourceDiff("aws_s3_bucket", "new-bucket")}
	res := c.Compare(nil, current)
	if len(res.NewDrift) != 0 {
		t.Fatalf("expected new drift suppressed, got %d", len(res.NewDrift))
	}
}

func TestCompare_IgnoreRemoved(t *testing.T) {
	c := NewComparator(CompareOptions{IgnoreRemoved: true})
	baseline := []drift.ResourceDiff{makeResourceDiff("aws_s3_bucket", "old-bucket")}
	res := c.Compare(baseline, nil)
	if len(res.ResolvedDrift) != 0 {
		t.Fatalf("expected resolved drift suppressed, got %d", len(res.ResolvedDrift))
	}
}

func TestCompare_MixedResults(t *testing.T) {
	c := NewComparator(CompareOptions{})
	baseline := []drift.ResourceDiff{
		makeResourceDiff("aws_s3_bucket", "old"),
		makeResourceDiff("aws_instance", "shared"),
	}
	current := []drift.ResourceDiff{
		makeResourceDiff("aws_instance", "shared"),
		makeResourceDiff("aws_s3_bucket", "new"),
	}
	res := c.Compare(baseline, current)
	if len(res.NewDrift) != 1 {
		t.Fatalf("expected 1 new, got %d", len(res.NewDrift))
	}
	if len(res.ResolvedDrift) != 1 {
		t.Fatalf("expected 1 resolved, got %d", len(res.ResolvedDrift))
	}
	if len(res.Persisted) != 1 {
		t.Fatalf("expected 1 persisted, got %d", len(res.Persisted))
	}
}
