package service

import (
	"context"
	"fmt"
	"log"

	"github.com/paulhalleux/workflow-engine-go/engine/internal"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/errors"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/models"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
)

type StepExecutionResult struct {
	NextStepIds *[]string
	Output      *map[string]interface{}
}

type StepTypeExecutor interface {
	Execute(exec *StepExecution) (*StepExecutionResult, error)
}

type StepExecution struct {
	StepInstanceID string
	StepDef        *models.WorkflowStepDefinition
	Input          *map[string]interface{}
}

type StepExecutor struct {
	stepInstanceService *StepInstanceService
	typeExecutors       map[models.StepType]StepTypeExecutor

	taskQueue chan *StepExecution
	sem       chan struct{}
}

func NewStepExecutor(
	config *internal.WorkflowEngineConfig,
	stepInstanceService *StepInstanceService,
) *StepExecutor {
	return &StepExecutor{
		stepInstanceService: stepInstanceService,
		typeExecutors:       make(map[models.StepType]StepTypeExecutor),

		taskQueue: make(chan *StepExecution, config.MaxWorkflowQueueSize),
		sem:       make(chan struct{}, config.MaxParallelWorkflows),
	}
}

func (we *StepExecutor) RegisterTypeExecutor(stepType models.StepType, executor StepTypeExecutor) {
	we.typeExecutors[stepType] = executor
}

func (we *StepExecutor) GetTypeExecutor(stepType models.StepType) (StepTypeExecutor, bool) {
	executor, exists := we.typeExecutors[stepType]
	return executor, exists
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

func (we *StepExecutor) failStepInstance(exec *StepExecution, messages string) {
	err := we.stepInstanceService.Update(&models.StepInstance{
		ID:           exec.StepInstanceID,
		Status:       models.StepStatusFailed,
		ErrorMessage: &messages,
	})

	if err != nil {
		log.Println("Error failing step instance:", err)
		return
	}
}

func (we *StepExecutor) startStepInstance(exec *StepExecution) {
	err := we.stepInstanceService.Update(&models.StepInstance{
		ID:     exec.StepInstanceID,
		Status: models.StepStatusRunning,
	})

	if err != nil {
		log.Println("Error starting step instance:", err)
		return
	}
}

func (we *StepExecutor) completeStepInstance(stepInstance *models.StepInstance, result *StepExecutionResult) {
	err := we.stepInstanceService.Update(&models.StepInstance{
		ID:     stepInstance.ID,
		Status: models.StepStatusCompleted,
		Output: utils.UnknownJsonFromMap(result.Output),
	})

	if err != nil {
		log.Println("Error completing step instance:", err)
		return
	}
}

func (we *StepExecutor) startStep(exec *StepExecution) {
	stepInstance, err := we.stepInstanceService.GetByID(exec.StepInstanceID)
	if err != nil {
		log.Println("Error fetching step instance:", err)
		we.failStepInstance(exec, "Failed to retrieve step instance")
		return
	}

	executor, exists := we.GetTypeExecutor(exec.StepDef.Type)
	if !exists {
		log.Println("No executor found for step type:", exec.StepDef.Type)
		we.failStepInstance(exec, fmt.Sprintf("No executor for step type: %s", exec.StepDef.Type))
		return
	}

	log.Println("Starting step instance:", stepInstance.ID)
	we.startStepInstance(exec)

	result, err := executor.Execute(exec)
	if err != nil {
		log.Println("Error executing step instance:", err)
		we.failStepInstance(exec, fmt.Sprintf("Step execution failed: %s", err.Error()))
		return
	}

	log.Println("Completed step instance:", stepInstance.ID, "with result:", result)
	we.completeStepInstance(stepInstance, result)

	if result.NextStepIds != nil {
		log.Println("Next steps to execute:", *result.NextStepIds)
	} else {
		log.Println("No next steps to execute.")
	}
}
