package worker

import (
	"errors"
	"log"
	"time"

	"github.com/paulhalleux/workflow-engine-go/internal/queue"
)

type WaitStepExecutor struct{}

func NewWaitStepExecutor() *WaitStepExecutor {
	return &WaitStepExecutor{}
}

func (w *WaitStepExecutor) execute(
	job *queue.StepJob,
) (*StepResult, error) {
	value := job.StepDefinition.Wait.Duration.GetValue(job.StepInstance.Input)
	parsedValue, ok := value.(string)
	if !ok {
		return nil, errors.New("invalid duration value")
	}

	durationTime, err := time.ParseDuration(parsedValue)
	if err != nil {
		return nil, err
	}

	log.Printf("Waiting for %s", durationTime)
	time.Sleep(durationTime)

	if job.StepDefinition.Wait.NextStepId != nil {
		return &StepResult{
			Output:      nil,
			NextStepIds: []string{*job.StepDefinition.Wait.NextStepId},
		}, nil
	}

	return &StepResult{
		Output:      nil,
		NextStepIds: []string{},
	}, nil
}
