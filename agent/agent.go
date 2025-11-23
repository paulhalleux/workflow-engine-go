package agent

import (
	"context"
	"log"

	"github.com/paulhalleux/workflow-engine-go/agent/internal/connector"
	"github.com/paulhalleux/workflow-engine-go/agent/internal/executor"
	"github.com/paulhalleux/workflow-engine-go/agent/internal/grpcserver"
	"github.com/paulhalleux/workflow-engine-go/agent/internal/models"
	"github.com/paulhalleux/workflow-engine-go/agent/internal/registry"
	"github.com/paulhalleux/workflow-engine-go/agent/internal/ticker"
	"github.com/paulhalleux/workflow-engine-go/proto"
	"github.com/swaggest/jsonschema-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TaskDefinition = models.TaskDefinition
type TaskExecutionRequest = models.TaskExecutionRequest
type TaskExecutionResult = models.TaskExecutionResult

type Agent struct {
	cfg       *Config
	ctx       context.Context
	registry  *registry.TaskDefinitionRegistry
	connector connector.EngineConnector
}

func NewAgent(ctx context.Context, cfg *Config) (*Agent, error) {
	reg := registry.NewTaskDefinitionRegistry()
	engineConnector, err := connector.NewEngineConnector(proto.AGENT_PROTOCOL_GRPC, cfg.EngineGrpcUrl)

	if err != nil {
		return nil, err
	}

	return &Agent{
		cfg:       cfg,
		ctx:       ctx,
		registry:  reg,
		connector: engineConnector,
	}, nil
}

func (a *Agent) Start() error {
	engineGrpcConnection, err := grpc.NewClient(a.cfg.EngineGrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to engine gRPC server: %v", err)
	}

	tick := ticker.NewEngineTicker(30, &connector.AgentInfo{
		Name:        a.cfg.Name,
		Version:     a.cfg.Version,
		Address:     a.cfg.GrpcAddress,
		Port:        a.cfg.GrpcPort,
		Definitions: a.registry.ToProto(),
	}, a.registry, a.connector)

	taskExecutor := executor.NewTaskExecutor(&executor.TaskExecutorConfig{
		MaxQueueSize:     100,
		MaxParallelTasks: 5,
	}, a.registry, engineGrpcConnection)

	grpcSrv := grpcserver.NewGrpcServer(
		a.cfg.GrpcAddress,
		a.cfg.GrpcPort,
		grpcserver.NewAgentService(),
	)

	go grpcSrv.Start()
	go taskExecutor.Start(a.ctx)
	go tick.Start()

	<-a.ctx.Done()
	log.Println("[agent] shutting down...")
	tick.Stop()

	if err := grpcSrv.Stop(); err != nil {
		log.Printf("grpc shutdown error: %v", err)
	}

	if err := a.connector.Close(); err != nil {
		log.Printf("engine connector close error: %v", err)
	}

	return nil
}

func (a *Agent) RegisterTaskDefinition(def models.TaskDefinition) {
	a.registry.Register(def)
}

func ReflectJsonSchema(v interface{}) *jsonschema.Schema {
	reflector := jsonschema.Reflector{}
	r, err := reflector.Reflect(v)
	if err != nil {
		return &jsonschema.Schema{}
	}
	return &r
}
