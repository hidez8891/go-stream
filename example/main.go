package main

import (
	"context"
	"fmt"
	"math/rand"

	stream ".."
)

func main() {
	exampleShowPrime()
	exampleShowFizzBuzz()
}

func exampleShowPrime() {
	const (
		jobs = 5
		size = 20
	)

	gen := func() interface{} {
		return rand.Int()
	}

	isPrime := func(v interface{}) bool {
		val := v.(int)
		if val == 2 || val == 3 {
			return true
		}
		if val < 2 || val%2 == 0 {
			return false
		}
		for p := 5; p < val/2; p = p + 2 {
			if val%p == 0 {
				return false
			}
		}
		return true
	}

	primes := stream.Connect(
		context.Background(),
		stream.InfinityFn(gen),
		stream.Parallel(
			jobs,
			stream.Filter(isPrime),
		),
		stream.Take(size),
	)

	for p := range primes {
		fmt.Println(p.(int))
	}
}

func exampleShowFizzBuzz() {
	const (
		size = 20
	)

	count := 0
	gen := func() interface{} {
		count++
		return count
	}

	primes := stream.Connect(
		context.Background(),
		stream.RepeatFn(size, gen),
		stream.If(
			func(v interface{}) bool { return v.(int)%15 == 0 },
			stream.Map(func(v interface{}) interface{} { return "FizzBuzz" }),
			stream.If(
				func(v interface{}) bool { return v.(int)%3 == 0 },
				stream.Map(func(v interface{}) interface{} { return "Fizz" }),
				stream.If(
					func(v interface{}) bool { return v.(int)%5 == 0 },
					stream.Map(func(v interface{}) interface{} { return "Buzz" }),
					stream.Map(func(v interface{}) interface{} { return v }),
				),
			),
		),
	)

	for p := range primes {
		fmt.Println(p)
	}
}
