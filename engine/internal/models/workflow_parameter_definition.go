package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
)

type WorkflowParameterDefinition struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Required    bool              `json:"required"`
	Default     interface{}       `json:"default"`
	Description *string           `json:"description"`
	Metadata    utils.UnknownJson `json:"metadata"`
} // @name WorkflowParameterDefinition

type WorkflowParameterDefinitionList []WorkflowParameterDefinition // @name WorkflowParameterDefinitionList

func (list *WorkflowParameterDefinitionList) Value() (driver.Value, error) {
	bytes, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (list *WorkflowParameterDefinitionList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, list)
}
