package agent

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type WorkflowAgentConfig struct {
	Name    string
	Address *string
	Port    string
}

type WorkflowAgent struct {
	Config  *WorkflowAgentConfig
	Context context.Context
}

func NewWorkflowAgent(config *WorkflowAgentConfig) *WorkflowAgent {
	ctx := context.Background()
	return &WorkflowAgent{
		Context: ctx,
		Config:  config,
	}
}

func (a *WorkflowAgent) Start() {
	fmt.Println("Workflow Agent started...")

	go func() {
		addr := joinHostPort(a.Config.Address, a.Config.Port)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		reflection.Register(grpcServer)

		// Register gRPC services

		// Start the gRPC server
		log.Printf("[Agent: %s] gRPC server running on %s", a.Config.Name, lis.Addr().String())
		if serveErr := grpcServer.Serve(lis); serveErr != nil {
			log.Fatalf("failed to serve gRPC server: %v", serveErr)
		}
	}()

	<-a.Context.Done()
}

func joinHostPort(address *string, port string) string {
	if address != nil {
		return net.JoinHostPort(*address, port)
	}
	return net.JoinHostPort("", port)
}
