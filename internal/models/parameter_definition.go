package models

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/datatypes"
)

// ParameterDefinition
// @Description A definition of a parameter, including its type, default value, and description.
type ParameterDefinition struct {
	DisplayName       string         `json:"displayName,omitempty"`
	Description       string         `json:"description,omitempty"`
	Type              string         `json:"type"`
	TypeSchemaName    string         `json:"typeSchemaName,omitempty"`
	TypeSchemaVersion string         `json:"typeSchemaVersion,omitempty"`
	DefaultValue      interface{}    `json:"defaultValue,omitempty"`
	Metadata          datatypes.JSON `json:"metadata,omitempty"`
} // @name ParameterDefinition

type ParameterDefinitionMap map[string]ParameterDefinition

func (s *ParameterDefinitionMap) Value() (driver.Value, error) {
	if s == nil || len(*s) == 0 {
		return "{}", nil
	}
	return json.Marshal(s)
}

func (s *ParameterDefinitionMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}
