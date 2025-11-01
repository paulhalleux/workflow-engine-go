package service

import (
	"github.com/paulhalleux/workflow-engine-go/engine/internal/models"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/persistence"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
)

type WorkflowDefinitionsService struct {
	persistence *persistence.Persistence
}

func NewWorkflowDefinitionsService(
	persistence *persistence.Persistence,
) *WorkflowDefinitionsService {
	return &WorkflowDefinitionsService{
		persistence: persistence,
	}
}

func (ws *WorkflowDefinitionsService) GetAll(scopeFactory *utils.GormScopeFactory) ([]models.WorkflowDefinition, error) {
	return ws.persistence.WorkflowDefinitions.GetAll(scopeFactory)
}

func (ws *WorkflowDefinitionsService) GetByID(id string) (*models.WorkflowDefinition, error) {
	return ws.persistence.WorkflowDefinitions.GetByID(id)
}

func (ws *WorkflowDefinitionsService) Create(definition *models.WorkflowDefinition) error {
	return ws.persistence.WorkflowDefinitions.Create(definition)
}

func (ws *WorkflowDefinitionsService) Update(definition *models.WorkflowDefinition) error {
	return ws.persistence.WorkflowDefinitions.Update(definition)
}

func (ws *WorkflowDefinitionsService) Delete(id string) error {
	return ws.persistence.WorkflowDefinitions.Delete(id)
}

func (ws *WorkflowDefinitionsService) Enable(id string) error {
	definition, err := ws.GetByID(id)
	if err != nil {
		return err
	}
	definition.IsEnabled = true
	return ws.Update(definition)
}

func (ws *WorkflowDefinitionsService) Disable(id string) error {
	definition, err := ws.GetByID(id)
	if err != nil {
		return err
	}
	definition.IsEnabled = false
	return ws.Update(definition)
}

func (ws *WorkflowDefinitionsService) Publish(id string) error {
	definition, err := ws.GetByID(id)
	if err != nil {
		return err
	}
	definition.IsDraft = false
	return ws.Update(definition)
}
