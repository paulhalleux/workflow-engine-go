package container

import (
	"context"
	"log"

	"github.com/paulhalleux/workflow-engine-go/internal/config"
	"github.com/paulhalleux/workflow-engine-go/internal/grpcapi"
	"github.com/paulhalleux/workflow-engine-go/internal/persistence"
	"github.com/paulhalleux/workflow-engine-go/internal/queue"
	"github.com/paulhalleux/workflow-engine-go/internal/services"
	"github.com/paulhalleux/workflow-engine-go/internal/worker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Container struct {
	Context    context.Context
	CancelFunc context.CancelFunc

	Config *config.Config

	AgentTaskExecutionChan chan grpcapi.TaskExecutionResult

	WorkflowService           services.WorkflowService
	WorkflowDefinitionService services.WorkflowDefinitionService
	WorkflowInstanceService   services.WorkflowInstanceService
	TypeSchemaService         services.TypeSchemaService
	StepInstanceService       services.StepInstanceService

	WorkflowExecutor *worker.WorkflowExecutor
	StepExecutor     *worker.StepExecutor

	WorkflowDefRepo  *persistence.WorkflowDefinitionsRepository
	WorkflowInstRepo *persistence.WorkflowInstancesRepository
	TypeSchemaRepo   *persistence.TypeSchemasRepository
	StepInstanceRepo *persistence.StepInstancesRepository
}

func NewContainer(cfg *config.Config) *Container {
	db := createDatabaseConnection(cfg)
	ctx, cancel := context.WithCancel(context.Background())

	workflowDefinitionsRepo := persistence.NewWorkflowDefinitionsRepository(db)
	workflowInstancesRepo := persistence.NewWorkflowInstancesRepository(db)
	typeSchemaRepo := persistence.NewTypeSchemasRepository(db)
	stepInstancesRepo := persistence.NewStepInstancesRepository(db)

	workflowQueue := queue.NewMemoryQueue[queue.WorkflowJob](cfg.QueueBuffer)
	stepQueue := queue.NewMemoryQueue[queue.StepJob](cfg.QueueBuffer)

	workflowService := services.NewWorkflowService(workflowDefinitionsRepo, workflowInstancesRepo, workflowQueue)
	workflowDefinitionService := services.NewWorkflowDefinitionService(workflowDefinitionsRepo)
	workflowInstanceService := services.NewWorkflowInstanceService(workflowInstancesRepo)
	typeSchemaService := services.NewTypeSchemaService(typeSchemaRepo)
	stepInstanceService := services.NewStepInstanceService(stepInstancesRepo)

	workflowExecutor := worker.NewWorkflowExecutor(stepInstanceService, workflowInstanceService, workflowDefinitionService, workflowQueue, stepQueue, cfg.MaxParallelWorkflows)
	stepExecutor := worker.NewStepExecutor(stepInstanceService, stepQueue, cfg.MaxParallelSteps)

	agentTaskExecutionChan := make(chan grpcapi.TaskExecutionResult, cfg.QueueBuffer)

	return &Container{
		WorkflowService:           workflowService,
		WorkflowDefinitionService: workflowDefinitionService,
		WorkflowInstanceService:   workflowInstanceService,
		TypeSchemaService:         typeSchemaService,
		StepInstanceService:       stepInstanceService,

		WorkflowExecutor: workflowExecutor,
		StepExecutor:     stepExecutor,

		WorkflowDefRepo:  workflowDefinitionsRepo,
		WorkflowInstRepo: workflowInstancesRepo,
		TypeSchemaRepo:   typeSchemaRepo,
		StepInstanceRepo: stepInstancesRepo,

		Context:    ctx,
		CancelFunc: cancel,
		Config:     cfg,

		AgentTaskExecutionChan: agentTaskExecutionChan,
	}
}

func createDatabaseConnection(cfg *config.Config) *gorm.DB {
	db, dbErr := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if dbErr != nil {
		log.Fatalf("failed to connect database: %v", dbErr)
	}
	return db
}
