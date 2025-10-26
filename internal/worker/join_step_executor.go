package worker

import (
	"log"
	"sync"

	"github.com/paulhalleux/workflow-engine-go/internal/queue"
	"gorm.io/gorm/utils"
)

type JoinStepExecutor struct{}

func NewJoinStepExecutor() *JoinStepExecutor {
	return &JoinStepExecutor{}
}

func (w *JoinStepExecutor) execute(
	job *queue.StepJob,
) (*StepResult, error) {
	wg := sync.WaitGroup{}
	wg.Add(len(job.StepDefinition.Join.IncomingStepIds))

	go func() {
		for {
			finishedStepId := <-job.StepFinishedCh
			if utils.Contains(job.StepDefinition.Join.IncomingStepIds, finishedStepId) {
				log.Printf("Join step '%s' received finished signal from step '%s'", job.StepDefinition.Id, finishedStepId)
				wg.Done()
			}
		}
	}()

	log.Printf("Waiting for tasks: %s", job.StepDefinition.Join.IncomingStepIds)
	wg.Wait()
	log.Printf("All incoming steps for join '%s' have finished", job.StepDefinition.Id)

	nextStepIds := make([]string, 0)
	if job.StepDefinition.Join.NextStepId != nil {
		nextStepIds = append(nextStepIds, *job.StepDefinition.Join.NextStepId)
	}

	return &StepResult{
		Output:      nil,
		NextStepIds: nextStepIds,
	}, nil
}
