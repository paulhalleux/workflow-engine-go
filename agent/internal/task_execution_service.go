package internal

import (
	"github.com/google/uuid"
	"github.com/paulhalleux/workflow-engine-go/agent/internal/proto"
)

type TaskExecutionService struct {
	taskExecutor *TaskExecutor
}

func NewTaskExecutionService(
	taskExecutor *TaskExecutor,
) *TaskExecutionService {
	return &TaskExecutionService{
		taskExecutor: taskExecutor,
	}
}

func (tes *TaskExecutionService) ExecuteTask(req *proto.StartTaskRequest) string {
	id := uuid.New()

	tes.taskExecutor.EnqueueTask(&TaskExecution{
		TaskID:      id.String(),
		TaskDefName: req.TaskName,
		Input:       req.InputParameters.AsMap(),
	})

	return id.String()
}
