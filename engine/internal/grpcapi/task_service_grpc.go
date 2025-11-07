package grpcapi

import (
	"context"
	"log"

	"github.com/paulhalleux/workflow-engine-go/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TaskServiceServer struct {
	proto.UnimplementedTaskServiceServer

	agentTaskChan map[string]chan *proto.NotifyTaskStatusRequest
}

func NewTaskServiceServer(
	agentTaskChan map[string]chan *proto.NotifyTaskStatusRequest,
) *TaskServiceServer {
	return &TaskServiceServer{
		agentTaskChan: agentTaskChan,
	}
}

func (s *TaskServiceServer) NotifyTaskStatus(_ context.Context, req *proto.NotifyTaskStatusRequest) (*emptypb.Empty, error) {
	log.Printf("NotifyTaskStatus called with TaskID: %s, Status: %v", req.TaskId, req.Status)
	if ch, ok := s.agentTaskChan[req.TaskId]; ok {
		ch <- req
	} else {
		log.Printf("No channel found for TaskID: %s", req.TaskId)
	}
	return &emptypb.Empty{}, nil
}

func (s *TaskServiceServer) NotifyTaskProgress(_ context.Context, req *proto.NotifyTaskProgressRequest) (*emptypb.Empty, error) {
	log.Printf("NotifyTaskProgress called with TaskID: %s, Progress: %f", req.TaskId, req.Progress)
	return nil, nil
}
