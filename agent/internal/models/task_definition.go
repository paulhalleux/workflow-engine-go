package models

import "github.com/swaggest/jsonschema-go"

type TaskHandler func(req *TaskExecutionRequest) TaskExecutionResult
type TaskDefinition struct {
	ID               string
	Name             string
	Description      string
	InputParameters  *jsonschema.Schema
	OutputParameters *jsonschema.Schema
	Handle           TaskHandler
}
