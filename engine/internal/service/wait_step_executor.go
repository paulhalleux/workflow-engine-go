package service

import "time"

type WaitStepExecutor struct{}

func NewWaitStepExecutor() *WaitStepExecutor {
	return &WaitStepExecutor{}
}

func (e *WaitStepExecutor) Execute(exec *StepExecution) (*StepExecutionResult, error) {
	config := exec.StepDef.WaitConfig
	var nextStepIds []string
	if config.NextStepID != nil {
		nextStepIds = append(nextStepIds, *config.NextStepID)
	}

	time.Sleep(time.Second * 5)

	return &StepExecutionResult{
		NextStepIds: &nextStepIds,
	}, nil
}
