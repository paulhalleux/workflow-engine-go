package queue

type memoryQueue struct {
	ch chan WorkflowJob
}

func NewMemoryQueue(buffer int) WorkflowQueue {
	return &memoryQueue{ch: make(chan WorkflowJob, buffer)}
}

func (q *memoryQueue) Enqueue(job WorkflowJob) error {
	q.ch <- job
	return nil
}

func (q *memoryQueue) Dequeue() <-chan WorkflowJob {
	return q.ch
}
