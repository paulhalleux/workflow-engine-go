package models

import (
	"time"

	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
)

type StepStatus string // @name StepStatus

const (
	StepStatusPending   StepStatus = "pending"
	StepStatusRunning   StepStatus = "running"
	StepStatusCompleted StepStatus = "completed"
	StepStatusFailed    StepStatus = "failed"
	StepStatusCancelled StepStatus = "cancelled"
)

type StepInstance struct {
	ID                   string             `gorm:"primaryKey" json:"id"`
	WorkflowDefinitionID string             `gorm:"type:uuid;not null" json:"workflowDefinitionId"`
	WorkflowInstanceID   string             `gorm:"type:uuid;not null" json:"workflowInstanceId"`
	StepID               string             `gorm:"type:string;not null" json:"stepId"`
	Status               StepStatus         `gorm:"type:varchar(50);not null" json:"status"`
	Input                *utils.UnknownJson `gorm:"type:jsonb" json:"input,omitempty"`
	Output               *utils.UnknownJson `gorm:"type:jsonb" json:"output,omitempty"`
	Progress             float64            `gorm:"type:float" json:"progress,omitempty"`
	ErrorMessage         *string            `gorm:"type:text" json:"errorMessage,omitempty"`
	CreatedAt            time.Time          `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt            time.Time          `gorm:"autoUpdateTime" json:"updatedAt"`
	StartedAt            *time.Time         `json:"startedAt,omitempty"`
	CompletedAt          *time.Time         `json:"completedAt,omitempty"`
} // @name StepInstance
