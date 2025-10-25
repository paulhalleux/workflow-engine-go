package queue

import (
	"github.com/paulhalleux/workflow-engine-go/internal/models"
)

type WorkflowJob struct {
	WorkflowInstance *models.WorkflowInstance
}

type WorkflowQueue Queue[WorkflowJob]
