package worker

import (
	"context"
	"log"

	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/queue"
)

type StepExecutor struct {
	jobs <-chan queue.StepJob
	sem  chan struct{}
}

type StepTypeExecutor interface {
	execute(job *queue.StepJob) (map[string]interface{}, error)
}

var executors = map[models.WorkflowStepType]StepTypeExecutor{}

func RegisterStepExecutor(stepType models.WorkflowStepType, executor StepTypeExecutor) {
	executors[stepType] = executor
}

func NewStepExecutor(q queue.StepQueue, maxParallelJobs int) *StepExecutor {
	return &StepExecutor{
		jobs: q.Dequeue(),
		sem:  make(chan struct{}, maxParallelJobs),
	}
}

func (e *StepExecutor) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case job := <-e.jobs:
				e.sem <- struct{}{}
				go func(j queue.StepJob) {
					defer func() { <-e.sem }()
					e.handleJob(j)
				}(job)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (e *StepExecutor) handleJob(job queue.StepJob) {
	executor, exists := executors[job.Step.Type]
	if !exists {
		return
	}
	log.Printf("Executing step %s of type %s", job.Step.Id, job.Step.Type)
	output, err := executor.execute(&job)
	if err != nil {
		log.Printf("Error executing step %s: %v", job.Step.Id, err)
		return
	}
	log.Printf("Step %s executed successfully with output: %v", job.Step.Id, output)
}
