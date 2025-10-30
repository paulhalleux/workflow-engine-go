package agent

import (
	"context"
	"log"
	"net"

	"github.com/paulhalleux/workflow-engine-go/agent/internal"
	"github.com/paulhalleux/workflow-engine-go/agent/internal/proto"
	"github.com/swaggest/jsonschema-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type WorkflowAgentConfig = internal.WorkflowAgentConfig
type TaskDefinition = internal.TaskDefinition

type WorkflowAgent struct {
	Config  *internal.WorkflowAgentConfig
	Context context.Context

	taskDefinitionRegistry *internal.TaskDefinitionRegistry
}

func NewWorkflowAgent(config *internal.WorkflowAgentConfig) *WorkflowAgent {
	ctx := context.Background()
	return &WorkflowAgent{
		Context:                ctx,
		Config:                 config,
		taskDefinitionRegistry: internal.NewTaskDefinitionRegistry(),
	}
}

func (a *WorkflowAgent) Start() {
	log.Printf("[Agent: %s] Starting workflow agent...", a.Config.Name)

	go func() {
		addr := joinHostPort(a.Config.Address, a.Config.Port)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		reflection.Register(grpcServer)

		// Register gRPC services
		proto.RegisterAgentServiceServer(grpcServer, internal.NewAgentServiceServer(
			a.Config,
			a.taskDefinitionRegistry,
		))

		// Start the gRPC server
		log.Printf("[Agent: %s] gRPC server running on %s", a.Config.Name, lis.Addr().String())
		if serveErr := grpcServer.Serve(lis); serveErr != nil {
			log.Fatalf("failed to serve gRPC server: %v", serveErr)
		}
	}()

	<-a.Context.Done()
}

func (a *WorkflowAgent) RegisterTaskDefinition(def TaskDefinition) {
	a.taskDefinitionRegistry.Register(def)
}

func ReflectJsonSchema(v interface{}) *jsonschema.Schema {
	reflector := jsonschema.Reflector{}
	r, err := reflector.Reflect(v)
	if err != nil {
		return &jsonschema.Schema{}
	}
	return &r
}

func joinHostPort(address *string, port string) string {
	if address != nil {
		return net.JoinHostPort(*address, port)
	}
	return net.JoinHostPort("", port)
}
