package grpcserver

import (
	"context"

	"github.com/paulhalleux/workflow-engine-go/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TaskService struct {
	proto.UnimplementedTaskServiceServer
}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (s *TaskService) NotifyTaskStatus(context.Context, *proto.NotifyTaskStatusRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyTaskStatus not implemented")
}

func (s *TaskService) NotifyTaskProgress(context.Context, *proto.NotifyTaskProgressRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyTaskProgress not implemented")
}
