package queue

import "github.com/google/uuid"

type WorkflowJob struct {
	InstanceId uuid.UUID
}

type WorkflowQueue interface {
	Enqueue(job WorkflowJob) error
	Dequeue() <-chan WorkflowJob
}
