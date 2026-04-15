// Package provider abstracts live cloud resource fetching behind a common
// interface so that the drift detector remains provider-agnostic.
//
// Usage:
//
//	reg := provider.NewRegistry()
//	awsProvider, err := provider.NewAWSProvider(ctx, "us-east-1")
//	if err != nil { ... }
//	reg.Register("aws", awsProvider)
//
//	attrs, err := reg.FetchAttributes(ctx, "aws", "aws_s3_bucket", "my-bucket")
//
// Supported resource types per provider:
//
//	AWS:
//	  - aws_s3_bucket
package provider
