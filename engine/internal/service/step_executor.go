package service

import (
	"context"
	"log"
	"time"

	"github.com/paulhalleux/workflow-engine-go/engine/internal"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/errors"
)

type StepExecution struct {
	StepInstanceID string
}

type StepExecutor struct {
	taskQueue chan *StepExecution
	sem       chan struct{}
}

func NewStepExecutor(
	config *internal.WorkflowEngineConfig,
) *StepExecutor {
	return &StepExecutor{
		taskQueue: make(chan *StepExecution, config.MaxWorkflowQueueSize),
		sem:       make(chan struct{}, config.MaxParallelWorkflows),
	}
}

func (we *StepExecutor) Enqueue(exec *StepExecution) error {
	select {
	case we.taskQueue <- exec:
		return nil
	default:
		return errors.ErrWorkflowQueueFull
	}
}

func (we *StepExecutor) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case exec := <-we.taskQueue:
			we.sem <- struct{}{}
			go func() {
				defer func() { <-we.sem }()
				we.startStep(exec)
			}()
		}
	}
}

func (we *StepExecutor) startStep(exec *StepExecution) {
	log.Println("Starting step execution:", exec.StepInstanceID)
	time.Sleep(2 * time.Second) // Simulate step execution
	log.Println("Completed step execution:", exec.StepInstanceID)
}
