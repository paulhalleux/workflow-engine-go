package persistence

import (
	"github.com/paulhalleux/workflow-engine-go/engine/internal/models"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
	"gorm.io/gorm"
)

type WorkflowDefinitionsRepo struct {
	db *gorm.DB
}

func NewWorkflowDefinitionsRepo(
	db *gorm.DB,
) *WorkflowDefinitionsRepo {
	return &WorkflowDefinitionsRepo{
		db: db,
	}
}

func (wr *WorkflowDefinitionsRepo) Create(definition *models.WorkflowDefinition) error {
	return wr.db.Create(definition).Error
}

func (wr *WorkflowDefinitionsRepo) GetByID(id string) (*models.WorkflowDefinition, error) {
	var definition models.WorkflowDefinition
	result := wr.db.Where("id = ?", id).First(&definition)
	if result.Error != nil {
		return nil, result.Error
	}
	return &definition, nil
}

func (wr *WorkflowDefinitionsRepo) GetAll(scopeFactory *utils.GormScopeFactory) ([]models.WorkflowDefinition, error) {
	var definitions []models.WorkflowDefinition
	result := utils.WithScope(wr.db, scopeFactory).Find(&definitions)
	if result.Error != nil {
		return nil, result.Error
	}
	return definitions, nil
}

func (wr *WorkflowDefinitionsRepo) Update(definition *models.WorkflowDefinition) error {
	return wr.db.Save(definition).Error
}

func (wr *WorkflowDefinitionsRepo) Delete(id string) error {
	return wr.db.Delete(&models.WorkflowDefinition{}, id).Error
}
