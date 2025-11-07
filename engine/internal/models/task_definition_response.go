package models

import "github.com/paulhalleux/workflow-engine-go/engine/internal/utils"

type TaskDefinitionResponse struct {
	Id               string            `json:"id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	InputParameters  utils.UnknownJson `json:"inputParameters"`
	OutputParameters utils.UnknownJson `json:"outputParameters"`
} // @name TaskDefinitionResponse
