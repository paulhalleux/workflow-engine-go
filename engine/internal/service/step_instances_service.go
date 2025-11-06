package service

import (
	"github.com/paulhalleux/workflow-engine-go/engine/internal/models"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/persistence"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
)

type StepInstanceService struct {
	persistence *persistence.Persistence
}

func NewStepInstanceService(
	persistence *persistence.Persistence,
) *StepInstanceService {
	return &StepInstanceService{
		persistence: persistence,
	}
}

func (ws *StepInstanceService) GetAll(scopeFactory *utils.GormScopeFactory) ([]models.StepInstance, error) {
	instances, err := ws.persistence.StepInstances.GetAll(scopeFactory)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (ws *StepInstanceService) GetByID(id string) (*models.StepInstance, error) {
	instance, err := ws.persistence.StepInstances.GetByID(id)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (ws *StepInstanceService) Update(instance *models.StepInstance) error {
	return ws.persistence.StepInstances.Update(instance)
}
