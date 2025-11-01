package service

import (
	"context"
	"log"
	"time"

	"github.com/paulhalleux/workflow-engine-go/engine/internal"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/errors"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/models"
)

type WorkflowExecution struct {
	WorkflowInstanceID string
}

type WorkflowExecutor struct {
	workflowInstancesService *WorkflowInstanceService

	taskQueue chan *WorkflowExecution
	sem       chan struct{}
}

func NewWorkflowExecutor(
	config *internal.WorkflowEngineConfig,
	workflowInstancesService *WorkflowInstanceService,
) *WorkflowExecutor {
	return &WorkflowExecutor{
		workflowInstancesService: workflowInstancesService,

		taskQueue: make(chan *WorkflowExecution, config.MaxWorkflowQueueSize),
		sem:       make(chan struct{}, config.MaxParallelWorkflows),
	}
}

func (we *WorkflowExecutor) Enqueue(exec *WorkflowExecution) error {
	select {
	case we.taskQueue <- exec:
		return nil
	default:
		return errors.ErrWorkflowQueueFull
	}
}

func (we *WorkflowExecutor) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case exec := <-we.taskQueue:
			we.sem <- struct{}{}
			go func() {
				defer func() { <-we.sem }()
				we.startWorkflow(exec)
			}()
		}
	}
}

func (we *WorkflowExecutor) failWorkflowInstance(exec *WorkflowExecution, message string) {
	err := we.workflowInstancesService.Update(&models.WorkflowInstance{
		ID:           exec.WorkflowInstanceID,
		Status:       models.WorkflowStatusFailed,
		ErrorMessage: &message,
	})

	if err != nil {
		log.Println("Error updating workflow status to failed:", err)
	}
}

func (we *WorkflowExecutor) startWorkflow(exec *WorkflowExecution) {
	instance, err := we.workflowInstancesService.GetByID(exec.WorkflowInstanceID)
	if err != nil {
		we.failWorkflowInstance(exec, "Failed to retrieve workflow instance")
		return
	}

	log.Println("Starting workflow execution:", instance)
	time.Sleep(2 * time.Second) // Simulate workflow execution
	log.Println("Completed workflow execution:", exec.WorkflowInstanceID)
}
