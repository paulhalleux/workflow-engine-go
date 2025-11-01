package agent

import (
	"context"
	"log"
	"net"

	"github.com/paulhalleux/workflow-engine-go/agent/internal"
	"github.com/paulhalleux/workflow-engine-go/agent/internal/connector"
	"github.com/paulhalleux/workflow-engine-go/agent/internal/ticker"
	"github.com/paulhalleux/workflow-engine-go/proto"
	"github.com/swaggest/jsonschema-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type WorkflowAgentConfig = internal.WorkflowAgentConfig
type TaskDefinition = internal.TaskDefinition
type TaskExecutionRequest = internal.TaskExecutionRequest
type TaskExecutionResult = internal.TaskExecutionResult

type WorkflowAgent struct {
	Config  *internal.WorkflowAgentConfig
	Context context.Context

	engineGrpcConnection *grpc.ClientConn

	taskExecutor           *internal.TaskExecutor
	taskExecutionService   *internal.TaskExecutionService
	taskDefinitionRegistry *internal.TaskDefinitionRegistry

	engineTicker    *ticker.EngineTicker
	engineConnector connector.EngineConnector
}

func NewWorkflowAgent(config *internal.WorkflowAgentConfig) *WorkflowAgent {
	ctx := context.Background()

	engineGrpcConnection, err := grpc.NewClient(config.EngineGrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to engine gRPC server: %v", err)
	}

	taskDefinitionRegistry := internal.NewTaskDefinitionRegistry()
	taskExecutor := internal.NewTaskExecutor(taskDefinitionRegistry, engineGrpcConnection)
	taskExecutionService := internal.NewTaskExecutionService(
		taskExecutor,
		engineGrpcConnection,
	)

	engineConnector, err := connector.NewEngineConnector(proto.AgentProtocol_GRPC, config.EngineGrpcUrl)
	if err != nil {
		log.Fatalf("Failed to create engine connector: %v", err)
	}

	engineTicker := ticker.NewEngineTicker(
		10,
		config,
		taskDefinitionRegistry,
		engineConnector,
	)

	return &WorkflowAgent{
		Context: ctx,
		Config:  config,

		engineGrpcConnection: engineGrpcConnection,

		taskExecutor:           taskExecutor,
		taskExecutionService:   taskExecutionService,
		taskDefinitionRegistry: taskDefinitionRegistry,

		engineTicker:    engineTicker,
		engineConnector: engineConnector,
	}
}

func (a *WorkflowAgent) Start() {
	log.Printf("[Agent: %s] Starting workflow agent...", a.Config.Name)
	defer func(engineGrpcConnection *grpc.ClientConn) {
		a.engineTicker.Stop()
		err := engineGrpcConnection.Close()
		if err != nil {
			log.Printf("[Agent: %s] Error closing engine gRPC connection: %v", a.Config.Name, err)
		}
	}(a.engineGrpcConnection)

	go a.engineTicker.Start()
	go a.taskExecutor.Start(a.Context)
	go a.startGrpcServer()

	_, err := a.engineConnector.RegisterAgent(a.Config, a.taskDefinitionRegistry)
	if err != nil {
		log.Printf("[Agent: %s] Failed to register agent with engine: %v, will retry every 10 seconds.", a.Config.Name, err.Error())
	} else {
		log.Printf("[Agent: %s] Registered with engine at %s", a.Config.Name, a.Config.EngineGrpcUrl)
	}

	<-a.Context.Done()
}

func (a *WorkflowAgent) startGrpcServer() {
	addr := joinHostPort(a.Config.GrpcAddress, a.Config.GrpcPort)
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
		a.taskExecutionService,
	))

	// Start the gRPC server
	log.Printf("[Agent: %s] gRPC server running on %s", a.Config.Name, lis.Addr().String())
	if serveErr := grpcServer.Serve(lis); serveErr != nil {
		log.Fatalf("failed to serve gRPC server: %v", serveErr)
	}
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
