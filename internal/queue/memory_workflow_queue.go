package queue

type memoryWorkflowQueue struct {
	ch chan WorkflowJob
}

func NewMemoryWorkflowQueue(buffer int) WorkflowQueue {
	return &memoryWorkflowQueue{ch: make(chan WorkflowJob, buffer)}
}

func (q *memoryWorkflowQueue) Enqueue(job WorkflowJob) error {
	q.ch <- job
	return nil
}

func (q *memoryWorkflowQueue) Dequeue() <-chan WorkflowJob {
	return q.ch
}
