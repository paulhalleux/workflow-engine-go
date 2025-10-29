package models

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/datatypes"
)

type WorkflowStepInputType string // @name WorkflowStepInputType
const (
	WorkflowStepInputTypeStepOutput    WorkflowStepInputType = "stepOutput"
	WorkflowStepInputTypeWorkflowInput WorkflowStepInputType = "workflowInput"
	WorkflowStepInputTypeConstant      WorkflowStepInputType = "constant"
	WorkflowStepInputTypeExpression    WorkflowStepInputType = "expression"
)

// WorkflowStepInput
// @Description Represents an input to a workflow step, which can be sourced from various types such as step outputs, workflow inputs, constants, expressions, or system values.
type WorkflowStepInput struct {
	Type          WorkflowStepInputType `json:"type,omitempty" validate:"required"`
	SourceStepId  string                `json:"sourceStepId,omitempty" validate:"required_if=type stepOutput"`
	SourceKey     string                `json:"sourceKey,omitempty" validate:"required_if=type stepOutput workflowInput"`
	ConstantValue interface{}           `json:"constantValue,omitempty" validate:"required_if=type constant"`
	Expression    string                `json:"expression,omitempty" validate:"required_if=type expression"`
} // @name WorkflowStepInput

type WorkflowStepInputMap map[string]WorkflowStepInput

func (w *WorkflowStepInputMap) Value() (driver.Value, error) {
	bytes, err := json.Marshal(w)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (w *WorkflowStepInputMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, w)
}

func (w *WorkflowStepInputMap) GetValueMap(stepOutputs map[string]datatypes.JSON, workflowInput datatypes.JSON) map[string]interface{} {
	result := make(map[string]interface{})
	for key, input := range *w {
		result[key] = input.GetValue(stepOutputs, workflowInput)
	}
	return result
}

func (w *WorkflowStepInput) GetValue(stepOutputs map[string]datatypes.JSON, workflowInput datatypes.JSON) interface{} {
	switch w.Type {
	case WorkflowStepInputTypeStepOutput:
		if stepOutput, exists := stepOutputs[w.SourceStepId]; exists {
			var outputMap map[string]interface{}
			if err := json.Unmarshal(stepOutput, &outputMap); err != nil {
				return nil
			}
			return outputMap[w.SourceKey]
		}
		return nil
	case WorkflowStepInputTypeWorkflowInput:
		var inputMap map[string]interface{}
		if err := json.Unmarshal(workflowInput, &inputMap); err != nil {
			return nil
		}
		return inputMap[w.SourceKey]
	case WorkflowStepInputTypeConstant:
		return w.ConstantValue
	case WorkflowStepInputTypeExpression:
		// Expression evaluation logic would go here
		return nil
	default:
		return nil
	}
}
