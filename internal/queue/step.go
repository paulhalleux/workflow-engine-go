package queue

import (
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/utils"
)

type StepJob struct {
	WorkflowFinishedCh chan<- WorkflowExecutionResult
	StepFinishedCh     chan string
	StepCounter        *utils.Counter
	StepDefinition     *models.WorkflowStep
	StepInstance       *models.StepInstance
	WorkflowDefinition *models.WorkflowDefinition
	WorkflowInstance   *models.WorkflowInstance
}

type StepQueue Queue[StepJob]
