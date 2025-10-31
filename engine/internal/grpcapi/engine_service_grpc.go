package grpcapi

import (
	"context"

	"github.com/paulhalleux/workflow-engine-go/engine/internal"
	"github.com/paulhalleux/workflow-engine-go/proto"
)

type EngineServiceServer struct {
	proto.UnimplementedEngineServiceServer

	agentRegistry *internal.AgentRegistry
}

func NewEngineServiceServer(
	agentRegistry *internal.AgentRegistry,
) *EngineServiceServer {
	return &EngineServiceServer{
		agentRegistry: agentRegistry,
	}
}

func (e EngineServiceServer) RegisterAgent(_ context.Context, req *proto.RegisterAgentRequest) (*proto.RegisterAgentResponse, error) {
	err := e.agentRegistry.RegisterAgent(
		req.Name,
		internal.RegisteredAgent{
			Name:           req.Name,
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
