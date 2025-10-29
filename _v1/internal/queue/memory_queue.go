package queue

import (
	"errors"
	"fmt"
)

type MemoryQueue[T interface{}] struct {
	buffer int
	ch     chan T
}

func NewMemoryQueue[T interface{}](buffer int) MemoryQueue[T] {
	return MemoryQueue[T]{
		ch:     make(chan T, buffer),
		buffer: buffer,
	}
}

func (q MemoryQueue[T]) Enqueue(job T) error {
	select {
	case q.ch <- job:
		return nil
	default:
		return errors.New(fmt.Sprintf("Queue is full, exceeded buffer of %d", q.buffer))
	}
}

func (q MemoryQueue[T]) Dequeue() <-chan T {
	return q.ch
}
