package container

import (
	"context"
	"log"

	"github.com/paulhalleux/workflow-engine-go/internal/config"
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

	WorkflowService           services.WorkflowService
	WorkflowDefinitionService services.WorkflowDefinitionService
	WorkflowInstanceService   services.WorkflowInstanceService
	TypeSchemaService         services.TypeSchemaService

	WorkflowExecutor *worker.WorkflowExecutor
	StepExecutor     *worker.StepExecutor

	WorkflowDefRepo  *persistence.WorkflowDefinitionsRepository
	WorkflowInstRepo *persistence.WorkflowInstancesRepository
	TypeSchemaRepo   *persistence.TypeSchemasRepository
}

func NewContainer(cfg *config.Config) *Container {
	db := createDatabaseConnection(cfg)
	ctx, cancel := context.WithCancel(context.Background())

	workflowDefinitionsRepo := persistence.NewWorkflowDefinitionsRepository(db)
	workflowInstancesRepo := persistence.NewWorkflowInstancesRepository(db)
	typeSchemaRepo := persistence.NewTypeSchemasRepository(db)

	workflowQueue := queue.NewMemoryQueue[queue.WorkflowJob](cfg.QueueBuffer)
	stepQueue := queue.NewMemoryQueue[queue.StepJob](cfg.QueueBuffer)

	workflowService := services.NewWorkflowService(workflowDefinitionsRepo, workflowInstancesRepo, workflowQueue)
	workflowDefinitionService := services.NewWorkflowDefinitionService(workflowDefinitionsRepo)
	workflowInstanceService := services.NewWorkflowInstanceService(workflowInstancesRepo)
	typeSchemaService := services.NewTypeSchemaService(typeSchemaRepo)

	workflowExecutor := worker.NewWorkflowExecutor(workflowInstanceService, workflowDefinitionService, workflowQueue, stepQueue, cfg.MaxParallelWorkflows)
	stepExecutor := worker.NewStepExecutor(stepQueue, cfg.MaxParallelSteps)

	return &Container{
		WorkflowService:           workflowService,
		WorkflowDefinitionService: workflowDefinitionService,
		WorkflowInstanceService:   workflowInstanceService,
		TypeSchemaService:         typeSchemaService,

		WorkflowExecutor: workflowExecutor,
		StepExecutor:     stepExecutor,

		WorkflowDefRepo:  workflowDefinitionsRepo,
		WorkflowInstRepo: workflowInstancesRepo,
		TypeSchemaRepo:   typeSchemaRepo,

		Context:    ctx,
		CancelFunc: cancel,
		Config:     cfg,
	}
}

func createDatabaseConnection(cfg *config.Config) *gorm.DB {
	db, dbErr := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if dbErr != nil {
		log.Fatalf("failed to connect database: %v", dbErr)
	}
	return db
}
