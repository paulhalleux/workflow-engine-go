package service

import (
	"log"

	"github.com/paulhalleux/workflow-engine-go/engine/internal"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/errors"
	"github.com/paulhalleux/workflow-engine-go/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type AgentStepExecutor struct {
	agentRegistry *internal.AgentRegistry
	agentTaskChan map[string]chan *proto.NotifyTaskStatusRequest
}

func NewAgentStepExecutor(
	agentRegistry *internal.AgentRegistry,
	agentTaskChan map[string]chan *proto.NotifyTaskStatusRequest,
) *AgentStepExecutor {
	return &AgentStepExecutor{
		agentRegistry: agentRegistry,
		agentTaskChan: agentTaskChan,
	}
}

func (e *AgentStepExecutor) Execute(exec *StepExecution) (*StepExecutionResult, error) {
	config := exec.StepDef.TaskConfig
	agent, has := e.agentRegistry.GetAgentByTaskName(config.TaskDefinitionID)
	if !has {
		return nil, errors.ErrAgentNotFoundForTaskDefinition
	}

	connector, has := e.agentRegistry.GetAgentConnector(agent.Name)
	if !has {
		return nil, errors.ErrAgentConnectorNotFound
	}

	log.Printf("Executing agent step with config: %+v", config)
	input, err := structpb.NewStruct(*exec.Input)
	if err != nil {
		return nil, err
	}

	taskRes, err := (*connector).StartTask(&proto.StartTaskRequest{
		TaskName:        config.TaskDefinitionID,
		InputParameters: input,
	})

	if err != nil {
		return nil, err
	}

	// Wait for task completion
	e.agentTaskChan[taskRes.TaskId] = make(chan *proto.NotifyTaskStatusRequest)
	result := <-e.agentTaskChan[taskRes.TaskId]
	delete(e.agentTaskChan, taskRes.TaskId)

	log.Printf("Received task completion for TaskID: %s, Status: %v, OutputParameters: %v", result.TaskId, result.Status, result.OutputParameters)

	var nextStepIds []string
	if config.NextStepID != nil {
		nextStepIds = append(nextStepIds, *config.NextStepID)
	}

	output := result.OutputParameters.AsMap()
	return &StepExecutionResult{
		NextStepIds: &nextStepIds,
		Output:      &output,
	}, nil
}
