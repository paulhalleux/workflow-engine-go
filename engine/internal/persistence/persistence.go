package persistence

import (
	"github.com/paulhalleux/workflow-engine-go/engine/internal/models"
	"gorm.io/gorm"
)

type Persistence struct {
	db                  *gorm.DB
	WorkflowDefinitions *WorkflowDefinitionsRepo
}

func NewPersistence(
	db *gorm.DB,
) *Persistence {
	workflowDefinitionsRepo := NewWorkflowDefinitionsRepo(db)
	return &Persistence{
		db: db,

		WorkflowDefinitions: workflowDefinitionsRepo,
	}
}

func (p *Persistence) Migrate() error {
	if err := p.db.AutoMigrate(
		&models.WorkflowDefinition{},
	); err != nil {
		return err
	}
	return nil
}
