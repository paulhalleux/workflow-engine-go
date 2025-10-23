package grpcapi

import (
	"context"

	"github.com/google/uuid"
	"github.com/paulhalleux/workflow-engine-go/internal/proto"
	"github.com/paulhalleux/workflow-engine-go/internal/services"
)

type WorkflowEngineServer struct {
	proto.UnimplementedWorkflowEngineServer
	svc services.WorkflowService
}

func NewWorkflowEngineServer(svc services.WorkflowService) *WorkflowEngineServer {
	return &WorkflowEngineServer{svc: svc}
}

func (s *WorkflowEngineServer) StartWorkflow(ctx context.Context, req *proto.StartWorkflowRequest) (*proto.StartWorkflowResponse, error) {
	defID, err := uuid.Parse(req.WorkflowDefinitionId)
	if err != nil {
		return nil, err
	}

	instanceID, err := s.svc.StartWorkflow(ctx, defID, req.Input, req.Metadata)
	if err != nil {
		return nil, err
	}

	return &proto.StartWorkflowResponse{WorkflowInstanceId: instanceID.String()}, nil
}
