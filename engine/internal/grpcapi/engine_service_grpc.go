package grpcapi

import (
	"context"
	"log"

	"github.com/paulhalleux/workflow-engine-go/proto"
)

type EngineServiceServer struct {
	proto.UnimplementedEngineServiceServer
}

func NewEngineServiceServer() *EngineServiceServer {
	return &EngineServiceServer{}
}

func (e EngineServiceServer) RegisterAgent(ctx context.Context, req *proto.RegisterAgentRequest) (*proto.RegisterAgentResponse, error) {
	log.Printf("Received RegisterAgent request, %v", req)
	return nil, nil
}
