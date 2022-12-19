package pipeline

import (
	"context"
	"sync"
)

// Join joins several channels to a one.
func Join[T any](ins ...<-chan T) chan T {
	r := make(chan T)

	var wg sync.WaitGroup
	for _, in := range ins {
		wg.Add(1)

		go func(in <-chan T) {
			defer wg.Done()

			for x := range in {
				r <- x
			}
		}(in)
	}
	go func() {
		wg.Wait()
		close(r)
	}()

	return r
}

func contextHanlder[T any](ctx context.Context, in <-chan T) <-chan T {
	r := make(chan T)

	go func() {
		for x := range in {
			if !sendToChan(ctx, r, x) {
				break
			}
		}
		close(r)
	}()

	return r
}

func sendToChan[T any](ctx context.Context, ch chan<- T, value T) bool {
	select {
	case <-ctx.Done():
		return false
	case ch <- value:
		return true
	}
}
