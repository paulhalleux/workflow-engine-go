package agent

import "github.com/google/uuid"

type TaskExecutionContext struct {
	ExecutionId uuid.UUID
	Input       interface{}
}

type Task interface {
	Execute(context TaskExecutionContext) (interface{}, error)
}
