package pipeline

import (
	"context"
)

type Pipeline[T any] struct {
	handlers []Handler[T]
}

func (p *Pipeline[T]) Use(handlers ...Handler[T]) {
	p.handlers = append(p.handlers, handlers...)
}

// Run creates a pipeline channel.
func (p Pipeline[T]) Run(in <-chan T) <-chan T {
	return p.RunContext(context.Background(), in)
}

// RunContext creates a pipeline channel with context
func (p Pipeline[T]) RunContext(ctx context.Context, in <-chan T) <-chan T {
	in = contextHanlder(ctx, in)
	for _, handler := range p.handlers {
		in = handler(ctx, in)
	}
	return in
}
