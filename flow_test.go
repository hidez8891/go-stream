package stream_test

import (
	"context"
	"sort"
	"testing"
	"time"

	stream "."
)

func Test_Connect(t *testing.T) {
	wantSize := 5

	count := -1
	gen := func() interface{} {
		count++
		return count
	}

	isEven := func(v interface{}) bool {
		return v.(int)%2 == 0
	}

	outch := stream.Connect(context.TODO(),
		stream.InfinityFn(gen),
		stream.Filter(isEven),
		stream.Take(wantSize),
	)

	results := make([]interface{}, 0)
	for v := range outch {
		results = append(results, v)
	}

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

func Test_Connect_CloseStream(t *testing.T) {
	var gench <-chan interface{}
	genStream := func(ctx context.Context, input <-chan interface{}) <-chan interface{} {
		gench = stream.InfinityFn(func() interface{} { return 0 })(ctx, input)
		return gench
	}

	outch := stream.Connect(context.TODO(),
		genStream,
		stream.Take(5),
	)

	for _ = range outch {
		// discard
	}
	select {
	case _, ok := <-gench:
		if ok {
			t.Fatalf("test fail: stream is not closed")
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("test timeout: stream is not closed")
	}
}

func Test_Serial(t *testing.T) {
	wantSize := 5

	count := -1
	gen := func() interface{} {
		count++
		return count
	}

	isEven := func(v interface{}) bool {
		return v.(int)%2 == 0
	}

	outch := stream.Serial(
		stream.InfinityFn(gen),
		stream.Filter(isEven),
		stream.Take(wantSize),
	)(context.TODO(), nil)

	results := make([]interface{}, 0)
	for v := range outch {
		results = append(results, v)
	}

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

func Test_Parallel(t *testing.T) {
	parallelSize := 3
	wantSize := 100

	count := -1
	gen := func() interface{} {
		count++
		return count
	}

	isEven := func(v interface{}) bool {
		return v.(int)%2 == 0
	}

	outch := stream.Serial(
		stream.RepeatFn(wantSize*2, gen),
		stream.Parallel(
			parallelSize,
			stream.Filter(isEven),
		),
	)(context.TODO(), nil)

	results := make([]interface{}, 0)
	for v := range outch {
		results = append(results, v)
	}

	if len(results) != wantSize {
		t.Fatalf("result size: want %d, get %d", wantSize, len(results))
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].(int) < results[j].(int)
	})
	for i, v := range results {
		want := i * 2
		if want != v.(int) {
			t.Fatalf("result value: want %d, get %d", want, v)
		}
	}
}

func Test_If(t *testing.T) {
	count := -1
	gen := func() interface{} {
		count++
		return count
	}

	isEven := func(v interface{}) bool {
		return v.(int)%2 == 0
	}

	neg := func(v interface{}) interface{} {
		return -v.(int)
	}

	id := func(v interface{}) interface{} {
		return v.(int)
	}

	wantSize := 100
	outch := stream.Connect(context.TODO(),
		stream.InfinityFn(gen),
		stream.If(
			isEven,
			stream.Serial(
				stream.Map(neg),
			),
			stream.Serial(
				stream.Map(id),
			),
		),
		stream.Take(wantSize),
	)

	results := make([]interface{}, 0)
	for v := range outch {
		results = append(results, v)
	}

	if len(results) != wantSize {
		t.Fatalf("result size: want %d, get %d", wantSize, len(results))
	}
	for i, v := range results {
		want := i
		if want%2 == 0 {
			want = -want
		}
		if want != v.(int) {
			t.Fatalf("result value: want %d, get %d", want, v)
		}
	}
}

func Test_StopIf(t *testing.T) {
	gen := func() interface{} {
		return 0
	}

	isCalled := false
	testFunc := func(v interface{}) bool {
		isCalled = true
		return true
	}

	outch := stream.Connect(context.TODO(),
		stream.InfinityFn(gen),
		stream.StopIf(func(_ interface{}) bool { return true }),
		stream.Filter(testFunc),
		stream.Take(1),
	)

	for _ = range outch {
		t.Fatalf("stream is readable: stream don't stoped")
	}
	if isCalled {
		t.Fatalf("function is called: stream don't stoped")
	}
}
