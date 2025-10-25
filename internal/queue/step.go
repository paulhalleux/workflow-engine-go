package queue

import (
	"github.com/paulhalleux/workflow-engine-go/internal/models"
)

type StepJob struct {
	WorkflowFinishedCh chan<- struct{}
	StepDefinition     *models.WorkflowStep
	StepInstance       *models.StepInstance
	WorkflowDefinition *models.WorkflowDefinition
	WorkflowInstance   *models.WorkflowInstance
}

type StepQueue Queue[StepJob]
