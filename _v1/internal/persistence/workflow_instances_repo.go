package persistence

import (
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"gorm.io/gorm"
)

type WorkflowInstancesRepository struct {
	db *gorm.DB
}

func NewWorkflowInstancesRepository(db *gorm.DB) *WorkflowInstancesRepository {
	return &WorkflowInstancesRepository{db}
}

func (r *WorkflowInstancesRepository) Create(wf *models.WorkflowInstance) error {
	return r.db.Create(wf).Error
}

func (r *WorkflowInstancesRepository) GetAll() ([]models.WorkflowInstance, error) {
	var wfs []models.WorkflowInstance
	err := r.db.Find(&wfs).Error
	return wfs, err
}

func (r *WorkflowInstancesRepository) GetById(id string) (*models.WorkflowInstance, error) {
	var wf models.WorkflowInstance
	err := r.db.First(&wf, "id = ?", id).Error
	return &wf, err
}

func (r *WorkflowInstancesRepository) Update(wf *models.WorkflowInstance) error {
	return r.db.Save(wf).Error
}

func (r *WorkflowInstancesRepository) Delete(id string) error {
	return r.db.Delete(&models.WorkflowInstance{}, "id = ?", id).Error
}
