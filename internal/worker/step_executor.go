package worker

import (
	"context"
	"log"
	"time"

	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/queue"
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
	stepQueue queue.StepQueue
	sem       chan struct{}
}

func NewStepExecutor(
	stepQueue queue.StepQueue,
	maxConcurrentWorkflows int,
) *StepExecutor {
	return &StepExecutor{
		stepQueue: stepQueue,
		sem:       make(chan struct{}, maxConcurrentWorkflows),
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

func (e *StepExecutor) handle(job *queue.StepJob) {
	executor, exists := executors[job.StepDefinition.Type]
	if !exists {
		log.Printf("No executor found for step type %s", job.StepDefinition.Type)
		job.WorkflowFinishedCh <- struct{}{}
		return
	}

	result, err := executor.execute(job)
	if err != nil || result == nil {
		log.Printf("Error executing step %s: %v", job.StepDefinition.Id, err)
		job.WorkflowFinishedCh <- struct{}{}
		return
	}

	if len(result.NextStepIds) == 0 {
		log.Printf("Workflow %s finished", job.WorkflowInstance.Id)
		job.WorkflowFinishedCh <- struct{}{}
		return
	}

	for _, nextStepId := range result.NextStepIds {
		stepDef := job.WorkflowDefinition.GetStepById(nextStepId)
		if stepDef == nil {
			job.WorkflowFinishedCh <- struct{}{}
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

		nextJob := queue.StepJob{
			WorkflowFinishedCh: job.WorkflowFinishedCh,
			StepDefinition:     stepDef,
			StepInstance:       instance,
			WorkflowDefinition: job.WorkflowDefinition,
			WorkflowInstance:   job.WorkflowInstance,
		}

		err := e.stepQueue.Enqueue(nextJob)
		if err != nil {
			log.Printf("Error enqueueing next step job: %v", err)
			continue
		}
	}
}
