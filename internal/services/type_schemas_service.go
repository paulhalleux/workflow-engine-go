package services

import (
	"github.com/paulhalleux/workflow-engine-go/internal/dto"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/persistence"
)

type TypeSchemaService interface {
	GetAllTypeSchemas() ([]*models.TypeSchema, error)
	GetTypeSchemaById(id string) (*models.TypeSchema, error)
	CreateTypeSchema(req *dto.CreateTypeSchemaRequest) (*models.TypeSchema, error)
	DeleteTypeSchema(id string) error
}

type typeSchemaService struct {
	repo *persistence.TypeSchemasRepository
}

func NewTypeSchemaService(repo *persistence.TypeSchemasRepository) TypeSchemaService {
	return &typeSchemaService{repo}
}

func (s *typeSchemaService) GetAllTypeSchemas() ([]*models.TypeSchema, error) {
	typeSchemas, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	result := make([]*models.TypeSchema, len(typeSchemas))
	for i := range typeSchemas {
		result[i] = &typeSchemas[i]
	}
	return result, nil
}

func (s *typeSchemaService) GetTypeSchemaById(id string) (*models.TypeSchema, error) {
	return s.repo.GetById(id)
}

func (s *typeSchemaService) CreateTypeSchema(req *dto.CreateTypeSchemaRequest) (*models.TypeSchema, error) {
	typeSchema := &models.TypeSchema{
		Name:    req.Name,
		Version: req.Version,
		Schema:  *req.Schema,
	}
	err := s.repo.Create(typeSchema)
	if err != nil {
		return nil, err
	}
	return typeSchema, nil
}

func (s *typeSchemaService) DeleteTypeSchema(id string) error {
	return s.repo.Delete(id)
}
