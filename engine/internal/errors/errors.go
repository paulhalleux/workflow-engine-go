package errors

type SimpleError struct {
	Message string
}

func (e *SimpleError) Error() string {
	return e.Message
}

var ErrAgentNotFoundForTaskDefinition = &SimpleError{"no registered agent found for the given task definition"}
var ErrWorkflowDefinitionNoSteps = &SimpleError{"workflow definition has no steps"}
var ErrWorkflowQueueFull = &SimpleError{"workflow queue is full"}
var ErrStepDefinitionNotFound = &SimpleError{"step definition not found"}
var ErrUnsupportedProtocol = &SimpleError{"unsupported agent protocol"}
var ErrAgentConnectorNotFound = &SimpleError{"could not find a way to connect to the specified agent"}
