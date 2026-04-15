// Package reporter produces structured, machine-readable reports of drift
// detection results.
//
// # Overview
//
// Use [Build] to assemble a [Report] from drift results and a summary, then
// use [Reporter.WriteJSON] to serialise it:
//
//	changes := detector.Detect(stateResources, liveResources)
//	sum     := summary.Compute(changes)
//	report  := reporter.Build(changes, sum)
//
//	r := reporter.New(os.Stdout)
//	if err := r.WriteJSON(report); err != nil {
//		log.Fatal(err)
//	}
//
// The JSON output is stable and suitable for piping into other tools.
package reporter
