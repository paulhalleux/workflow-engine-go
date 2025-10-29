package main

import (
	"errors"
	"time"

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
	time.Sleep(3 * time.Second)

	if context.Input["message"] == nil {
		return nil, errors.New("message is empty")
	}

	return map[string]interface{}{
		"EchoedMessage": context.Input["message"],
	}, nil
}
