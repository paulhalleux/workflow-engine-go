package engine

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/paulhalleux/workflow-engine-go/engine/internal"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/grpcapi"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/httpapi"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/persistence"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/service"
	"github.com/paulhalleux/workflow-engine-go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type WorkflowEngineConfig = internal.WorkflowEngineConfig

type Engine struct {
	Config  *internal.WorkflowEngineConfig
	Context context.Context

	persistence   *persistence.Persistence
	agentRegistry *internal.AgentRegistry
	db            *gorm.DB

	workflowDefinitionService *service.WorkflowDefinitionsService
	workflowInstanceService   *service.WorkflowInstanceService
	workflowExecutionService  *service.WorkflowExecutionService
	stepExecutionService      *service.StepExecutionService

	workflowExecutor *service.WorkflowExecutor
	stepExecutor     *service.StepExecutor
}

func NewEngine(
	config *internal.WorkflowEngineConfig,
) *Engine {
	ctx := context.Background()

	agentRegistry := internal.NewAgentRegistry()
	database, err := createDatabase(config)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	pers := persistence.NewPersistence(database)
	if err := pers.Migrate(); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	workflowDefinitionService := service.NewWorkflowDefinitionsService(pers)
	workflowInstanceService := service.NewWorkflowInstanceService(workflowDefinitionService, pers)

	stepExecutor := service.NewStepExecutor(config)
	stepExecutionService := service.NewStepExecutionService(
		pers,
		stepExecutor,
	)

	workflowExecutor := service.NewWorkflowExecutor(config, workflowDefinitionService, workflowInstanceService, stepExecutionService)
	workflowExecutionService := service.NewWorkflowExecutionService(
		workflowDefinitionService,
		pers,
		workflowExecutor,
	)

	return &Engine{
		Config:  config,
		Context: ctx,

		persistence:   pers,
		agentRegistry: agentRegistry,
		db:            database,

		workflowDefinitionService: workflowDefinitionService,
		workflowInstanceService:   workflowInstanceService,
		workflowExecutionService:  workflowExecutionService,
		stepExecutionService:      stepExecutionService,

		workflowExecutor: workflowExecutor,
		stepExecutor:     stepExecutor,
	}
}

func (e *Engine) Start() {
	log.Printf("[Engine] Starting workflow engine...")

	go e.startGrpcServer()
	go e.startHttpServer()
	go e.workflowExecutor.Start(e.Context)

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
	proto.RegisterEngineServiceServer(grpcServer, grpcapi.NewEngineServiceServer(e.agentRegistry, e.workflowExecutionService))
	proto.RegisterTaskServiceServer(grpcServer, grpcapi.NewTaskServiceServer())

	// Start the gRPC server
	log.Printf("[Engine] gRPC server running on %s", lis.Addr().String())
	if serveErr := grpcServer.Serve(lis); serveErr != nil {
		log.Fatalf("failed to serve gRPC server: %v", serveErr)
	}
}

func (e *Engine) startHttpServer() {
	r := gin.Default()
	api := r.Group("/api")

	workflowDefinitionHandlers := httpapi.NewWorkflowDefinitionsHandlers(e.workflowDefinitionService)
	workflowInstanceHandlers := httpapi.NewWorkflowInstancesHandlers(e.workflowInstanceService)

	// Define HTTP routes and handlers here
	workflowDefinitionHandlers.Register(api)
	workflowInstanceHandlers.Register(api)

	addr := joinHostPort(e.Config.HttpAddress, e.Config.HttpPort)
	log.Printf("[Engine] HTTP server running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to run HTTP server: %v", err)
	}
}

func createDatabase(config *WorkflowEngineConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DbHost, config.DbPort, config.DbUser, config.DbPassword, config.DbName, config.DbSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func joinHostPort(host *string, port string) string {
	if host == nil || *host == "" {
		return fmt.Sprintf(":%s", port)
	}
	return fmt.Sprintf("%s:%s", *host, port)
}
