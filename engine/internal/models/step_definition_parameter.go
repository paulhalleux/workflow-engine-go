package models

type StepParameterType string // @name StepParameterType

const (
	StepParameterTypeConstant   StepParameterType = "constant"
	StepParameterTypeWorkflow   StepParameterType = "workflowInput"
	StepParameterTypeTaskOutput StepParameterType = "taskOutput"
)

type StepDefinitionParameter struct {
	Type  StepParameterType `json:"type" validate:"required"`
	Value interface{}       `json:"value" validate:"required"`
} // @name StepDefinitionParameter

type StepDefinitionParameters map[string]StepDefinitionParameter // @name StepDefinitionParameters
