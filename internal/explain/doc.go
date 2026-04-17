// Package explain provides drift explanation logic for driftctl-diff.
//
// It analyses detected drift changes and produces human-readable explanations
// with severity levels (info, warning, critical) based on configurable rules.
//
// Default rules flag sensitive attributes (passwords, policies, ports) at
// elevated severity to help operators prioritise remediation.
//
// Usage:
//
//	e := explain.New()
//	explanations := e.Explain(results)
//	explain.NewPrinter(os.Stdout).Print(explanations)
package explain
