package worker

import (
	"log"

	"github.com/paulhalleux/workflow-engine-go/internal/queue"
)

type ForkStepExecutor struct{}

func NewForkStepExecutor() *ForkStepExecutor {
	return &ForkStepExecutor{}
}

func (w *ForkStepExecutor) execute(
	job *queue.StepJob,
) (*StepResult, error) {
	nextStepIds := make([]string, 0)
	nextStepIds = append(nextStepIds, job.StepDefinition.Fork.JoinStepId)
	for _, branch := range job.StepDefinition.Fork.Branches {
		if branch.NextStepId == "" {
			continue
		}
		nextStepIds = append(nextStepIds, branch.NextStepId)
	}

	log.Printf("Forking for %s", nextStepIds)

	return &StepResult{
		Output:      nil,
		NextStepIds: nextStepIds,
	}, nil
}
