package agent

import (
	"log"
	"net"
	"strconv"

	"github.com/paulhalleux/workflow-engine-go/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcInfo struct {
	Address string
	Port    int
}

type Agent struct {
	Name  string
	Grpc  GrpcInfo
	Tasks map[string]*Task
}

func (a *Agent) Start() {
	go startGrpcServer(a)
}

func (a *Agent) RegisterTask(name string, task Task) {
	if a.Tasks == nil {
		a.Tasks = make(map[string]*Task)
	}
	a.Tasks[name] = &task
}

func (a *Agent) GetTask(name string) (*Task, bool) {
	task, exists := a.Tasks[name]
	return task, exists
}

func startGrpcServer(a *Agent) {
	addr := net.JoinHostPort(a.Grpc.Address, strconv.Itoa(a.Grpc.Port))
	lis, _ := net.Listen("tcp", addr)
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// Register gRPC services
	proto.RegisterAgentServiceServer(grpcServer, NewAgentServiceServer())

	// Start the gRPC server
	log.Printf("[Agent: %s] gRPC server running on %s", a.Name, addr)
	if serveErr := grpcServer.Serve(lis); serveErr != nil {
		log.Fatalf("failed to serve gRPC server: %v", serveErr)
	}
}
