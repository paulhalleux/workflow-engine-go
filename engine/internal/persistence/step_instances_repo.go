package persistence

import (
	"github.com/paulhalleux/workflow-engine-go/engine/internal/models"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
	"gorm.io/gorm"
)

type StepInstancesRepo struct {
	db *gorm.DB
}

func NewStepInstancesRepo(
	db *gorm.DB,
) *StepInstancesRepo {
	return &StepInstancesRepo{
		db: db,
	}
}

func (wr *StepInstancesRepo) Create(definition *models.StepInstance) error {
	return wr.db.Create(definition).Error
}

func (wr *StepInstancesRepo) GetByID(id string) (*models.StepInstance, error) {
	var definition models.StepInstance
	result := wr.db.Where("id = ?", id).First(&definition)
	if result.Error != nil {
		return nil, result.Error
	}
	return &definition, nil
}

func (wr *StepInstancesRepo) GetAll(scopeFactory *utils.GormScopeFactory) ([]models.StepInstance, error) {
	var definitions []models.StepInstance
	result := utils.WithScope(wr.db, scopeFactory).Find(&definitions)
	if result.Error != nil {
		return nil, result.Error
	}
	return definitions, nil
}

func (wr *StepInstancesRepo) Update(definition *models.StepInstance) error {
	return wr.db.Where("id = ?", definition.ID).Updates(definition).Error
}

func (wr *StepInstancesRepo) Delete(id string) error {
	return wr.db.Delete(&models.StepInstance{}, id).Error
}
