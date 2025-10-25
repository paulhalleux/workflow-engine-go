package services

import (
	"github.com/google/uuid"
	"github.com/paulhalleux/workflow-engine-go/internal/dto"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/persistence"
)

type StepInstanceService interface {
	CreateStepInstance(stepInstance *models.StepInstance) (*models.StepInstance, error)
	UpdateStepInstance(id uuid.UUID, stepInstance *dto.UpdateStepInstanceRequest) (*models.StepInstance, error)
}

type stepInstanceService struct {
	repo *persistence.StepInstancesRepository
}

func NewStepInstanceService(repo *persistence.StepInstancesRepository) StepInstanceService {
	return &stepInstanceService{repo}
}

func (s *stepInstanceService) CreateStepInstance(stepInstance *models.StepInstance) (*models.StepInstance, error) {
	err := s.repo.Create(stepInstance)
	if err != nil {
		return nil, err
	}
	return stepInstance, nil
}

func (s *stepInstanceService) UpdateStepInstance(id uuid.UUID, stepInstance *dto.UpdateStepInstanceRequest) (*models.StepInstance, error) {
	existingStepInstance, err := s.repo.GetById(id.String())
	if err != nil {
		return nil, err
	}
	updatePartialStepInstance(existingStepInstance, stepInstance)
	err = s.repo.Update(existingStepInstance)
	if err != nil {
		return nil, err
	}
	return existingStepInstance, nil
}

func updatePartialStepInstance(stepInstance *models.StepInstance, dto *dto.UpdateStepInstanceRequest) {
	if dto.Status != nil {
		stepInstance.Status = *dto.Status
	}
	if dto.Output != nil {
		stepInstance.Output = *dto.Output
	}
	if dto.Metadata != nil {
		stepInstance.Metadata = *dto.Metadata
	}
	if dto.StartedAt != nil {
		stepInstance.StartedAt = dto.StartedAt
	}
	if dto.CompletedAt != nil {
		stepInstance.CompletedAt = dto.CompletedAt
	}
}
