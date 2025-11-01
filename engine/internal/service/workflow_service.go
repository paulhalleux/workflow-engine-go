package service

import (
	"github.com/paulhalleux/workflow-engine-go/engine/internal/models"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/persistence"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
)

type WorkflowService struct {
	workflowDefinitionsService *WorkflowDefinitionsService
	persistence                *persistence.Persistence
}

func NewWorkflowService(
	workflowDefinitionsService *WorkflowDefinitionsService,
	persistence *persistence.Persistence,
) *WorkflowService {
	return &WorkflowService{
		workflowDefinitionsService: workflowDefinitionsService,
		persistence:                persistence,
	}
}

func (ws *WorkflowService) StartWorkflow(definitionID string, parameters map[string]interface{}) (*string, error) {
	def, err := ws.workflowDefinitionsService.GetByID(definitionID)
	if err != nil {
		return nil, err
	}

	instance := def.NewInstance(&parameters)
	err = ws.persistence.WorkflowInstances.Create(instance)
	if err != nil {
		return nil, err
	}

	return &instance.ID, nil
}

func (ws *WorkflowService) GetAll(scopeFactory *utils.GormScopeFactory) ([]models.WorkflowInstance, error) {
	instances, err := ws.persistence.WorkflowInstances.GetAll(scopeFactory)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func (ws *WorkflowService) GetByID(id string) (*models.WorkflowInstance, error) {
	instance, err := ws.persistence.WorkflowInstances.GetByID(id)
	if err != nil {
		return nil, err
	}
	return instance, nil
}
