package persistance

import (
	"github.com/paulhalleux/workflow-engine-go/engine-new/internal/models"
	"github.com/paulhalleux/workflow-engine-go/utils/expr"
	"github.com/paulhalleux/workflow-engine-go/utils/pagination"
	"gorm.io/gorm"
)

type WorkflowDefinitionRepository interface {
	GetAll(pagination pagination.Pagination) (*pagination.PaginatedResult[models.WorkflowDefinition], error)
	Search(expression expr.Expression, pagination pagination.Pagination) (*pagination.PaginatedResult[models.WorkflowDefinition], error)
	GetByID(id string) (*models.WorkflowDefinition, error)
	Create(definition *models.WorkflowDefinition) (*models.WorkflowDefinition, error)
	Update(definition *models.WorkflowDefinition) (*models.WorkflowDefinition, error)
	Delete(id string) error
	Enable(id string) error
	Disable(id string) error
}

type workflowDefinitionRepository struct {
	db *gorm.DB
}

func NewWorkflowDefinitionRepository(
	db *gorm.DB,
) WorkflowDefinitionRepository {
	return &workflowDefinitionRepository{
		db: db,
	}
}

func (r *workflowDefinitionRepository) GetAll(
	pg pagination.Pagination,
) (*pagination.PaginatedResult[models.WorkflowDefinition], error) {
	var totalCount int64
	if err := r.db.Model(&models.WorkflowDefinition{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	definitions := make([]models.WorkflowDefinition, 0)
	result := pg.ToGorm(r.db).Find(&definitions)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pagination.PaginatedResult[models.WorkflowDefinition]{
		TotalCount: totalCount,
		Items:      definitions,
	}, nil
}

func (r *workflowDefinitionRepository) Search(expression expr.Expression, pg pagination.Pagination) (*pagination.PaginatedResult[models.WorkflowDefinition], error) {
	var totalCount int64
	if err := expression.ToGorm(r.db.Model(&models.WorkflowDefinition{}), false).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	definitions := make([]models.WorkflowDefinition, 0)
	result := pg.ToGorm(expression.ToGorm(r.db, false)).Find(&definitions)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pagination.PaginatedResult[models.WorkflowDefinition]{
		TotalCount: totalCount,
		Items:      definitions,
	}, nil
}

func (r *workflowDefinitionRepository) GetByID(id string) (*models.WorkflowDefinition, error) {
	definition := &models.WorkflowDefinition{}
	result := r.db.First(definition, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return definition, nil
}

func (r *workflowDefinitionRepository) Create(definition *models.WorkflowDefinition) (*models.WorkflowDefinition, error) {
	result := r.db.Create(definition)
	if result.Error != nil {
		return nil, result.Error
	}
	return definition, nil
}

func (r *workflowDefinitionRepository) Update(definition *models.WorkflowDefinition) (*models.WorkflowDefinition, error) {
	result := r.db.Save(definition)
	if result.Error != nil {
		return nil, result.Error
	}
	return definition, nil
}

func (r *workflowDefinitionRepository) Delete(id string) error {
	result := r.db.Delete(&models.WorkflowDefinition{}, "id = ?", id)
	return result.Error
}

func (r *workflowDefinitionRepository) Enable(id string) error {
	result := r.db.Model(&models.WorkflowDefinition{}).Where("id = ?", id).Update("is_enabled", true)
	return result.Error
}

func (r *workflowDefinitionRepository) Disable(id string) error {
	result := r.db.Model(&models.WorkflowDefinition{}).Where("id = ?", id).Update("is_enabled", false)
	return result.Error
}
