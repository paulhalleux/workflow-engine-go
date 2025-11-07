package service

import (
	"log"
)

type ForkStepExecutor struct{}

func NewForkStepExecutor() *ForkStepExecutor {
	return &ForkStepExecutor{}
}

func (e *ForkStepExecutor) Execute(exec *StepExecution) (*StepExecutionResult, error) {
	config := exec.StepDef.ForkConfig
	nextStepIds := make([]string, len(config.Branches))
	for i, branch := range config.Branches {
		nextStepIds[i] = branch.NextStepID
	}

	log.Printf("Executing fork step with branches: %+v", nextStepIds)

	return &StepExecutionResult{
		NextStepIds: &nextStepIds,
	}, nil
}
