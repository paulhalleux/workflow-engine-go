package persistence

import (
	"github.com/paulhalleux/workflow-engine-go/engine/internal/models"
	"gorm.io/gorm"
)

type Persistence struct {
	db                  *gorm.DB
	WorkflowDefinitions *WorkflowDefinitionsRepo
	WorkflowInstances   *WorkflowInstancesRepo
}

func NewPersistence(
	db *gorm.DB,
) *Persistence {
	workflowDefinitionsRepo := NewWorkflowDefinitionsRepo(db)
	workflowInstancesRepo := NewWorkflowInstancesRepo(db)

	return &Persistence{
		db: db,

		WorkflowDefinitions: workflowDefinitionsRepo,
		WorkflowInstances:   workflowInstancesRepo,
	}
}

func (p *Persistence) Migrate() error {
	if err := p.db.AutoMigrate(
		&models.WorkflowDefinition{},
		&models.WorkflowInstance{},
	); err != nil {
		return err
	}
	return nil
}
