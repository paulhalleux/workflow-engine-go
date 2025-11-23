package grpcserver

import (
	"context"

	"github.com/paulhalleux/workflow-engine-go/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AgentService struct {
	proto.UnimplementedAgentServiceServer
}

func NewAgentService() *AgentService {
	return &AgentService{}
}

func (s *AgentService) StartTask(context.Context, *proto.StartTaskRequest) (*proto.TaskActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartTask not implemented")
}
func (s *AgentService) GetTaskStatus(context.Context, *proto.TaskActionRequest) (*proto.GetTaskStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTaskStatus not implemented")
}
func (s *AgentService) StopTask(context.Context, *proto.TaskActionRequest) (*proto.TaskActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopTask not implemented")
}
func (s *AgentService) PauseTask(context.Context, *proto.TaskActionRequest) (*proto.TaskActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PauseTask not implemented")
}
func (s *AgentService) ResumeTask(context.Context, *proto.TaskActionRequest) (*proto.TaskActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResumeTask not implemented")
}
func (s *AgentService) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
