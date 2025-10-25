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
	"gorm.io/datatypes"
)

type StepResult struct {
	Output      map[string]interface{}
	NextStepIds []string
}

type StepTypeExecutor interface {
	execute(job *queue.StepJob) (*StepResult, error)
}

var executors = map[models.WorkflowStepType]StepTypeExecutor{}

func RegisterStepExecutor(stepType models.WorkflowStepType, executor StepTypeExecutor) {
	executors[stepType] = executor
}

type StepExecutor struct {
	stepInstanceService services.StepInstanceService
	stepQueue           queue.StepQueue
	sem                 chan struct{}
}

func NewStepExecutor(
	stepInstanceService services.StepInstanceService,
	stepQueue queue.StepQueue,
	maxConcurrentWorkflows int,
) *StepExecutor {
	return &StepExecutor{
		stepInstanceService: stepInstanceService,
		stepQueue:           stepQueue,
		sem:                 make(chan struct{}, maxConcurrentWorkflows),
	}
}

func (e *StepExecutor) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case job := <-e.stepQueue.Dequeue():
				e.sem <- struct{}{}
				go func(j *queue.StepJob) {
					defer func() { <-e.sem }()
					e.handle(j)
				}(&job)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (e *StepExecutor) failWithError(job *queue.StepJob, err error) {
	newStatus := models.StepInstanceStatusFailed
	_, updateErr := e.stepInstanceService.UpdateStepInstance(job.StepInstance.Id, &dto.UpdateStepInstanceRequest{
		Status: &newStatus,
	})

	if updateErr != nil {
		log.Printf("Error updating step instance %s to failed: %v", job.StepInstance.Id, err)
	}

	job.WorkflowFinishedCh <- queue.WorkflowExecutionResult{
		Status: models.WorkflowInstanceStatusFailed,
		Error:  errors.Join(err, errors.New("step execution failed")),
	}
}

func (e *StepExecutor) markStepAsRunning(job *queue.StepJob) error {
	newStatus := models.StepInstanceStatusRunning
	now := time.Now()
	_, err := e.stepInstanceService.UpdateStepInstance(job.StepInstance.Id, &dto.UpdateStepInstanceRequest{
		Status:    &newStatus,
		StartedAt: &now,
	})
	return err
}

func (e *StepExecutor) markStepAsCompleted(job *queue.StepJob, output map[string]interface{}) error {
	newStatus := models.StepInstanceStatusCompleted
	now := time.Now()
	_, err := e.stepInstanceService.UpdateStepInstance(job.StepInstance.Id, &dto.UpdateStepInstanceRequest{
		Status:      &newStatus,
		Output:      (*datatypes.JSONMap)(&output),
		CompletedAt: &now,
	})
	return err
}

func (e *StepExecutor) handle(job *queue.StepJob) {
	executor, exists := executors[job.StepDefinition.Type]
	if !exists {
		log.Printf("No executor found for step type %s", job.StepDefinition.Type)
		job.WorkflowFinishedCh <- queue.WorkflowExecutionResult{
			Status: models.WorkflowInstanceStatusFailed,
			Error:  nil,
		}
		return
	}

	err := e.markStepAsRunning(job)
	if err != nil {
		log.Printf("Error marking step %s as running: %v", job.StepDefinition.Id, err)
		e.failWithError(job, err)
		return
	}

	result, err := executor.execute(job)
	if err != nil || result == nil {
		log.Printf("Error executing step %s: %v", job.StepDefinition.Id, err)
		e.failWithError(job, err)
		return
	}

	err = e.markStepAsCompleted(job, result.Output)
	if err != nil {
		log.Printf("Error marking step %s as completed: %v", job.StepDefinition.Id, err)
		e.failWithError(job, err)
		return
	}

	if len(result.NextStepIds) == 0 {
		log.Printf("Workflow %s finished", job.WorkflowInstance.Id)
		job.WorkflowFinishedCh <- queue.WorkflowExecutionResult{
			Status: models.WorkflowInstanceStatusCompleted,
			Error:  nil,
		}
		return
	}

	// Execute next steps
	for _, nextStepId := range result.NextStepIds {
		stepDef := job.WorkflowDefinition.GetStepById(nextStepId)
		if stepDef == nil {
			job.WorkflowFinishedCh <- queue.WorkflowExecutionResult{
				Status: models.WorkflowInstanceStatusFailed,
				Error:  nil,
			}
			log.Printf("Next step definition %s not found", nextStepId)
			break
		}

		now := time.Now()
		inputData := stepDef.Input.GetValueMap(nil, job.WorkflowInstance.Input)
		instance := &models.StepInstance{
			WorkflowInstanceId: job.WorkflowInstance.Id,
			StepId:             stepDef.Id,
			Status:             models.StepInstanceStatusPending,
			Input:              inputData,
			CreatedAt:          now,
			UpdatedAt:          now,
		}

		createdInstance, err := e.stepInstanceService.CreateStepInstance(instance)
		if err != nil {
			log.Printf("Error creating step instance for step %s: %v", stepDef.Id, err)
			job.WorkflowFinishedCh <- queue.WorkflowExecutionResult{
				Status: models.WorkflowInstanceStatusFailed,
				Error:  errors.Join(err, errors.New("failed to create step instance")),
			}
			return
		}

		nextJob := queue.StepJob{
			WorkflowFinishedCh: job.WorkflowFinishedCh,
			StepDefinition:     stepDef,
			StepInstance:       createdInstance,
			WorkflowDefinition: job.WorkflowDefinition,
			WorkflowInstance:   job.WorkflowInstance,
		}

		err = e.stepQueue.Enqueue(nextJob)
		if err != nil {
			log.Printf("Error enqueueing next step %s: %v", stepDef.Id, err)
			job.WorkflowFinishedCh <- queue.WorkflowExecutionResult{
				Status: models.WorkflowInstanceStatusFailed,
				Error:  errors.Join(err, errors.New("failed to enqueue next step")),
			}
			return
		}
	}
}
