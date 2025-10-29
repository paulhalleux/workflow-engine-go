package persistence

import (
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"gorm.io/gorm"
)

type StepInstancesRepository struct {
	db *gorm.DB
}

func NewStepInstancesRepository(db *gorm.DB) *StepInstancesRepository {
	return &StepInstancesRepository{db}
}

func (r *StepInstancesRepository) Create(wf *models.StepInstance) error {
	return r.db.Create(wf).Error
}

func (r *StepInstancesRepository) GetAll() ([]models.StepInstance, error) {
	var wfs []models.StepInstance
	err := r.db.Find(&wfs).Error
	return wfs, err
}

func (r *StepInstancesRepository) GetById(id string) (*models.StepInstance, error) {
	var wf models.StepInstance
	err := r.db.First(&wf, "id = ?", id).Error
	return &wf, err
}

func (r *StepInstancesRepository) Update(wf *models.StepInstance) error {
	return r.db.Save(wf).Error
}

func (r *StepInstancesRepository) Delete(id string) error {
	return r.db.Delete(&models.StepInstance{}, "id = ?", id).Error
}
