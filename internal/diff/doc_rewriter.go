// Package diff provides utilities for processing and presenting infrastructure
// drift results produced by driftctl-diff.
//
// # Rewriter
//
// The Rewriter applies user-defined string-substitution rules to the attribute
// values of drift results. This is useful when environment-specific tokens
// (e.g. account IDs, region aliases, or temporary suffixes) should be
// normalised before the diff is displayed or exported.
//
// Rules are matched by ResourceType and Attribute name. Both fields support a
// trailing wildcard ("*") to match any value with a given prefix, or the
// literal "*" to match everything.
//
// Example:
//
//	opts := diff.RewriteOptions{
//		Rules: []diff.RewriteRule{
//			{ResourceType: "aws_s3*", Attribute: "*", Find: "staging", Replace: "production"},
//		},
//	}
//	rewriter := diff.NewRewriter(opts)
//	rewritten := rewriter.Rewrite(results)
//
// The Rewriter never mutates its input; it always returns a freshly allocated
// slice of DriftResults.
package diff
