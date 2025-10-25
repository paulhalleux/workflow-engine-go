package grpcapi

import (
	"context"

	"github.com/google/uuid"
	"github.com/paulhalleux/workflow-engine-go/internal/proto"
	"github.com/paulhalleux/workflow-engine-go/internal/services"
)

type WorkflowServiceServer struct {
	proto.UnimplementedWorkflowServiceServer
	svc services.WorkflowService
}

func NewWorkflowServiceServer(svc services.WorkflowService) *WorkflowServiceServer {
	return &WorkflowServiceServer{svc: svc}
}

func (s *WorkflowServiceServer) StartWorkflow(ctx context.Context, req *proto.StartWorkflowRequest) (*proto.StartWorkflowResponse, error) {
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
