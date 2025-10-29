package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type WorkflowInstanceStatus string // @name WorkflowInstanceStatus
const (
	WorkflowInstanceStatusPending   WorkflowInstanceStatus = "pending"
	WorkflowInstanceStatusRunning   WorkflowInstanceStatus = "running"
	WorkflowInstanceStatusCompleted WorkflowInstanceStatus = "completed"
	WorkflowInstanceStatusFailed    WorkflowInstanceStatus = "failed"
	WorkflowInstanceStatusCanceled  WorkflowInstanceStatus = "canceled"
	WorkflowInstanceStatusPaused    WorkflowInstanceStatus = "paused"
)

type WorkflowInstance struct {
	Id                   uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	WorkflowDefinitionId uuid.UUID              `gorm:"type:uuid;not null" json:"workflowDefinitionId"`
	Status               WorkflowInstanceStatus `gorm:"type:varchar(50);not null" json:"status"`
	StartedAt            *time.Time             `json:"startedAt,omitempty"`
	CompletedAt          *time.Time             `json:"completedAt,omitempty"`
	Input                datatypes.JSON         `gorm:"type:jsonb" json:"input"`
	Output               datatypes.JSON         `gorm:"type:jsonb" json:"output"`
	CreatedAt            time.Time              `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt            time.Time              `gorm:"autoUpdateTime" json:"updatedAt"`
	Metadata             datatypes.JSON         `gorm:"type:jsonb" json:"metadata"`
} // @name WorkflowInstance
