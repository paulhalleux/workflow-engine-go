package queue

import (
	"github.com/paulhalleux/workflow-engine-go/internal/models"
)

type WorkflowJob struct {
	WorkflowInstance *models.WorkflowInstance
}

type WorkflowExecutionResult struct {
	Status models.WorkflowInstanceStatus
	Error  error
}

type WorkflowQueue Queue[WorkflowJob]
