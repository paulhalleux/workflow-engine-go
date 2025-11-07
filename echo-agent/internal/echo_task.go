package internal

import (
	"time"

	"github.com/paulhalleux/workflow-engine-go/agent"
)

type inputParameters struct {
	Message string `json:"message" required:"true"`
}

type outputParameters struct {
	Message string `json:"message" required:"true"`
}

func NewEchoTaskDefinition() agent.TaskDefinition {
	return agent.TaskDefinition{
		ID:               "echo",
		Name:             "Echo Task",
		Description:      "A task that echoes a message",
		InputParameters:  agent.ReflectJsonSchema(inputParameters{}),
		OutputParameters: agent.ReflectJsonSchema(outputParameters{}),
		Handle: func(req *agent.TaskExecutionRequest) agent.TaskExecutionResult {
			time.Sleep(5 * time.Second) // Simulate some processing time
			return agent.TaskExecutionResult{
				Output: &map[string]interface{}{
					"message": req.Input["message"],
				},
				Error: nil,
			}
		},
	}
}
