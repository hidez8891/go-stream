package stream

import (
	"context"
)

func InfinityFn(f func() interface{}) StreamFunc {
	return func(ctx context.Context, _ <-chan interface{}) <-chan interface{} {
		outch := make(chan interface{})

		go func() {
			defer close(outch)

			for {
				select {
				case <-ctx.Done():
					return

				case outch <- f():
					// continue
				}
			}
		}()

		return outch
	}
}

func RepeatFn(n int, f func() interface{}) StreamFunc {
	return func(ctx context.Context, _ <-chan interface{}) <-chan interface{} {
		outch := make(chan interface{})

		go func() {
			defer close(outch)

			for i := 0; i < n; i++ {
				select {
				case <-ctx.Done():
					return

				case outch <- f():
					// continue
				}
			}
		}()

		return outch
	}
}
