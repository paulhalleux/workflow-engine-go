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
	Context                   context.Context
	CancelFunc                context.CancelFunc
	WorkflowService           services.WorkflowService
	WorkflowDefinitionService services.WorkflowDefinitionService
	WorkflowInstanceService   services.WorkflowInstanceService
	WorkflowQueue             queue.WorkflowQueue
	StepQueue                 queue.StepQueue
	WorkflowExecutor          *worker.WorkflowExecutor
	StepExecutor              *worker.StepExecutor
	WorkflowDefRepo           *persistence.WorkflowDefinitionsRepository
	WorkflowInstRepo          *persistence.WorkflowInstancesRepository
	Config                    *config.Config
}

func NewContainer(cfg *config.Config) *Container {
	db := createDatabaseConnection(cfg)
	ctx, cancel := context.WithCancel(context.Background())

	wfdRepo := persistence.NewWorkflowDefinitionsRepository(db)
	wfiRepo := persistence.NewWorkflowInstancesRepository(db)

	wfQueue := queue.NewMemoryWorkflowQueue(cfg.QueueBuffer)
	stepQueue := queue.NewMemoryStepQueue(cfg.QueueBuffer)

	wfService := services.NewWorkflowService(wfdRepo, wfiRepo, wfQueue)
	wfdService := services.NewWorkflowDefinitionService(wfdRepo)
	wfiService := services.NewWorkflowInstanceService(wfiRepo)

	wfExecutor := worker.NewWorkflowExecutor(wfiService, wfdService, stepQueue, wfQueue, cfg.MaxParallelWorkflows)
	stepExecutor := worker.NewStepExecutor(stepQueue, cfg.MaxParallelSteps)

	return &Container{
		WorkflowService:           wfService,
		WorkflowDefinitionService: wfdService,
		WorkflowInstanceService:   wfiService,
		WorkflowExecutor:          wfExecutor,
		WorkflowQueue:             wfQueue,
		StepQueue:                 stepQueue,
		StepExecutor:              stepExecutor,
		WorkflowDefRepo:           wfdRepo,
		WorkflowInstRepo:          wfiRepo,
		Context:                   ctx,
		CancelFunc:                cancel,
		Config:                    cfg,
	}
}

func createDatabaseConnection(cfg *config.Config) *gorm.DB {
	db, dbErr := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if dbErr != nil {
		log.Fatalf("failed to connect database: %v", dbErr)
	}
	return db
}
