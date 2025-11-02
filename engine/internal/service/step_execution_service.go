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

func (ws *StepExecutionService) StartStep(stepDefinition *models.WorkflowStepDefinition, workflowDefinition *models.WorkflowDefinition) (*string, error) {
	return nil, nil
}
