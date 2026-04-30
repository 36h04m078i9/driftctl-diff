// Package diff provides utilities for processing and presenting infrastructure
// drift results produced by driftctl-diff.
//
// # Linker
//
// The Linker finds relationships between drifted resources by inspecting shared
// attribute values. For example, two EC2 instances that both drifted on their
// subnet_id attribute and share the same live value are considered linked.
//
// Usage:
//
//	opts := diff.DefaultLinkerOptions()
//	opts.LinkByAttribute = "subnet_id"
//	linker := diff.NewLinker(opts)
//	result := linker.Link(driftResults)
//
//	printer := diff.NewLinkerPrinter(os.Stdout)
//	printer.Print(result)
package diff
