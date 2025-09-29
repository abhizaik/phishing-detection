package analyzer

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// runTasks runs all tasks in parallel, collects timings and non-fatal errors
func runTasks(ctx context.Context, in *Input, tasks []Task) (*Output, []error) {
	out := &Output{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for _, t := range tasks {
		t := t
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				start := time.Now()
				if err := t.Run(in, out); err != nil {
					mu.Lock()
					errs = append(errs, fmt.Errorf("%s: %w", t.Name(), err))
					mu.Unlock()
				}
				out.setTiming(t.Name(), time.Since(start))
			}
		}()
	}

	wg.Wait()
	return out, errs
}
