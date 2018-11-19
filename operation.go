package stream

import (
	"context"
)

func Take(n int) StreamFunc {
	return func(ctx context.Context, input <-chan interface{}) <-chan interface{} {
		outch := make(chan interface{})

		go func() {
			defer close(outch)

			ch := orDone(ctx.Done(), input)
			for i := 0; i < n; i++ {
				v, ok := <-ch
				if !ok {
					break
				}
				outch <- v
			}
		}()

		return outch
	}
}

func Filter(f func(interface{}) bool) StreamFunc {
	return func(ctx context.Context, input <-chan interface{}) <-chan interface{} {
		outch := make(chan interface{})

		go func() {
			defer close(outch)

			for v := range orDone(ctx.Done(), input) {
				if f(v) {
					outch <- v
				}
			}
		}()

		return outch
	}
}

func Map(f func(interface{}) interface{}) StreamFunc {
	return func(ctx context.Context, input <-chan interface{}) <-chan interface{} {
		outch := make(chan interface{})

		go func() {
			defer close(outch)

			for v := range orDone(ctx.Done(), input) {
				outch <- f(v)
			}
		}()

		return outch
	}
}
