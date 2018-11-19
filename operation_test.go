package stream_test

import (
	"context"
	"testing"

	stream "."
)

func Test_Take(t *testing.T) {
	input := make(chan interface{})
	go func() {
		defer close(input)

		count := 0
		for {
			input <- count
			count++
		}
	}()

	wantSize := 5
	outch := stream.Take(wantSize)(context.TODO(), input)

	results := make([]interface{}, 0)
	for v := range outch {
		results = append(results, v)
	}

	if len(results) != wantSize {
		t.Fatalf("result size: want %d, get %d", wantSize, len(results))
	}
	for i, v := range results {
		if i != v.(int) {
			t.Fatalf("result value: want %q, get %q", i, v)
		}
	}
}

func Test_Filter(t *testing.T) {
	allsize := 10
	input := make(chan interface{})
	go func() {
		defer close(input)

		for i := 0; i < allsize; i++ {
			input <- i
		}
	}()

	outch := stream.Filter(func(v interface{}) bool {
		return v.(int)%2 == 0
	})(context.TODO(), input)

	results := make([]interface{}, 0)
	for v := range outch {
		results = append(results, v)
	}

	wantSize := allsize / 2
	if len(results) != wantSize {
		t.Fatalf("result size: want %d, get %d", wantSize, len(results))
	}
	for i, v := range results {
		want := i * 2
		if want != v.(int) {
			t.Fatalf("result value: want %q, get %q", want, v)
		}
	}
}

func Test_Map(t *testing.T) {
	allsize := 10
	input := make(chan interface{})
	go func() {
		defer close(input)

		for i := 0; i < allsize; i++ {
			input <- i
		}
	}()

	outch := stream.Map(func(v interface{}) interface{} {
		return v.(int) / 2
	})(context.TODO(), input)

	results := make([]interface{}, 0)
	for v := range outch {
		results = append(results, v)
	}

	wantSize := allsize
	if len(results) != wantSize {
		t.Fatalf("result size: want %d, get %d", wantSize, len(results))
	}
	for i, v := range results {
		want := i / 2
		if want != v.(int) {
			t.Fatalf("result value: want %q, get %q", want, v)
		}
	}
}
