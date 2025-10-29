package dto

import "gorm.io/datatypes"

type CreateTypeSchemaRequest struct {
	Name    string          `json:"name" validate:"required"`
	Version string          `json:"version" validate:"required,wf_version"`
	Schema  *datatypes.JSON `json:"schema" validate:"required"`
} // @name CreateTypeSchemaRequest
