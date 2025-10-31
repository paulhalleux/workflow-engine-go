package internal

import (
	"context"

	"github.com/paulhalleux/workflow-engine-go/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AgentServiceServer struct {
	proto.UnimplementedAgentServiceServer

	config                 *WorkflowAgentConfig
	taskDefinitionRegistry *TaskDefinitionRegistry
	taskExecutionService   *TaskExecutionService
}

func NewAgentServiceServer(
	config *WorkflowAgentConfig,
	taskDefinitionRegistry *TaskDefinitionRegistry,
	taskExecutionService *TaskExecutionService,
) *AgentServiceServer {
	return &AgentServiceServer{
		config:                 config,
		taskDefinitionRegistry: taskDefinitionRegistry,
		taskExecutionService:   taskExecutionService,
	}
}

func (a *AgentServiceServer) StartTask(_ context.Context, req *proto.StartTaskRequest) (*proto.TaskActionResponse, error) {
	id := a.taskExecutionService.ExecuteTask(req)
	return &proto.TaskActionResponse{
		TaskId:  id,
		Success: true,
		Message: nil,
	}, nil
}

func (a *AgentServiceServer) GetTaskStatus(context.Context, *proto.TaskActionRequest) (*proto.GetTaskStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTaskStatus not implemented")
}

func (a *AgentServiceServer) StopTask(context.Context, *proto.TaskActionRequest) (*proto.TaskActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopTask not implemented")
}

func (a *AgentServiceServer) PauseTask(context.Context, *proto.TaskActionRequest) (*proto.TaskActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PauseTask not implemented")
}

func (a *AgentServiceServer) ResumeTask(context.Context, *proto.TaskActionRequest) (*proto.TaskActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResumeTask not implemented")
}

func (a *AgentServiceServer) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
