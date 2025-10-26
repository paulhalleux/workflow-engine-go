package agent

import "github.com/google/uuid"

type TaskExecutionContext struct {
	ExecutionId uuid.UUID
	Input       map[string]interface{}
	TaskId      string
	Task        Task
}

type Task interface {
	Execute(context TaskExecutionContext) (map[string]interface{}, error)
}
