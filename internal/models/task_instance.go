package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type TaskInstanceStatus string // @name TaskInstanceStatus
const (
	TaskInstanceStatusPending   TaskInstanceStatus = "pending"
	TaskInstanceStatusRunning   TaskInstanceStatus = "running"
	TaskInstanceStatusCompleted TaskInstanceStatus = "completed"
	TaskInstanceStatusFailed    TaskInstanceStatus = "failed"
	TaskInstanceStatusCanceled  TaskInstanceStatus = "canceled"
	TaskInstanceStatusPaused    TaskInstanceStatus = "paused"
)

type TaskInstance struct {
	Id                 uuid.UUID          `gorm:"type:uuid;primaryKey" json:"id"`
	WorkflowInstanceId uuid.UUID          `gorm:"type:uuid;not null" json:"workflowInstanceId"`
	StepId             uuid.UUID          `gorm:"type:uuid;not null" json:"stepId"`
	Status             TaskInstanceStatus `gorm:"type:varchar(50);not null" json:"status"`
	CreatedAt          int64              `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt          int64              `gorm:"autoUpdateTime" json:"updatedAt"`
	StartedAt          *int64             `json:"startedAt,omitempty"`
	CompletedAt        *int64             `json:"completedAt,omitempty"`
	Input              datatypes.JSON     `gorm:"type:text" json:"input"`
	Output             datatypes.JSON     `gorm:"type:text" json:"output"`
	Metadata           datatypes.JSON     `gorm:"type:jsonb" json:"metadata"`
} // @name TaskInstance
