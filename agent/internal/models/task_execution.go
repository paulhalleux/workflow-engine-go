package models

type TaskExecutionResult struct {
	Output *map[string]interface{}
	Error  *error
}

type TaskExecutionRequest struct {
	Input map[string]interface{}
}
