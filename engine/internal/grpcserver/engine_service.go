package grpcserver

import (
	"context"

	"github.com/paulhalleux/workflow-engine-go/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EngineService struct {
	proto.UnimplementedEngineServiceServer
}

func NewEngineService() *EngineService {
	return &EngineService{}
}

func (s *EngineService) RegisterAgent(context.Context, *proto.RegisterAgentRequest) (*proto.RegisterAgentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterAgent not implemented")
}

func (s *EngineService) Ping(context.Context, *proto.EnginePingRequest) (*proto.EnginePingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}

func (s *EngineService) StartWorkflow(context.Context, *proto.StartWorkflowRequest) (*proto.StartWorkflowResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartWorkflow not implemented")
}
