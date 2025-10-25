package agent

import (
	"context"

	"github.com/paulhalleux/workflow-engine-go/internal/proto"
)

type ServiceServer struct {
	proto.UnimplementedAgentServiceServer
}

func NewAgentServiceServer() *ServiceServer {
	return &ServiceServer{}
}

func (s *ServiceServer) StartTask(ctx context.Context, req *proto.StartTaskRequest) (*proto.StartTaskResponse, error) {
	return nil, nil
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
