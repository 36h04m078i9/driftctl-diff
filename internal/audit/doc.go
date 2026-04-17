// Package audit provides append-only audit logging for driftctl-diff scan runs.
//
// Each scan can record an Entry containing the state file path, the number of
// total and drifted resources, and the number of drifted attributes. Entries
// are written as newline-delimited JSON so they can be easily parsed by
// downstream tooling (e.g. jq, Splunk, CloudWatch Logs).
//
// Usage:
//
//	l, err := audit.New("/var/log/driftctl-diff/audit.jsonl")
//	if err != nil { ... }
//	l.Record(audit.Entry{
//		StateFile:  "terraform.tfstate",
//		Drifted:    3,
//		Total:      12,
//		Attributes: 5,
//	})
package audit
