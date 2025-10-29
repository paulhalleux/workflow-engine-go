package agent

import (
	"context"
	"log"
	"net"

	"github.com/paulhalleux/workflow-engine-go/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type GrpcInfo struct {
	Address string
	Port    string
}

type Agent struct {
	Name           string
	Context        context.Context
	Grpc           GrpcInfo
	WorkflowClient proto.WorkflowServiceClient
	Tasks          map[string]*Task
	Queue          *TaskQueue
}

func NewAgent(name string, grpcAddress string, grpcPort string) *Agent {
	conn, err := grpc.NewClient("engine:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to workflow server: %v", err)
	}

	workflowClient := proto.NewWorkflowServiceClient(conn)
	ctx := context.Background()

	return &Agent{
		Name: name,
		Grpc: GrpcInfo{
			Address: grpcAddress,
			Port:    grpcPort,
		},
		Context:        ctx,
		Tasks:          make(map[string]*Task),
		Queue:          NewTaskQueue(100),
		WorkflowClient: workflowClient,
	}
}

func (a *Agent) Start() {
	taskExecutor := NewTaskExecutor(a, 10)
	taskExecutor.Start(a.Context)
	go startGrpcServer(a)
}

func (a *Agent) RegisterTask(name string, task Task) {
	a.Tasks[name] = &task
}

func (a *Agent) GetTask(name string) (*Task, bool) {
	task, exists := a.Tasks[name]
	return task, exists
}

func startGrpcServer(a *Agent) {
	addr := net.JoinHostPort(a.Grpc.Address, a.Grpc.Port)
	lis, _ := net.Listen("tcp", addr)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// Register gRPC services
	proto.RegisterAgentServiceServer(grpcServer, NewAgentServiceServer(a))

	// Start the gRPC server
	log.Printf("[Agent: %s] gRPC server running on %s", a.Name, lis.Addr().String())
	if serveErr := grpcServer.Serve(lis); serveErr != nil {
		log.Fatalf("failed to serve gRPC server: %v", serveErr)
	}
}
