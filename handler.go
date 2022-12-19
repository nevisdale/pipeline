package pipeline

import (
	"context"
	"sync"
)

type Handler[T any] func(ctx context.Context, in <-chan T) <-chan T

// HandlerFunc is a function for handling a value. You can return modified value.
// If a handler returns value with false, values will not send to a next channel.
type HandlerFunc[T any] func(ctx context.Context, x T) (_ T, valid bool)

// NewHandler creates a new handler for a Pipeline.
func NewHandler[T any](h HandlerFunc[T], opts ...HandlerOpt) Handler[T] {
	ho := hopt{
		concurrent: 1,
		cap:        0,
	}
	for _, opt := range opts {
		opt(&ho)
	}
	return func(ctx context.Context, in <-chan T) <-chan T {
		r := make(chan T, ho.cap)

		var wg sync.WaitGroup
		for i := 0; i < ho.concurrent; i++ {
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

type hopt struct {
	concurrent int
	cap        int
}

type HandlerOpt func(*hopt)

// WithConcurrent sets goroutines count. If n is less than 1, do nothing.
func WithConcurrent(n int) HandlerOpt {
	return func(h *hopt) {
		if n > 0 {
			h.concurrent = n
		}
	}
}

// WithCapacity sets channel capacity. If n is less than 0, do nothing.
func WithCapacity(n int) HandlerOpt {
	return func(h *hopt) {
		if n >= 0 {
			h.cap = n
		}
	}
}
