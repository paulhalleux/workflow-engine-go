package agent

import (
	"context"

	"github.com/google/uuid"
	"github.com/paulhalleux/workflow-engine-go/internal/proto"
)

type ServiceServer struct {
	proto.UnimplementedAgentServiceServer
	agent *Agent
}

func NewAgentServiceServer(agent *Agent) *ServiceServer {
	return &ServiceServer{
		agent: agent,
	}
}

func (s *ServiceServer) StartTask(_ context.Context, req *proto.StartTaskRequest) (*proto.StartTaskResponse, error) {
	task, exists := s.agent.GetTask(req.TaskId)
	if !exists {
		return &proto.StartTaskResponse{
			Success: false,
			Message: "task not registered on agent",
		}, nil
	}

	executionContext := TaskExecutionContext{
		TaskId:      req.TaskId,
		Task:        *task,
		ExecutionId: uuid.New(),
		Input:       req.Parameters,
	}

	err := s.agent.Queue.Enqueue(&executionContext)
	if err != nil {
		return &proto.StartTaskResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &proto.StartTaskResponse{
		Success:     true,
		ExecutionId: executionContext.ExecutionId.String(),
	}, nil
}

func (s *ServiceServer) StopTask(ctx context.Context, req *proto.TaskActionRequest) (*proto.TaskActionResponse, error) {
	return nil, nil
}

func (s *ServiceServer) PauseTask(ctx context.Context, req *proto.TaskActionRequest) (*proto.TaskActionResponse, error) {
	return nil, nil
}

func (s *ServiceServer) ResumeTask(ctx context.Context, req *proto.TaskActionRequest) (*proto.TaskActionResponse, error) {
	return nil, nil
}
