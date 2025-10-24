package worker

import (
	"context"
	"log"
	"time"

	"github.com/paulhalleux/workflow-engine-go/internal/dto"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/queue"
	"github.com/paulhalleux/workflow-engine-go/internal/services"
	"gorm.io/datatypes"
)

type WorkflowExecutor struct {
	workflowInstanceService   services.WorkflowInstanceService
	workflowDefinitionService services.WorkflowDefinitionService
	jobs                      <-chan queue.WorkflowJob
	sem                       chan struct{}
}

func NewWorkflowExecutor(
	workflowInstanceService services.WorkflowInstanceService,
	workflowDefinitionService services.WorkflowDefinitionService,
	q queue.WorkflowQueue,
	maxParallelJobs int,
) *WorkflowExecutor {
	return &WorkflowExecutor{
		workflowInstanceService:   workflowInstanceService,
		workflowDefinitionService: workflowDefinitionService,
		jobs:                      q.Dequeue(),
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
	log.Printf("Executing workflow instance %s", job.InstanceId)

	instance, err := e.workflowInstanceService.GetWorkflowInstanceById(job.InstanceId.String())
	if err != nil {
		log.Printf("Error loading instance: %v", err)
		return
	}

	def, err := e.workflowDefinitionService.GetWorkflowDefinitionById(instance.WorkflowDefinitionId.String())
	if err != nil {
		log.Printf("Error loading definition: %v", err)
		return
	}

	status := models.WorkflowInstanceStatusRunning
	if _, err := e.workflowInstanceService.UpdateWorkflowInstance(
		job.InstanceId.String(),
		&dto.UpdateWorkflowInstanceRequest{
			Status: &status,
		},
	); err != nil {
		log.Printf("Error updating instance status: %v", err)
		return
	}

	log.Printf("Workflow %s started with %d steps", def.Name, len(*def.Steps))
	time.Sleep(5 * time.Second)

	status = models.WorkflowInstanceStatusCompleted
	output := datatypes.JSON(`{"result": "success"}`)
	if _, err := e.workflowInstanceService.UpdateWorkflowInstance(
		job.InstanceId.String(),
		&dto.UpdateWorkflowInstanceRequest{
			Status: &status,
			Output: &output,
		},
	); err != nil {
		log.Printf("Error updating instance status: %v", err)
		return
	}
}
