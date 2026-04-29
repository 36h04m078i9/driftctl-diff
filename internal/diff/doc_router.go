// Package diff provides utilities for processing and presenting infrastructure
// drift results produced by driftctl-diff.
//
// # Router
//
// The Router type dispatches [drift.DriftResult] slices to different
// [io.Writer] destinations based on a metadata label attached to each result.
//
// A common use-case is routing drift results by environment so that prod,
// staging, and dev findings are written to separate files or streams:
//
//	prodFile, _ := os.Create("prod.txt")
//	stagingFile, _ := os.Create("staging.txt")
//
//	router := diff.NewRouter(diff.RouteOptions{
//		LabelKey: "env",
//		Routes: map[string]io.Writer{
//			"prod":    prodFile,
//			"staging": stagingFile,
//		},
//	})
//	if err := router.Route(results); err != nil {
//		log.Fatal(err)
//	}
//
// Results that carry no matching label are sent to RouteOptions.DefaultWriter
// (defaults to os.Stdout).
//
// RouterPrinter can be used to display a summary table of how many resources
// were dispatched to each route.
package diff
