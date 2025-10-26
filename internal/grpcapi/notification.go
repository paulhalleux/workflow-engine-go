package grpcapi

import "github.com/google/uuid"

type TaskExecutionResultType string

const (
	TaskExecutionResultTypeCompletion TaskExecutionResultType = "completion"
	TaskExecutionResultTypeFailure    TaskExecutionResultType = "failure"
	TaskExecutionResultTypeProgress   TaskExecutionResultType = "progress"
)

type TaskExecutionResult struct {
	Type        TaskExecutionResultType
	ExecutionId uuid.UUID
	Error       error
	Output      map[string]interface{}
	Progress    float32
}
