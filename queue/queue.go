// Package queue is a super simple FIFO concurrent safe queue implementation in Go.
package queue

import (
	"errors"
)

// Queue is a channel based concurrent safe queue.
type Queue[T any] struct {
	queue chan T
}

// New returns an empty concurrent safe queue. Panics if size is 0 or less.
func New[T any](size int) *Queue[T] {
	if size < 1 {
		panic(errors.New("queue size must be greater than 0"))
	}
	return &Queue[T]{
		queue: make(chan T, size),
	}
}

// Enqueue puts the given value v at the tail of the queue. If the queue is full, the operation is a no-op.
func (q *Queue[T]) Enqueue(v T) {
	select {
	case q.queue <- v:
	default:
	}
}

// Dequeue removes and returns the value at the head of the queue.
// It returns zero value of T if the queue is empty.
func (q *Queue[T]) Dequeue() T {
	select {
	case v := <-q.queue:
		return v
	default:
		var t T
		return t
	}
}

// HasNext checks if there is a value in the queue. If there is, it returns true and the value can be accessed by Dequeue().
func (q *Queue[T]) HasNext() bool {
	return q.Len() > 0
}

// Len Returns the current length of queue.
func (q *Queue[T]) Len() int {
	return len(q.queue)
}

// Cap returns the capacity of the queue.
func (q *Queue[T]) Cap() int {
	return cap(q.queue)
}
