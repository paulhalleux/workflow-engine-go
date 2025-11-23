package grpcserver

import (
	"context"

	"github.com/paulhalleux/workflow-engine-go/engine-new/internal/registry"
	"github.com/paulhalleux/workflow-engine-go/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EngineService struct {
	proto.UnimplementedEngineServiceServer

	agentRegistry *registry.AgentRegistry
}

func NewEngineService(
	agentRegistry *registry.AgentRegistry,
) *EngineService {
	return &EngineService{
		agentRegistry: agentRegistry,
	}
}

func (s *EngineService) RegisterAgent(_ context.Context, req *proto.RegisterAgentRequest) (*proto.RegisterAgentResponse, error) {
	err := s.agentRegistry.RegisterAgent(
		req.Name,
		registry.RegisteredAgent{
			Name:           req.Name,
			Version:        req.Version,
			Address:        req.Address,
			Port:           req.Port,
			Protocol:       req.Protocol,
			SupportedTasks: req.SupportedTasks,
		},
	)

	if err != nil {
		message := err.Error()
		return &proto.RegisterAgentResponse{
			Success: false,
			Message: &message,
		}, nil
	}

	return &proto.RegisterAgentResponse{
		Success: true,
	}, nil
}

func (s *EngineService) Ping(_ context.Context, req *proto.EnginePingRequest) (*proto.EnginePingResponse, error) {
	_, know := s.agentRegistry.GetAgent(req.Name)
	return &proto.EnginePingResponse{
		KnowAgent: know,
	}, nil
}

func (s *EngineService) StartWorkflow(context.Context, *proto.StartWorkflowRequest) (*proto.StartWorkflowResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartWorkflow not implemented")
}
