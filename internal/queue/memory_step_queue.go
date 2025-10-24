package queue

type memoryStepQueue struct {
	ch chan StepJob
}

func NewMemoryStepQueue(buffer int) StepQueue {
	return &memoryStepQueue{ch: make(chan StepJob, buffer)}
}

func (q *memoryStepQueue) Enqueue(job StepJob) error {
	q.ch <- job
	return nil
}

func (q *memoryStepQueue) Dequeue() <-chan StepJob {
	return q.ch
}
