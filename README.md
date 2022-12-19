# pipeline is a golang package for pipeline building based on golang channels and generics.

# Description
The goal is pipeline building using handlers.

# Example

```golang
package main

import (
	"context"
	"fmt"

	"github.com/nevisdale/pipeline"
)

func main() {
	numbers := make(chan int)
	go func() {
		for i := -100; i <= 100; i++ {
			numbers <- i
		}
		close(numbers)
	}()

	var p pipeline.Pipeline[int]
	p = p.Steps(
		pipeline.NewHandler(positive),
		pipeline.NewHandler(square),
		pipeline.NewHandler(lessOrEqual100),
	)

	for v := range p.Run(numbers) {
		fmt.Printf("v: %v\n", v)
	}
}

func positive(ctx context.Context, x int) (int, bool) {
	if x < 0 {
		return 0, false
	}
	return x, true
}

func square(ctx context.Context, x int) (int, bool) {
	return x * x, true
}

func lessOrEqual100(ctx context.Context, x int) (int, bool) {
	if x > 100 {
		return 0, false
	}
	return x, true
}
```

# License
[MIT](LICENSE)
