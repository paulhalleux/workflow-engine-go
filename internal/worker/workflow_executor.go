package worker

import (
	"context"
	"log"
	"time"

	"github.com/paulhalleux/workflow-engine-go/internal/dto"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/queue"
	"github.com/paulhalleux/workflow-engine-go/internal/services"
)

type WorkflowExecutor struct {
	workflowInstancesService  services.WorkflowInstanceService
	workflowDefinitionService services.WorkflowDefinitionService
	workflowQueue             queue.WorkflowQueue
	stepQueue                 queue.StepQueue
	sem                       chan struct{}
}

func NewWorkflowExecutor(
	workflowInstancesService services.WorkflowInstanceService,
	workflowDefinitionService services.WorkflowDefinitionService,
	workflowQueue queue.WorkflowQueue,
	stepQueue queue.StepQueue,
	maxConcurrentWorkflows int,
) *WorkflowExecutor {
	return &WorkflowExecutor{
		workflowInstancesService:  workflowInstancesService,
		workflowDefinitionService: workflowDefinitionService,
		workflowQueue:             workflowQueue,
		stepQueue:                 stepQueue,
		sem:                       make(chan struct{}, maxConcurrentWorkflows),
	}
}

func (e *WorkflowExecutor) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case job := <-e.workflowQueue.Dequeue():
				e.sem <- struct{}{}
				go func(j queue.WorkflowJob) {
					defer func() { <-e.sem }()
					e.handle(j)
				}(job)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (e *WorkflowExecutor) handle(job queue.WorkflowJob) {
	workflowDefinition, err := e.workflowDefinitionService.GetWorkflowDefinitionById(job.WorkflowInstance.WorkflowDefinitionId.String())
	if err != nil {
		log.Printf("Error loading workflow definition: %v", err)
		return
	}

	newStatus := models.WorkflowInstanceStatusRunning
	startTime := time.Now()

	_, err = e.workflowInstancesService.UpdateWorkflowInstance(
		job.WorkflowInstance.Id.String(),
		&dto.UpdateWorkflowInstanceRequest{
			Status:    &newStatus,
			StartedAt: &startTime,
		},
	)

	if err != nil {
		log.Printf("Error updating workflow instance status while stating it: %v", err)
		return
	}

	firstStep := workflowDefinition.GetFirstStep()
	if firstStep == nil {
		log.Printf("Workflow definition %s has no steps defined", workflowDefinition.Id.String())
		return
	}

	now := time.Now()
	inputData := firstStep.Input.GetValueMap(nil, job.WorkflowInstance.Input)
	instance := &models.StepInstance{
		WorkflowInstanceId: job.WorkflowInstance.Id,
		StepId:             firstStep.Id,
		Status:             models.StepInstanceStatusPending,
		Input:              inputData,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	finishedCh := make(chan struct{})
	err = e.stepQueue.Enqueue(queue.StepJob{
		WorkflowFinishedCh: finishedCh,
		WorkflowInstance:   job.WorkflowInstance,
		StepDefinition:     firstStep,
		StepInstance:       instance,
		WorkflowDefinition: workflowDefinition,
	})

	if err != nil {
		log.Printf("Error enqueueing first step of workflow instance: %v", err)
		return
	}

	<-finishedCh
	log.Printf("Workflow instance %s completed", job.WorkflowInstance.Id.String())

	newStatus = models.WorkflowInstanceStatusCompleted
	endTime := time.Now()
	_, err = e.workflowInstancesService.UpdateWorkflowInstance(
		job.WorkflowInstance.Id.String(),
		&dto.UpdateWorkflowInstanceRequest{
			Status:      &newStatus,
			CompletedAt: &endTime,
		},
	)

	if err != nil {
		log.Printf("Error updating workflow instance status while completing it: %v", err)
		return
	}
}
