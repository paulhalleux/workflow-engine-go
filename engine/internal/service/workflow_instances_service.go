package service

import (
	"github.com/paulhalleux/workflow-engine-go/engine/internal/models"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/persistence"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
)

type WorkflowInstanceService struct {
	persistence *persistence.Persistence
}

func NewWorkflowInstanceService(
	persistence *persistence.Persistence,
) *WorkflowInstanceService {
	return &WorkflowInstanceService{
		persistence: persistence,
	}
}

func (ws *WorkflowInstanceService) GetAll(scopeFactory *utils.GormScopeFactory) ([]models.WorkflowInstance, error) {
	instances, err := ws.persistence.WorkflowInstances.GetAll(scopeFactory)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (ws *WorkflowInstanceService) GetByID(id string) (*models.WorkflowInstance, error) {
	instance, err := ws.persistence.WorkflowInstances.GetByID(id)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (ws *WorkflowInstanceService) Update(instance *models.WorkflowInstance) error {
	return ws.persistence.WorkflowInstances.Update(instance)
}
