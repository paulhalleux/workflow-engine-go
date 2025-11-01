package persistence

import (
	"github.com/paulhalleux/workflow-engine-go/engine/internal/models"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
	"gorm.io/gorm"
)

type WorkflowInstancesRepo struct {
	db *gorm.DB
}

func NewWorkflowInstancesRepo(
	db *gorm.DB,
) *WorkflowInstancesRepo {
	return &WorkflowInstancesRepo{
		db: db,
	}
}

func (wr *WorkflowInstancesRepo) Create(definition *models.WorkflowInstance) error {
	return wr.db.Create(definition).Error
}

func (wr *WorkflowInstancesRepo) GetByID(id string) (*models.WorkflowInstance, error) {
	var definition models.WorkflowInstance
	result := wr.db.Where("id = ?", id).First(&definition)
	if result.Error != nil {
		return nil, result.Error
	}
	return &definition, nil
}

func (wr *WorkflowInstancesRepo) GetAll(scopeFactory *utils.GormScopeFactory) ([]models.WorkflowInstance, error) {
	var definitions []models.WorkflowInstance
	result := utils.WithScope(wr.db, scopeFactory).Find(&definitions)
	if result.Error != nil {
		return nil, result.Error
	}
	return definitions, nil
}

func (wr *WorkflowInstancesRepo) Update(definition *models.WorkflowInstance) error {
	return wr.db.Save(definition).Error
}

func (wr *WorkflowInstancesRepo) Delete(id string) error {
	return wr.db.Delete(&models.WorkflowInstance{}, id).Error
}
