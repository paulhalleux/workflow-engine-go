package queue

type Queue[T interface{}] interface {
	Enqueue(job T) error
	Dequeue() <-chan T
}
