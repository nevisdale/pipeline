package pipeline

import (
	"context"
	"sync"
)

type Handler[T any] func(ctx context.Context, in <-chan T) <-chan T

// HandlerFunc is a function for handling a value. You can return modified value.
// If a handler returns value with false, values will not send to a next channel.
type HandlerFunc[T any] func(ctx context.Context, x T) (_ T, passed bool)

// NewHandler creates a new handler for a Pipeline.
func NewHandler[T any](h HandlerFunc[T]) Handler[T] {
	return NewHandlerConcurrent(h, 1)
}

// NewHandler creates a new concurrent handler for a Pipeline.
func NewHandlerConcurrent[T any](h HandlerFunc[T], concurrent int) Handler[T] {
	if concurrent < 1 {
		concurrent = 1
	}

	return func(ctx context.Context, in <-chan T) <-chan T {
		r := make(chan T)

		var wg sync.WaitGroup
		for i := 0; i < concurrent; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for x := range in {
					var passed bool
					if x, passed = h(ctx, x); passed {
						r <- x
					}
				}
			}()
		}
		go func() {
			wg.Wait()
			close(r)
		}()
		return r
	}
}
