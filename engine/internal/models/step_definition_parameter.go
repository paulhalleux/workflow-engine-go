package models

type StepParameterType string

const (
	StepParameterTypeConstant   StepParameterType = "constant"
	StepParameterTypeWorkflow   StepParameterType = "workflow"
	StepParameterTypeTaskOutput StepParameterType = "taskOutput"
)

type StepDefinitionParameter struct {
	Type  StepParameterType `json:"type"`
	Value interface{}       `json:"value"`
}

type StepDefinitionParameters map[string]StepDefinitionParameter
