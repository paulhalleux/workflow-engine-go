package internal

import (
	"context"

	"github.com/google/uuid"
	"github.com/paulhalleux/workflow-engine-go/proto"
	"google.golang.org/grpc"
)

type TaskExecutionService struct {
	taskExecutor     *TaskExecutor
	engineConnection *grpc.ClientConn
}

func NewTaskExecutionService(
	taskExecutor *TaskExecutor,
	engineConnection *grpc.ClientConn,
) *TaskExecutionService {
	return &TaskExecutionService{
		taskExecutor:     taskExecutor,
		engineConnection: engineConnection,
	}
}

func (tes *TaskExecutionService) ExecuteTask(req *proto.StartTaskRequest) string {
	id := uuid.New()

	tes.taskExecutor.EnqueueTask(&TaskExecution{
		TaskID:      id.String(),
		TaskDefName: req.TaskName,
		Input:       req.InputParameters.AsMap(),
	})

	client := proto.NewTaskServiceClient(tes.engineConnection)
	_, _ = client.NotifyTaskStatus(
		context.Background(),
		&proto.NotifyTaskStatusRequest{
			TaskId: id.String(),
			Status: proto.TaskStatus_PENDING,
		},
	)

	return id.String()
}
