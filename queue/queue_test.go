package queue_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/tigorlazuardi/maleo/queue"
)

func TestQueue(t *testing.T) {
	q := queue.New[int](5000)
	count := uint64(0)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(5000)
	for i := 1; i <= 5000; i++ {
		q.Enqueue(i)
		go func() {
			<-ctx.Done()
			j := q.Dequeue()
			if j == 0 {
				t.Error("unexpected 0 value from queue. there should be no 0 value.")
			}
			atomic.AddUint64(&count, 1)
			wg.Done()
		}()
	}
	if q.Len() != 5000 {
		t.Errorf("expected queue to have 5000 length, but got %d length", q.Len())
	}
	cancel()
	wg.Wait()
	if q.Len() != 0 {
		t.Errorf("expected queue to have 0 length, but got %d length", q.Len())
	}
	if count != 5000 {
		t.Errorf("expected count to be 5000, but got %d", count)
	}
}

func TestQueue2(t *testing.T) {
	const queueInsert = 5
	const wantCount = 2
	q := queue.New[int](wantCount)
	if q.Cap() != wantCount {
		t.Errorf("expected queue to have %d capacity, but got %d capacity", wantCount, q.Cap())
	}
	for i := 1; i <= queueInsert; i++ {
		q.Enqueue(i)
	}
	if q.Len() != wantCount {
		t.Errorf("expected queue to have %d length, but got %d length", wantCount, q.Len())
	}
	start := 1
	for q.HasNext() {
		got := q.Dequeue()
		if got != start {
			t.Errorf("expected %d, but got %d", start, got)
		}
		start++
	}
	shouldEmpty := q.Dequeue()
	if shouldEmpty != 0 {
		t.Errorf("expected 0, but got %d", shouldEmpty)
	}
}

func BenchmarkQueue(b *testing.B) {
	q := queue.New[int](b.N)
	wg := sync.WaitGroup{}
	wg.Add(b.N * 2)
	for i := 0; i < b.N; i++ {
		go func(i int) {
			q.Enqueue(i)
			wg.Done()
		}(i)
		go func() {
			q.Dequeue()
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestNewPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got nil")
		}
	}()
	queue.New[int](0)
}
