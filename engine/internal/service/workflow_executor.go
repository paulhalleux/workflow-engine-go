package service

import (
	"context"
	"log"

	"github.com/paulhalleux/workflow-engine-go/engine/internal"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/errors"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/models"
)

type StepResult struct {
	StepID string
	Output *map[string]interface{}
}

type WorkflowExecutionResult struct {
	Success bool
	Message *string
}

type WorkflowExecution struct {
	WorkflowInstanceID string
}

type WorkflowExecutor struct {
	workflowInstancesService   *WorkflowInstanceService
	workflowDefinitionsService *WorkflowDefinitionsService
	stepExecutionService       *StepExecutionService

	workflowChan           *map[string]chan *WorkflowExecutionResult
	workflowStepOutputChan *map[string]chan *StepResult

	taskQueue chan *WorkflowExecution
	sem       chan struct{}
}

func NewWorkflowExecutor(
	config *internal.WorkflowEngineConfig,
	workflowDefinitionsService *WorkflowDefinitionsService,
	workflowInstancesService *WorkflowInstanceService,
	stepExecutionService *StepExecutionService,
	workflowChan *map[string]chan *WorkflowExecutionResult,
	workflowStepOutputChan *map[string]chan *StepResult,
) *WorkflowExecutor {
	return &WorkflowExecutor{
		workflowInstancesService:   workflowInstancesService,
		workflowDefinitionsService: workflowDefinitionsService,
		stepExecutionService:       stepExecutionService,

		workflowChan:           workflowChan,
		workflowStepOutputChan: workflowStepOutputChan,

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

func (we *WorkflowExecutor) startWorkflowInstance(exec *WorkflowExecution) {
	err := we.workflowInstancesService.Update(&models.WorkflowInstance{
		ID:     exec.WorkflowInstanceID,
		Status: models.WorkflowStatusRunning,
	})

	if err != nil {
		log.Println("Error updating workflow status to running:", err)
	}
}

func (we *WorkflowExecutor) completeWorkflowInstance(instance *models.WorkflowInstance, taskOutputs *map[string]*map[string]interface{}) {
	err := we.workflowInstancesService.Update(&models.WorkflowInstance{
		ID:     instance.ID,
		Status: models.WorkflowStatusCompleted,
		Output: nil,
	})

	if err != nil {
		log.Println("Error updating workflow status to completed:", err)
	}
}

func (we *WorkflowExecutor) startWorkflow(exec *WorkflowExecution) {
	instance, err := we.workflowInstancesService.GetByID(exec.WorkflowInstanceID)
	if err != nil {
		we.failWorkflowInstance(exec, "Failed to retrieve workflow instance")
		return
	}

	definition, err := we.workflowDefinitionsService.GetByID(instance.WorkflowDefinitionID)
	if err != nil {
		we.failWorkflowInstance(exec, "Failed to retrieve workflow definition")
		return
	}

	firstStep, err := definition.GetFirstStep()
	if err != nil {
		we.failWorkflowInstance(exec, err.Error())
		return
	}

	stepOutputMap := make(map[string]*map[string]interface{})
	stepOutputChan := make(chan *StepResult, 100)

	(*we.workflowStepOutputChan)[instance.ID] = stepOutputChan
	defer delete(*we.workflowStepOutputChan, instance.ID)

	(*we.workflowChan)[instance.ID] = make(chan *WorkflowExecutionResult)
	defer delete(*we.workflowChan, instance.ID)

	log.Printf("Starting workflow instance %s with first step %s", instance.ID, firstStep.Name)
	id, err := we.stepExecutionService.StartStep(firstStep, instance)
	if err != nil {
		we.failWorkflowInstance(exec, "failed to start first step")
		return
	}

	log.Printf("Started first step instance %s for workflow instance %s", *id, instance.ID)
	we.startWorkflowInstance(exec)

	for {
		select {
		case stepOutput := <-stepOutputChan:
			log.Printf("Workflow instance %s received output from step %s", instance.ID, stepOutput.StepID)
			stepOutputMap[stepOutput.StepID] = stepOutput.Output
		case result := <-(*we.workflowChan)[instance.ID]:
			log.Printf("Workflow instance %s received stop signal", instance.ID)
			if result.Success {
				we.completeWorkflowInstance(instance, &stepOutputMap)
			} else {
				we.failWorkflowInstance(exec, *result.Message)
			}
			return
		}
	}
}
