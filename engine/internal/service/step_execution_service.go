package service

import (
	"sync"

	"github.com/paulhalleux/workflow-engine-go/engine/internal/errors"
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
	stepDefinitionID string,
	workflowInstance *models.WorkflowInstance,
	workflowDefinition *models.WorkflowDefinition,
	waitGroup *sync.WaitGroup,
) (*string, error) {
	stepDefinition, ok := workflowDefinition.GetStepByID(stepDefinitionID)
	if !ok {
		return nil, errors.ErrStepDefinitionNotFound
	}

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
		Next: func(stepId string) error {
			_, err := ws.StartStep(stepId, workflowInstance, workflowDefinition, waitGroup)
			if err != nil {
				return err
			}
			return nil
		},
		End: func() {
			waitGroup.Done()
		},
	}

	waitGroup.Add(1)
	err = ws.stepExecutor.Enqueue(exec)
	if err != nil {
		return nil, err
	}

	return &instance.ID, nil
}
