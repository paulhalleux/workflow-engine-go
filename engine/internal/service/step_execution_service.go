package service

import (
	"github.com/paulhalleux/workflow-engine-go/engine/internal/models"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/persistence"
)

type StepExecutionService struct {
	persistence  *persistence.Persistence
	stepExecutor *StepExecutor
}

func NewStepExecutionService(
	persistence *persistence.Persistence,
	stepExecutor *StepExecutor,
) *StepExecutionService {
	return &StepExecutionService{
		persistence:  persistence,
		stepExecutor: stepExecutor,
	}
}

func (ws *StepExecutionService) StartStep(
	stepDefinition *models.WorkflowStepDefinition,
	workflowInstance *models.WorkflowInstance,
) (*string, error) {
	instance := stepDefinition.NewInstance(
		workflowInstance.WorkflowDefinitionID,
		workflowInstance.ID,
		stepDefinition.Parameters.ToResolved(
			workflowInstance.Input.ToMap(),
			&map[string]map[string]interface{}{},
		),
	)

	err := ws.persistence.StepInstances.Create(instance)
	if err != nil {
		return nil, err
	}

	exec := &StepExecution{
		StepInstanceID:     instance.ID,
		WorkflowInstanceID: workflowInstance.ID,
		StepDef:            stepDefinition,
		Input:              instance.Input.ToMap(),
	}

	err = ws.stepExecutor.Enqueue(exec)
	if err != nil {
		return nil, err
	}

	return &instance.ID, nil
}
