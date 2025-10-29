package agent

import (
	"context"
	"log"

	"github.com/paulhalleux/workflow-engine-go/internal/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type TaskExecutor struct {
	agent *Agent
	sem   chan struct{}
}

func NewTaskExecutor(agent *Agent, maxConcurrentTasks int) *TaskExecutor {
	return &TaskExecutor{
		agent: agent,
		sem:   make(chan struct{}, maxConcurrentTasks),
	}
}

func (e *TaskExecutor) Start(context context.Context) {
	go func() {
		for taskContext := range e.agent.Queue.Dequeue() {
			e.sem <- struct{}{}
			go func(tc *TaskExecutionContext) {
				defer func() { <-e.sem }()
				e.handleTask(context, tc)
			}(taskContext)
		}
	}()
}

func (e *TaskExecutor) handleTask(context context.Context, taskContext *TaskExecutionContext) {
	output, err := taskContext.Task.Execute(*taskContext)
	if err != nil {
		log.Printf("[TaskExecutor] Failed to execute task %s: %v", taskContext.TaskId, err)
		_, err = e.agent.WorkflowClient.NotifyTaskFailure(context, &proto.TaskFailureNotification{
			ExecutionId:  taskContext.ExecutionId.String(),
			ErrorMessage: err.Error(),
		})

		if err != nil {
			log.Printf("[TaskExecutor] Failed to notify task failure: %v", err)
		}

		return
	}

	outputStruct, err := structpb.NewStruct(output)
	_, err = e.agent.WorkflowClient.NotifyTaskCompletion(context, &proto.TaskCompletionNotification{
		ExecutionId: taskContext.ExecutionId.String(),
		Output:      outputStruct,
	})

	if err != nil {
		log.Printf("[TaskExecutor] Failed to notify task completion: %v", err)
		return
	}
}
