package persistence

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/paulhalleux/workflow-engine-go/internal/api/dto"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"gorm.io/gorm"
)

type WorkflowDefinitionsRepository struct {
	db *gorm.DB
}

func NewWorkflowRepository(db *gorm.DB) *WorkflowDefinitionsRepository {
	return &WorkflowDefinitionsRepository{db}
}

func (r *WorkflowDefinitionsRepository) Create(wf *models.WorkflowDefinition) error {
	return r.db.Create(wf).Error
}

func (r *WorkflowDefinitionsRepository) GetAll() ([]models.WorkflowDefinition, error) {
	var wfs []models.WorkflowDefinition
	err := r.db.Find(&wfs).Error
	return wfs, err
}

func (r *WorkflowDefinitionsRepository) GetById(id string) (*models.WorkflowDefinition, error) {
	var wf models.WorkflowDefinition
	err := r.db.First(&wf, "id = ?", id).Error
	return &wf, err
}

func (r *WorkflowDefinitionsRepository) Update(wf *models.WorkflowDefinition) error {
	return r.db.Save(wf).Error
}

func (r *WorkflowDefinitionsRepository) Delete(id string) error {
	return r.db.Delete(&models.WorkflowDefinition{}, "id = ?", id).Error
}

func (r *WorkflowDefinitionsRepository) Search(req *dto.SearchWorkflowDefinitionsRequest) ([]models.WorkflowDefinition, error) {
	var wfs []models.WorkflowDefinition
	query := r.db.Model(&models.WorkflowDefinition{})

	if req.Name != nil {
		query = query.Where("name = ?", *req.Name)
	}
	if req.Version != nil {
		query = query.Where("version = ?", *req.Version)
	}
	if req.IsEnabled != nil {
		query = query.Where("is_enabled = ?", *req.IsEnabled)
	}
	if req.MajorVersion != nil {
		query = query.Where("version LIKE ?", fmt.Sprintf("%d.%%", *req.MajorVersion))
	}
	if req.MinorVersion != nil {
		query = query.Where("version LIKE ?", fmt.Sprintf("%%.%d.%%", *req.MinorVersion))
	}
	if req.PatchVersion != nil {
		query = query.Where("version LIKE ?", fmt.Sprintf("%%.%%.%d", *req.PatchVersion))
	}
	if req.IsRelease != nil {
		if *req.IsRelease {
			query = query.Where("version NOT LIKE ?", "%draft%")
		} else {
			query = query.Where("version LIKE ?", "%draft%")
		}
	}

	err := query.Find(&wfs).Error
	return wfs, err
}

func (r *WorkflowDefinitionsRepository) IsDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	errors.As(err, &pgErr)
	return pgErr.Code == "23505"
}
