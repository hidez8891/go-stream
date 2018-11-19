package stream_test

import (
	"context"
	"testing"

	stream "."
)

func Test_InfinityFn(t *testing.T) {
	ctx, done := context.WithCancel(context.TODO())

	count := -1
	f := func() interface{} {
		count++
		return count
	}
	gench := stream.InfinityFn(f)(ctx, nil)

	results := make([]int, 0)
	stopSize := 5
	for v := range gench {
		v := v.(int)
		if v == stopSize {
			done()
		}
		results = append(results, v)
	}

	// len(results) >= stopSize
	// this case will not be tested

	for i, v := range results {
		if i != v {
			t.Fatalf("result value: want %q, get %q", i, v)
		}
	}
}

func Test_RepeatFn(t *testing.T) {
	wantSize := 5

	count := -1
	f := func() interface{} {
		count++
		return count
	}
	gench := stream.RepeatFn(wantSize, f)(context.TODO(), nil)

	results := make([]interface{}, 0)
	for v := range gench {
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
