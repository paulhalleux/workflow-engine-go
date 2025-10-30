package internal

import "github.com/paulhalleux/workflow-engine-go/agent"

type inputParameters struct {
	Message string `json:"message" required:"true"`
}

type outputParameters struct {
	Message string `json:"message" required:"true"`
}

func NewEchoTaskDefinition() agent.TaskDefinition {
	return agent.TaskDefinition{
		ID:               "echo-task",
		Name:             "Echo Task",
		Description:      "A task that echoes a message",
		InputParameters:  agent.ReflectJsonSchema(inputParameters{}),
		OutputParameters: agent.ReflectJsonSchema(outputParameters{}),
		Handle:           func() {},
	}
}
