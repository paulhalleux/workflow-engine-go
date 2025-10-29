package persistence

import (
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"gorm.io/gorm"
)

type TypeSchemasRepository struct {
	db *gorm.DB
}

func NewTypeSchemasRepository(db *gorm.DB) *TypeSchemasRepository {
	return &TypeSchemasRepository{db}
}

func (r *TypeSchemasRepository) Create(typeSchema *models.TypeSchema) error {
	return r.db.Create(typeSchema).Error
}

func (r *TypeSchemasRepository) GetAll() ([]models.TypeSchema, error) {
	var typeSchemas []models.TypeSchema
	err := r.db.Find(&typeSchemas).Error
	return typeSchemas, err
}

func (r *TypeSchemasRepository) GetAllByName(name string) ([]models.TypeSchema, error) {
	var typeSchemas []models.TypeSchema
	err := r.db.Where("name = ?", name).Find(&typeSchemas).Error
	return typeSchemas, err
}

func (r *TypeSchemasRepository) GetByNameAndVersion(name string, version string) (*models.TypeSchema, error) {
	var typeSchema models.TypeSchema
	err := r.db.Where("name = ? AND version = ?", name, version).First(&typeSchema).Error
	if err != nil {
		return nil, err
	}
	return &typeSchema, nil
}

func (r *TypeSchemasRepository) GetById(id string) (*models.TypeSchema, error) {
	var typeSchema models.TypeSchema
	err := r.db.Where("id = ?", id).First(&typeSchema).Error
	if err != nil {
		return nil, err
	}
	return &typeSchema, nil
}

func (r *TypeSchemasRepository) Update(typeSchema *models.TypeSchema) error {
	return r.db.Save(typeSchema).Error
}

func (r *TypeSchemasRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.TypeSchema{}).Error
}
