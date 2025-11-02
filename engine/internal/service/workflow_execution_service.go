package service

import (
	"github.com/paulhalleux/workflow-engine-go/engine/internal/persistence"
)

type WorkflowExecutionService struct {
	workflowDefinitionsService *WorkflowDefinitionsService
	persistence                *persistence.Persistence
	workflowExecutor           *WorkflowExecutor
}

func NewWorkflowExecutionService(
	workflowDefinitionsService *WorkflowDefinitionsService,
	persistence *persistence.Persistence,
	workflowExecutor *WorkflowExecutor,
) *WorkflowExecutionService {
	return &WorkflowExecutionService{
		workflowDefinitionsService: workflowDefinitionsService,
		persistence:                persistence,
		workflowExecutor:           workflowExecutor,
	}
}

func (ws *WorkflowExecutionService) StartWorkflow(definitionID string, parameters map[string]interface{}) (*string, error) {
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
