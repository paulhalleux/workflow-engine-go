package worker

import (
	"context"
	"log"
	"time"

	"github.com/paulhalleux/workflow-engine-go/internal/persistence"
	"github.com/paulhalleux/workflow-engine-go/internal/queue"
)

type WorkflowExecutor struct {
	wfdRepo *persistence.WorkflowDefinitionsRepository
	wfiRepo *persistence.WorkflowInstancesRepository
	jobs    <-chan queue.WorkflowJob
	sem     chan struct{}
}

func NewWorkflowExecutor(
	wfdRepo *persistence.WorkflowDefinitionsRepository,
	wfiRepo *persistence.WorkflowInstancesRepository,
	q queue.WorkflowQueue,
	maxParallelJobs int,
) *WorkflowExecutor {
	return &WorkflowExecutor{
		wfdRepo: wfdRepo,
		wfiRepo: wfiRepo,
		jobs:    q.Dequeue(),
		sem:     make(chan struct{}, maxParallelJobs),
	}
}

func (e *WorkflowExecutor) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case job := <-e.jobs:
				e.sem <- struct{}{}
				go func(j queue.WorkflowJob) {
					defer func() { <-e.sem }()
					e.handleJob(j)
				}(job)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (e *WorkflowExecutor) handleJob(job queue.WorkflowJob) {
	log.Printf("Executing workflow instance %s", job.InstanceId)

	// 1. Load instance + definition
	instance, err := e.wfiRepo.GetById(job.InstanceId.String())
	if err != nil {
		log.Printf("Error loading instance: %v", err)
		return
	}

	def, err := e.wfdRepo.GetById(instance.WorkflowDefinitionId.String())
	if err != nil {
		log.Printf("Error loading definition: %v", err)
		return
	}

	// 2. Parse steps (you already store as JSON)
	// 3. Execute first step(s)
	log.Printf("Workflow %s started with %d steps", def.Name, len(*def.Steps))

	// wait  for 5 seconds to simulate work
	time.Sleep(5 * time.Second)
}
