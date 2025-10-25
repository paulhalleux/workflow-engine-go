package worker

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/paulhalleux/workflow-engine-go/internal/dto"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/queue"
	"github.com/paulhalleux/workflow-engine-go/internal/services"
	"gorm.io/datatypes"
)

type WorkflowExecutor struct {
	workflowInstanceService   services.WorkflowInstanceService
	workflowDefinitionService services.WorkflowDefinitionService
	stepQueue                 queue.StepQueue
	jobs                      <-chan queue.WorkflowJob
	sem                       chan struct{}
}

func NewWorkflowExecutor(
	workflowInstanceService services.WorkflowInstanceService,
	workflowDefinitionService services.WorkflowDefinitionService,
	stepQueue queue.StepQueue,
	workflowQueue queue.WorkflowQueue,
	maxParallelJobs int,
) *WorkflowExecutor {
	return &WorkflowExecutor{
		workflowInstanceService:   workflowInstanceService,
		workflowDefinitionService: workflowDefinitionService,
		stepQueue:                 stepQueue,
		jobs:                      workflowQueue.Dequeue(),
		sem:                       make(chan struct{}, maxParallelJobs),
	}
}

func (e *WorkflowExecutor) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case job := <-e.jobs:
				e.sem <- struct{}{}
				go func(j queue.WorkflowJob) {
					defer func() { <-e.sem }()
					e.handleJob(j)
				}(job)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (e *WorkflowExecutor) handleJob(job queue.WorkflowJob) {
	def, err := e.workflowDefinitionService.GetWorkflowDefinitionById(job.Instance.WorkflowDefinitionId.String())
	if err != nil {
		log.Printf("Error loading definition: %v", err)
		return
	}

	status := models.WorkflowInstanceStatusRunning
	now := time.Now()
	if _, err := e.workflowInstanceService.UpdateWorkflowInstance(
		job.Instance.Id.String(),
		&dto.UpdateWorkflowInstanceRequest{
			Status:    &status,
			StartedAt: &now,
		},
	); err != nil {
		log.Printf("Error updating instance status: %v", err)
		return
	}

	log.Printf("Workflow %s started with %d steps", def.Name, len(*def.Steps))

	for _, step := range *def.Steps {
		now := time.Now()
		inputData := step.Input.GetValueMap(nil, job.Instance.Input)
		taskInstance := models.StepInstance{
			Id:                 uuid.New(),
			StepId:             step.Id,
			Status:             models.StepInstanceStatusPending,
			WorkflowInstanceId: job.Instance.Id,
			Input:              inputData,
			Output:             datatypes.JSONMap{},
			Metadata:           datatypes.JSONMap{},
			CreatedAt:          now,
			UpdatedAt:          now,
		}

		stepJob := queue.StepJob{
			Step:               &step,
			Instance:           &taskInstance,
			WorkflowDefinition: def,
			WorkflowInstance:   job.Instance,
		}

		log.Printf("Enqueuing step %s of workflow %s", step.Id, def.Name)
		if err := e.stepQueue.Enqueue(stepJob); err != nil {
		}

		status = models.WorkflowInstanceStatusCompleted
		output := datatypes.JSON(`{"result": "success"}`)
		now = time.Now()
		if _, err := e.workflowInstanceService.UpdateWorkflowInstance(
			job.Instance.Id.String(),
			&dto.UpdateWorkflowInstanceRequest{
				Status:      &status,
				Output:      &output,
				CompletedAt: &now,
			},
		); err != nil {
			log.Printf("Error updating instance status: %v", err)
			return
		}
	}
}
