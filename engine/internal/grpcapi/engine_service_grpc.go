package grpcapi

import (
	"context"
	"fmt"

	"github.com/paulhalleux/workflow-engine-go/engine/internal"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/service"
	"github.com/paulhalleux/workflow-engine-go/proto"
)

type EngineServiceServer struct {
	proto.UnimplementedEngineServiceServer

	agentRegistry   *internal.AgentRegistry
	workflowService *service.WorkflowExecutionService
}

func NewEngineServiceServer(
	agentRegistry *internal.AgentRegistry,
	workflowService *service.WorkflowExecutionService,
) *EngineServiceServer {
	return &EngineServiceServer{
		agentRegistry:   agentRegistry,
		workflowService: workflowService,
	}
}

func (e *EngineServiceServer) RegisterAgent(_ context.Context, req *proto.RegisterAgentRequest) (*proto.RegisterAgentResponse, error) {
	err := e.agentRegistry.RegisterAgent(
		req.Name,
		internal.RegisteredAgent{
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

func (e *EngineServiceServer) Ping(_ context.Context, req *proto.EnginePingRequest) (*proto.EnginePingResponse, error) {
	_, know := e.agentRegistry.GetAgent(req.Name)
	return &proto.EnginePingResponse{
		KnowAgent: know,
	}, nil
}

func (e *EngineServiceServer) StartWorkflow(_ context.Context, req *proto.StartWorkflowRequest) (*proto.StartWorkflowResponse, error) {
	workflowID, err := e.workflowService.StartWorkflow(
		req.WorkflowDefinitionId,
		req.InputParameters.AsMap(),
	)

	if err != nil {
		message := fmt.Sprintf("failed to start workflow: %v", err)
		return &proto.StartWorkflowResponse{
			Success: false,
			Message: &message,
		}, nil
	}

	return &proto.StartWorkflowResponse{
		WorkflowInstanceId: workflowID,
		Success:            true,
	}, nil
}
