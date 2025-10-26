package main

import (
	"log"

	"github.com/paulhalleux/workflow-engine-go/internal/agent"
)

type EchoInput struct {
	Message string
}

type EchoTask struct{}

func NewEchoTask() *EchoTask {
	return &EchoTask{}
}

func (t *EchoTask) Execute(context agent.TaskExecutionContext) (map[string]interface{}, error) {
	log.Printf("EchoTask Execute")

	input, ok := context.Input.(EchoInput)
	if !ok {
		return nil, nil
	}

	return map[string]interface{}{
		"EchoedMessage": input.Message,
	}, nil
}
