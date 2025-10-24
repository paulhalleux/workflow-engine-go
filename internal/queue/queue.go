package queue

import (
	"github.com/paulhalleux/workflow-engine-go/internal/models"
)

type WorkflowJob struct {
	Instance *models.WorkflowInstance
}

type StepJob struct {
	Instance           *models.StepInstance
	WorkflowInstance   *models.WorkflowInstance
	WorkflowDefinition *models.WorkflowDefinition
	Step               *models.WorkflowStep
}

type WorkflowQueue interface {
	Enqueue(job WorkflowJob) error
	Dequeue() <-chan WorkflowJob
}

type StepQueue interface {
	Enqueue(job StepJob) error
	Dequeue() <-chan StepJob
}
