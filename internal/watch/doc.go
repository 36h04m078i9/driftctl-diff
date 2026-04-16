// Package watch implements continuous drift monitoring by periodically
// re-running detection and surfacing changes through a result channel.
//
// Basic usage:
//
//	w := watch.New(myRunner, 5*time.Minute)
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	go w.Start(ctx)
//	for r := range w.Results() {
//		if r.Err != nil { log.Println(r.Err); continue }
//		fmt.Printf("%d drifted resources at %s\n", len(r.Changes), r.At)
//	}
package watch
