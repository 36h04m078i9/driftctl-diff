// Package policy provides severity-based policy evaluation for drift results.
//
// Users define rules in a YAML file mapping resource types (and optional
// attributes) to a severity level (low, medium, high). The Evaluator checks
// each drift.Change against the rule set and returns Violations for any
// matches. MaxSeverity can be used to determine the overall exit code or
// alert level for a drift run.
//
// Example policy file:
//
//	rules:
//	  - resource_type: aws_s3_bucket
//	    attribute: acl
//	    severity: high
//	  - resource_type: "*"
//	    severity: low
package policy
