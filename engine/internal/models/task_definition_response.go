package models

type TaskDefinitionResponse struct {
	Id               string                 `json:"id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	InputParameters  map[string]interface{} `json:"inputParameters"`
	OutputParameters map[string]interface{} `json:"outputParameters"`
}
