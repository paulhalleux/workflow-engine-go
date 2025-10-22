package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type WorkflowDefinition struct {
	Id          uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string           `gorm:"not null" json:"name"`
	Description string           `json:"description"`
	Version     string           `gorm:"not null" json:"version"`
	Metadata    datatypes.JSON   `gorm:"type:jsonb" json:"metadata"`
	Steps       WorkflowStepList `gorm:"type:jsonb" json:"steps"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
