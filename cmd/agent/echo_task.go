package main

import "github.com/paulhalleux/workflow-engine-go/internal/agent"

type EchoInput struct {
	Message string
}

type EchoOutput struct {
	EchoedMessage string
}

type EchoTask struct{}

func NewEchoTask() *EchoTask {
	return &EchoTask{}
}

func (t *EchoTask) Execute(context agent.TaskExecutionContext) (interface{}, error) {
	input, ok := context.Input.(EchoInput)
	if !ok {
		// TODO: fail
		return &EchoOutput{}, nil
	}

	return &EchoOutput{
		EchoedMessage: input.Message,
	}, nil
}
