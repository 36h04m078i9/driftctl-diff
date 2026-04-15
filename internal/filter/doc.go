// Package filter implements include/exclude rules for scoping drift detection
// to a specific subset of Terraform-managed resources.
//
// Rules are matched by resource type and an optional resource ID pattern.
// Resource IDs support a trailing "*" wildcard to match by prefix, e.g.
// "prod-*" matches "prod-api", "prod-db", etc.
//
// Exclude rules always take precedence over include rules. When no include
// rules are registered, all resources not matched by an exclude rule are
// considered allowed.
//
// Example usage:
//
//	f := filter.New()
//	f.AddInclude("aws_s3_bucket", "*")         // only S3 buckets
//	f.AddExclude("aws_s3_bucket", "logs-*")    // except log buckets
//	f.Allow("aws_s3_bucket", "my-app")         // → true
//	f.Allow("aws_s3_bucket", "logs-archive")   // → false
package filter
