package worker

import (
	"context"
	"log"

	"github.com/paulhalleux/workflow-engine-go/internal/grpcapi"
	"github.com/paulhalleux/workflow-engine-go/internal/proto"
	"github.com/paulhalleux/workflow-engine-go/internal/queue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

type AgentStepExecutor struct {
	executionChan chan grpcapi.TaskExecutionResult
	client        proto.AgentServiceClient
}

func NewAgentStepExecutor(executionChan chan grpcapi.TaskExecutionResult) *AgentStepExecutor {
	conn, err := grpc.NewClient("agent:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to agent: %v", err)
	}

	client := proto.NewAgentServiceClient(conn)

	return &AgentStepExecutor{executionChan, client}
}

func (w *AgentStepExecutor) execute(
	job *queue.StepJob,
) (*StepResult, error) {
	ctx := context.Background()
	input, err := structpb.NewStruct(job.StepInstance.Input)
	if err != nil {
		return nil, err
	}

	execution, err := w.client.StartTask(ctx, &proto.StartTaskRequest{
		TaskId:     job.StepDefinition.Task.TaskDefinitionId,
		Parameters: input,
	})

	if err != nil {
		return nil, err
	}

	log.Printf("Started execution for step definition: %v", job.StepDefinition.Id)
	var output map[string]interface{}
	for {
		executionResult := <-w.executionChan
		if executionResult.ExecutionId.String() == execution.ExecutionId {
			log.Printf("Step execution finished with output: %v", output)
			if executionResult.Type == grpcapi.TaskExecutionResultTypeCompletion {
				output = executionResult.Output
				break
			} else if executionResult.Type == grpcapi.TaskExecutionResultTypeFailure {
				return nil, executionResult.Error
			} else {
				break
			}
		}
	}

	var nextStepIds []string
	if job.StepDefinition.Task.NextStepId != nil {
		nextStepIds = append(nextStepIds, *job.StepDefinition.Task.NextStepId)
	}

	return &StepResult{
		Output:      output,
		NextStepIds: nextStepIds,
	}, nil
}
