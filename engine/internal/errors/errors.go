package errors

type UnsupportedProtocolError struct{}

func (e *UnsupportedProtocolError) Error() string {
	return "unsupported agent protocol"
}

var ErrUnsupportedProtocol = &UnsupportedProtocolError{}

type WorkflowQueueFullError struct{}

func (e *WorkflowQueueFullError) Error() string {
	return "workflow queue is full"
}

var ErrWorkflowQueueFull = &WorkflowQueueFullError{}

type WorkflowDefinitionNoStepsError struct{}

func (e *WorkflowDefinitionNoStepsError) Error() string {
	return "workflow definition has no steps"
}

var ErrWorkflowDefinitionNoSteps = &WorkflowDefinitionNoStepsError{}
