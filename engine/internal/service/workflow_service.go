package service

import (
	"github.com/paulhalleux/workflow-engine-go/engine/internal/persistence"
)

type WorkflowService struct {
	workflowDefinitionsService *WorkflowDefinitionsService
	persistence                *persistence.Persistence
	workflowExecutor           *WorkflowExecutor
}

func NewWorkflowService(
	workflowDefinitionsService *WorkflowDefinitionsService,
	persistence *persistence.Persistence,
	workflowExecutor *WorkflowExecutor,
) *WorkflowService {
	return &WorkflowService{
		workflowDefinitionsService: workflowDefinitionsService,
		persistence:                persistence,
		workflowExecutor:           workflowExecutor,
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

	exec := &WorkflowExecution{
		WorkflowInstanceID: instance.ID,
	}
	err = ws.workflowExecutor.Enqueue(exec)
	if err != nil {
		return nil, err
	}

	return &instance.ID, nil
}
