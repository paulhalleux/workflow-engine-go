package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// WorkflowDefinition
// @Description A workflow definition defines the structure and behavior of a workflow, including its steps and metadata.
type WorkflowDefinition struct {
	Id          uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string            `gorm:"not null" json:"name"`
	Description string            `json:"description"`
	Version     string            `gorm:"not null" json:"version"`
	Steps       *WorkflowStepList `gorm:"type:jsonb" json:"steps"`
	IsEnabled   bool              `gorm:"default:false" json:"isEnabled"`
	Metadata    datatypes.JSON    `gorm:"type:jsonb" json:"metadata"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
} // @name WorkflowDefinition
