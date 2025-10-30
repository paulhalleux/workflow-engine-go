package internal

import (
	"context"

	"github.com/paulhalleux/workflow-engine-go/agent/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AgentServiceServer struct {
	proto.UnimplementedAgentServiceServer

	config                 *WorkflowAgentConfig
	taskDefinitionRegistry *TaskDefinitionRegistry
}

func NewAgentServiceServer(
	config *WorkflowAgentConfig,
	taskDefinitionRegistry *TaskDefinitionRegistry,
) *AgentServiceServer {
	return &AgentServiceServer{
		config:                 config,
		taskDefinitionRegistry: taskDefinitionRegistry,
	}
}

func (a *AgentServiceServer) StartTask(context.Context, *proto.StartTaskRequest) (*proto.TaskActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartTask not implemented")
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

func (a *AgentServiceServer) Ping(context.Context, *emptypb.Empty) (*proto.PingAgentResponse, error) {
	return &proto.PingAgentResponse{
		Name:           a.config.Name,
		Version:        a.config.Version,
		SupportedTasks: a.taskDefinitionRegistry.ToProto(),
	}, nil
}
