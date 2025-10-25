package worker

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/paulhalleux/workflow-engine-go/internal/dto"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/queue"
	"github.com/paulhalleux/workflow-engine-go/internal/services"
)

type WorkflowExecutor struct {
	stepInstanceService       services.StepInstanceService
	workflowInstancesService  services.WorkflowInstanceService
	workflowDefinitionService services.WorkflowDefinitionService
	workflowQueue             queue.WorkflowQueue
	stepQueue                 queue.StepQueue
	sem                       chan struct{}
}

func NewWorkflowExecutor(
	stepInstanceService services.StepInstanceService,
	workflowInstancesService services.WorkflowInstanceService,
	workflowDefinitionService services.WorkflowDefinitionService,
	workflowQueue queue.WorkflowQueue,
	stepQueue queue.StepQueue,
	maxConcurrentWorkflows int,
) *WorkflowExecutor {
	return &WorkflowExecutor{
		stepInstanceService:       stepInstanceService,
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

func (e *WorkflowExecutor) failWithError(job queue.WorkflowJob, err error) {
	log.Printf("Workflow instance %s failed with error: %v", job.WorkflowInstance.Id.String(), err)

	finishedAt := time.Now()
	failedStatus := models.WorkflowInstanceStatusFailed

	_, updateErr := e.workflowInstancesService.UpdateWorkflowInstance(
		job.WorkflowInstance.Id.String(),
		&dto.UpdateWorkflowInstanceRequest{
			Status:      &failedStatus,
			CompletedAt: &finishedAt,
		},
	)

	if updateErr != nil {
		log.Printf("Error updating workflow instance %s to failed status: %v", job.WorkflowInstance.Id.String(), updateErr)
	}
}

func (e *WorkflowExecutor) handle(job queue.WorkflowJob) {
	workflowDefinition, err := e.workflowDefinitionService.GetWorkflowDefinitionById(job.WorkflowInstance.WorkflowDefinitionId.String())
	if err != nil {
		log.Printf("Error loading workflow definition: %v", err)
		e.failWithError(job, err)
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
		e.failWithError(job, err)
		return
	}

	firstStep := workflowDefinition.GetFirstStep()
	if firstStep == nil {
		log.Printf("Workflow definition %s has no steps defined", workflowDefinition.Id.String())
		e.failWithError(job, errors.New("no steps defined in workflow"))
		return
	}

	now := time.Now()
	inputData := firstStep.Input.GetValueMap(nil, job.WorkflowInstance.Input)
	instance := models.StepInstance{
		WorkflowInstanceId: job.WorkflowInstance.Id,
		StepId:             firstStep.Id,
		Status:             models.StepInstanceStatusPending,
		Input:              inputData,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	_, err = e.stepInstanceService.CreateStepInstance(&instance)
	if err != nil {
		log.Printf("Error creating first step instance: %v", err)
		e.failWithError(job, err)
		return
	}

	finishedCh := make(chan queue.WorkflowExecutionResult)
	err = e.stepQueue.Enqueue(queue.StepJob{
		WorkflowFinishedCh: finishedCh,
		WorkflowInstance:   job.WorkflowInstance,
		StepDefinition:     firstStep,
		StepInstance:       &instance,
		WorkflowDefinition: workflowDefinition,
	})

	if err != nil {
		log.Printf("Error enqueueing first step of workflow instance: %v", err)
		e.failWithError(job, err)
		return
	}

	result := <-finishedCh
	if result.Error != nil {
		log.Printf("Workflow instance %s completed with error: %v", job.WorkflowInstance.Id.String(), result.Error)
	} else {
		log.Printf("Workflow instance %s completed with status %s", job.WorkflowInstance.Id.String(), result.Status)
	}

	endTime := time.Now()
	_, err = e.workflowInstancesService.UpdateWorkflowInstance(
		job.WorkflowInstance.Id.String(),
		&dto.UpdateWorkflowInstanceRequest{
			Status:      &result.Status,
			CompletedAt: &endTime,
		},
	)

	if err != nil {
		log.Printf("Error updating workflow instance status while completing it: %v", err)
		return
	}
}
