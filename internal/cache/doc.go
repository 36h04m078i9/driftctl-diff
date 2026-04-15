// Package cache provides an in-memory, thread-safe caching layer for live
// cloud resource attributes fetched during a drift-detection run.
//
// # Overview
//
// Fetching attributes from a cloud provider (e.g. AWS) involves network I/O
// that can be slow and subject to API rate-limits. When the same resource is
// referenced more than once during a single run (e.g. in multiple Terraform
// modules), caching avoids redundant calls.
//
// # Usage
//
// Use cache.New() for a bare key/value store, or wrap any AttributeFetcher
// with cache.NewCachedProvider to get transparent caching:
//
//	  inner := provider.NewAWSProvider(cfg)
//	  cached := cache.NewCachedProvider(inner)
//	  attrs, err := cached.FetchAttributes("aws_s3_bucket", "my-bucket")
//
// The cache is not persisted between runs; it is flushed automatically when
// the process exits.
package cache
