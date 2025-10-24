package worker

import (
	"time"

	"github.com/paulhalleux/workflow-engine-go/internal/queue"
)

type WaitStepExecutor struct{}

func NewWaitStepExecutor() *WaitStepExecutor {
	return &WaitStepExecutor{}
}

func (w *WaitStepExecutor) execute(
	job *queue.StepJob,
) (map[string]interface{}, error) {
	durationTime, err := time.ParseDuration(job.Step.Wait.Duration)
	if err != nil {
		return nil, err
	}

	time.Sleep(durationTime)
	return nil, nil
}
