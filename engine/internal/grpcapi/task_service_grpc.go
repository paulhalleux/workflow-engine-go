package grpcapi

import (
	"context"
	"log"

	"github.com/paulhalleux/workflow-engine-go/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TaskServiceServer struct {
	proto.UnimplementedTaskServiceServer
}

func NewTaskServiceServer() *TaskServiceServer {
	return &TaskServiceServer{}
}

func (s *TaskServiceServer) NotifyTaskStatus(ctx context.Context, req *proto.NotifyTaskStatusRequest) (*emptypb.Empty, error) {
	log.Printf("NotifyTaskStatus called with TaskID: %s, Status: %v, OutputParameters: %v", req.TaskId, req.Status, req.OutputParameters)
	return nil, nil
}

func (s *TaskServiceServer) NotifyTaskProgress(ctx context.Context, req *proto.NotifyTaskProgressRequest) (*emptypb.Empty, error) {
	log.Printf("NotifyTaskProgress called with TaskID: %s, Progress: %f", req.TaskId, req.Progress)
	return nil, nil
}
