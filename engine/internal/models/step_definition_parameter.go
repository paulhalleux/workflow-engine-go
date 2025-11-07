package models

type StepParameterType string // @name StepParameterType

const (
	StepParameterTypeConstant   StepParameterType = "constant"
	StepParameterTypeWorkflow   StepParameterType = "workflow"
	StepParameterTypeTaskOutput StepParameterType = "taskOutput"
)

type StepDefinitionParameter struct {
	Type  StepParameterType `json:"type"`
	Value interface{}       `json:"value"`
} // @name StepDefinitionParameter

type StepDefinitionParameters map[string]StepDefinitionParameter // @name StepDefinitionParameters

func (sdp StepDefinitionParameters) ToResolved(
	workflowInput *map[string]interface{},
	taskOutputs *map[string]map[string]interface{},
) *map[string]interface{} {
	result := make(map[string]interface{})

	for key, param := range sdp {
		result[key] = param.ToResolved(workflowInput, taskOutputs)
	}

	return &result
}

func (sdp StepDefinitionParameter) ToResolved(
	workflowInput *map[string]interface{},
	taskOutputs *map[string]map[string]interface{},
) interface{} {
	switch sdp.Type {
	case StepParameterTypeConstant:
		return sdp.Value
	case StepParameterTypeWorkflow:
		if workflowInput == nil {
			return nil
		}

		if key, ok := sdp.Value.(string); ok {
			return (*workflowInput)[key]
		}
	case StepParameterTypeTaskOutput:
		if ref, ok := sdp.Value.(map[string]interface{}); ok {
			taskID, ok1 := ref["taskId"].(string)
			outputKey, ok2 := ref["outputKey"].(string)
			if ok1 && ok2 {
				if taskOutput, exists := (*taskOutputs)[taskID]; exists {
					return taskOutput[outputKey]
				}
			}
		}
	}
	return nil
}
