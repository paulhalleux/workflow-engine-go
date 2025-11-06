package persistence

import (
	"github.com/paulhalleux/workflow-engine-go/engine/internal/models"
	"gorm.io/gorm"
)

type Persistence struct {
	db                  *gorm.DB
	WorkflowDefinitions *WorkflowDefinitionsRepo
	WorkflowInstances   *WorkflowInstancesRepo
	StepInstances       *StepInstancesRepo
}

func NewPersistence(
	db *gorm.DB,
) *Persistence {
	workflowDefinitionsRepo := NewWorkflowDefinitionsRepo(db)
	workflowInstancesRepo := NewWorkflowInstancesRepo(db)
	stepInstancesRepo := NewStepInstancesRepo(db)

	return &Persistence{
		db: db,

		WorkflowDefinitions: workflowDefinitionsRepo,
		WorkflowInstances:   workflowInstancesRepo,
		StepInstances:       stepInstancesRepo,
	}
}

func (p *Persistence) Migrate() error {
	if err := p.db.AutoMigrate(
		&models.WorkflowDefinition{},
		&models.WorkflowInstance{},
		&models.StepInstance{},
	); err != nil {
		return err
	}
	return nil
}
