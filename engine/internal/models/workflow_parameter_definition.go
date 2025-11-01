package models

import (
	"database/sql/driver"
	"encoding/json"
)

type WorkflowParameterDefinition struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Required    bool                   `json:"required"`
	Default     interface{}            `json:"default"`
	Description *string                `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type WorkflowParameterDefinitionList []WorkflowParameterDefinition

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
