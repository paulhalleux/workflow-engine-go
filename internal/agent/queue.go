package agent

import "errors"

type TaskQueue struct {
	tasksChannel chan *TaskExecutionContext
}

func NewTaskQueue(buffer int) *TaskQueue {
	return &TaskQueue{
		tasksChannel: make(chan *TaskExecutionContext, buffer),
	}
}

func (q *TaskQueue) Enqueue(taskContext *TaskExecutionContext) error {
	select {
	case q.tasksChannel <- taskContext:
		return nil
	default:
		return errors.New("task queue is full")
	}
}

func (q *TaskQueue) Dequeue() <-chan *TaskExecutionContext {
	return q.tasksChannel
}
