// Package diff provides utilities for working with drift results, including
// pagination, searching, sorting, highlighting, and grouping.
//
// Grouper partitions []drift.Result by resource type or change kind, making
// it easy to render grouped output in CLI formatters or reports.
//
// Example usage:
//
//	g := diff.NewGrouper()
//	groups := g.GroupByType(results)
//	for _, grp := range groups {
//		fmt.Println(grp.ResourceType, len(grp.Results))
//	}
package diff
