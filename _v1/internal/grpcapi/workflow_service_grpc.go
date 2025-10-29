package grpcapi

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/paulhalleux/workflow-engine-go/internal/proto"
	"github.com/paulhalleux/workflow-engine-go/internal/services"
	"google.golang.org/protobuf/types/known/emptypb"
)

type WorkflowServiceServer struct {
	proto.UnimplementedWorkflowServiceServer
	svc           services.WorkflowService
	agentTaskChan chan<- TaskExecutionResult
}

func NewWorkflowServiceServer(svc services.WorkflowService, agentTaskChan chan<- TaskExecutionResult) *WorkflowServiceServer {
	return &WorkflowServiceServer{
		svc:           svc,
		agentTaskChan: agentTaskChan,
	}
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

func (s *WorkflowServiceServer) NotifyTaskCompletion(_ context.Context, notif *proto.TaskCompletionNotification) (*emptypb.Empty, error) {
	log.Printf("Received task completion notification: ExecutionId=%s", notif.ExecutionId)
	s.agentTaskChan <- TaskExecutionResult{
		Type:        TaskExecutionResultTypeCompletion,
		ExecutionId: uuid.MustParse(notif.ExecutionId),
		Output:      notif.Output.AsMap(),
	}
	return &emptypb.Empty{}, nil
}

func (s *WorkflowServiceServer) NotifyTaskFailure(_ context.Context, notif *proto.TaskFailureNotification) (*emptypb.Empty, error) {
	log.Printf("Received task failure notification")
	s.agentTaskChan <- TaskExecutionResult{
		Type:        TaskExecutionResultTypeFailure,
		ExecutionId: uuid.MustParse(notif.ExecutionId),
		Error:       errors.New(notif.ErrorMessage),
	}
	return &emptypb.Empty{}, nil
}

func (s *WorkflowServiceServer) NotifyTaskProgress(_ context.Context, notif *proto.TaskProgressNotification) (*emptypb.Empty, error) {
	log.Printf("Received task progress notification: %d", notif.ProgressPercentage)
	s.agentTaskChan <- TaskExecutionResult{
		Type:        TaskExecutionResultTypeProgress,
		ExecutionId: uuid.MustParse(notif.ExecutionId),
		Progress:    float32(notif.ProgressPercentage) / 100.0,
	}
	return &emptypb.Empty{}, nil
}
