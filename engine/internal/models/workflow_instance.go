package models

import (
	"time"

	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
)

type WorkflowStatus string // @name WorkflowStatus

const (
	WorkflowStatusPending   WorkflowStatus = "PENDING"
	WorkflowStatusRunning   WorkflowStatus = "RUNNING"
	WorkflowStatusCompleted WorkflowStatus = "COMPLETED"
	WorkflowStatusFailed    WorkflowStatus = "FAILED"
	WorkflowStatusCancelled WorkflowStatus = "CANCELLED"
)

type WorkflowInstance struct {
	ID                   string             `gorm:"primaryKey" json:"id"`
	WorkflowDefinitionID string             `gorm:"type:uuid;not null" json:"workflowDefinitionId"`
	Status               WorkflowStatus     `gorm:"type:varchar(50);not null" json:"status"`
	Input                *utils.UnknownJson `gorm:"type:jsonb" json:"input,omitempty"`
	Output               *utils.UnknownJson `gorm:"type:jsonb" json:"output,omitempty"`
	ErrorMessage         *string            `gorm:"type:text" json:"errorMessage,omitempty"`
	CreatedAt            time.Time          `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt            time.Time          `gorm:"autoUpdateTime" json:"updatedAt"`
	StartedAt            *time.Time         `json:"startedAt,omitempty"`
	CompletedAt          *time.Time         `json:"completedAt,omitempty"`
} // @name WorkflowInstance
