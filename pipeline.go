package pipeline

import (
	"context"
	"sync"
)

type Pipeline[T any] struct {
	handlers []Handler[T]
}

func (p *Pipeline[T]) Use(handlers ...Handler[T]) {
	p.handlers = append(p.handlers, handlers...)
}

func (p Pipeline[T]) Build(ctx context.Context, ins ...<-chan T) <-chan T {
	r := make(chan T)

	var wg sync.WaitGroup
	for _, in := range ins {
		wg.Add(1)

		go func(in <-chan T) {
			defer wg.Done()

			for v := range in {
				r <- v
			}
		}(in)
	}
	go func() {
		wg.Wait()
		close(r)
	}()

	var ro <-chan T = r
	for _, handler := range p.handlers {
		ro = handler(ctx, ro)
	}
	return ro
}
