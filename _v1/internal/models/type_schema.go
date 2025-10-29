package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// TypeSchema
// @Description A schema definition for a specific type, including its name, version, and schema (JSON schema).
type TypeSchema struct {
	Id      uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name    string         `json:"name"`
	Version string         `json:"version"`
	Schema  datatypes.JSON `json:"schema"`
} // @name TypeSchema
