package agent

import (
	"log"
	"net"
	"strconv"

	"github.com/paulhalleux/workflow-engine-go/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type GrpcInfo struct {
	Address *string
	Port    int
}

type Agent struct {
	Name           string
	Grpc           GrpcInfo
	WorkflowClient proto.WorkflowServiceClient
	Tasks          map[string]*Task
	Queue          *TaskQueue
}

func NewAgent(name string, grpcAddress *string, grpcPort int) *Agent {
	conn, err := grpc.NewClient("127.0.0.1", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to workflow server: %v", err)
	}

	workflowClient := proto.NewWorkflowServiceClient(conn)

	return &Agent{
		Name: name,
		Grpc: GrpcInfo{
			Address: grpcAddress,
			Port:    grpcPort,
		},
		Tasks:          make(map[string]*Task),
		Queue:          NewTaskQueue(100),
		WorkflowClient: workflowClient,
	}
}

func (a *Agent) Start() {
	taskExecutor := NewTaskExecutor(a.Queue, 10)
	taskExecutor.Start()

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
	var lis net.Listener
	if a.Grpc.Address != nil {
		lis, _ = net.Listen("tcp", *a.Grpc.Address+":"+strconv.Itoa(a.Grpc.Port))
	} else {
		lis, _ = net.Listen("tcp", ":"+strconv.Itoa(a.Grpc.Port))
	}

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
