package service

import (
	"log"
	"sync"
)

type JoinStepExecutor struct{}

func NewJoinStepExecutor() *JoinStepExecutor {
	return &JoinStepExecutor{}
}

func (e *JoinStepExecutor) Execute(exec *StepExecution) (*StepExecutionResult, error) {
	config := exec.StepDef.JoinConfig

	count := len(config.IncomingStepIDs)
	wg := sync.WaitGroup{}
	wg.Add(len(config.IncomingStepIDs))

	go func() {
		for {
			stepID := <-exec.WorkflowExecution.StepCompletionChan
			if contains(config.IncomingStepIDs, stepID) {
				wg.Done()
				count--
			}

			if count == 0 {
				break
			}
		}
	}()

	log.Printf("Waiting for tasks: %+v", config.IncomingStepIDs)
	wg.Wait()
	log.Printf("All incoming steps for join '%s' have finished", exec.StepDef.StepDefinitionID)

	var nextStepIds []string
	if config.NextStepID != nil {
		nextStepIds = append(nextStepIds, *config.NextStepID)
	}

	return &StepExecutionResult{
		NextStepIds: &nextStepIds,
	}, nil
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
