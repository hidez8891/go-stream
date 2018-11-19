package stream

import (
	"context"
	"sync"
)

type StreamFunc func(ctx context.Context, input <-chan interface{}) <-chan interface{}

func Connect(ctx context.Context, fs ...StreamFunc) <-chan interface{} {
	subctx, done := context.WithCancel(ctx)
	ch := Serial(fs...)(subctx, nil)

	outch := make(chan interface{})
	go func() {
		defer close(outch)
		defer done()

		for v := range orDone(ctx.Done(), ch) {
			outch <- v
		}
	}()

	return outch
}

func Serial(fs ...StreamFunc) StreamFunc {
	return func(ctx context.Context, input <-chan interface{}) <-chan interface{} {
		outch := make(chan interface{})

		var ch = fs[0](ctx, input)
		for _, f := range fs[1:] {
			ch = f(ctx, ch)
		}

		go func() {
			defer close(outch)
			for v := range orDone(ctx.Done(), ch) {
				outch <- v
			}
		}()

		return outch
	}
}

func Parallel(n int, fs ...StreamFunc) StreamFunc {
	return func(ctx context.Context, input <-chan interface{}) <-chan interface{} {
		outch := make(chan interface{})
		wg := &sync.WaitGroup{}

		wg.Add(n)
		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()

				ch := Serial(fs...)(ctx, input)
				for v := range orDone(ctx.Done(), ch) {
					outch <- v
				}
			}()
		}

		go func() {
			defer close(outch)
			wg.Wait()
		}()

		return outch
	}
}

func If(f func(interface{}) bool, st1, st2 StreamFunc) StreamFunc {
	return func(ctx context.Context, input <-chan interface{}) <-chan interface{} {
		outch := make(chan interface{})
		input1 := make(chan interface{})
		input2 := make(chan interface{})

		heartbeat := make(chan struct{}, 1)
		heartbeat <- struct{}{}

		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()

			ch := st1(ctx, input1)
			for v := range orDone(ctx.Done(), ch) {
				outch <- v
				heartbeat <- struct{}{}
			}
		}()

		go func() {
			defer wg.Done()

			ch := st2(ctx, input2)
			for v := range orDone(ctx.Done(), ch) {
				outch <- v
				heartbeat <- struct{}{}
			}
		}()

		go func() {
			defer close(heartbeat)
			defer close(outch)
			wg.Wait()
		}()

		go func() {
			defer close(input1)
			defer close(input2)

			for v := range orDone(ctx.Done(), input) {
				if _, ok := <-heartbeat; !ok {
					return
				}
				if f(v) {
					select {
					case <-ctx.Done():
					case input1 <- v:
					}
				} else {
					select {
					case <-ctx.Done():
					case input2 <- v:
					}
				}
			}
		}()

		return outch
	}
}

func StopIf(f func(interface{}) bool) StreamFunc {
	return func(ctx context.Context, input <-chan interface{}) <-chan interface{} {
		outch := make(chan interface{})

		go func() {
			defer close(outch)

			for v := range orDone(ctx.Done(), input) {
				if f(v) {
					return
				}
				outch <- v
			}
		}()

		return outch
	}
}
