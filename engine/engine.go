package engine

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/paulhalleux/workflow-engine-go/engine/internal"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/grpcapi"
	"github.com/paulhalleux/workflow-engine-go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type WorkflowEngineConfig = internal.WorkflowEngineConfig

type Engine struct {
	Config  *internal.WorkflowEngineConfig
	Context context.Context

	agentRegistry *internal.AgentRegistry
}

func NewEngine(
	config *internal.WorkflowEngineConfig,
) *Engine {
	ctx := context.Background()

	agentRegistry := internal.NewAgentRegistry()

	return &Engine{
		Config:  config,
		Context: ctx,

		agentRegistry: agentRegistry,
	}
}

func (e *Engine) Start() {
	log.Printf("[Engine] Starting workflow engine...")

	go e.startGrpcServer()

	<-e.Context.Done()
	log.Printf("[Engine] Shutting down workflow engine...")
}

func (e *Engine) startGrpcServer() {
	addr := joinHostPort(e.Config.GrpcAddress, e.Config.GrpcPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// Register gRPC services
	proto.RegisterEngineServiceServer(grpcServer, grpcapi.NewEngineServiceServer(e.agentRegistry))
	proto.RegisterTaskServiceServer(grpcServer, grpcapi.NewTaskServiceServer())

	// Start the gRPC server
	log.Printf("[Engine] gRPC server running on %s", lis.Addr().String())
	if serveErr := grpcServer.Serve(lis); serveErr != nil {
		log.Fatalf("failed to serve gRPC server: %v", serveErr)
	}
}

func joinHostPort(host *string, port string) string {
	if host == nil || *host == "" {
		return fmt.Sprintf(":%s", port)
	}
	return fmt.Sprintf("%s:%s", *host, port)
}
