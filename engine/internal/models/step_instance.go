package models

import (
	"time"

	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
)

type StepStatus string

const (
	StepStatusPending   StepStatus = "PENDING"
	StepStatusRunning   StepStatus = "RUNNING"
	StepStatusCompleted StepStatus = "COMPLETED"
	StepStatusFailed    StepStatus = "FAILED"
	StepStatusCancelled StepStatus = "CANCELLED"
)

type StepInstance struct {
	ID                   string             `gorm:"primaryKey" json:"id"`
	WorkflowDefinitionID string             `gorm:"type:uuid;not null" json:"workflowDefinitionId"`
	StepID               string             `gorm:"type:string;not null" json:"stepId"`
	Status               StepStatus         `gorm:"type:varchar(50);not null" json:"status"`
	Input                *utils.UnknownJson `gorm:"type:jsonb" json:"input,omitempty"`
	Output               *utils.UnknownJson `gorm:"type:jsonb" json:"output,omitempty"`
	Progress             *float32           `gorm:"type:float" json:"progress,omitempty"`
	ErrorMessage         *string            `gorm:"type:text" json:"errorMessage,omitempty"`
	CreatedAt            time.Time          `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt            time.Time          `gorm:"autoUpdateTime" json:"updatedAt"`
	StartedAt            *time.Time         `json:"startedAt,omitempty"`
	CompletedAt          *time.Time         `json:"completedAt,omitempty"`
}
