package worker

import (
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
) (map[string]interface{}, error) {
	value := job.Step.Wait.Duration.GetValue(job.Instance.Input)
	parsedValue, ok := value.(string)
	if !ok {
		return nil, nil
	}

	durationTime, err := time.ParseDuration(parsedValue)
	if err != nil {
		return nil, err
	}

	log.Printf("Waiting for %s", durationTime)
	time.Sleep(durationTime)
	return nil, nil
}
