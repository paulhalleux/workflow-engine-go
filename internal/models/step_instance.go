package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type StepInstanceStatus string // @name StepInstanceStatus
const (
	StepInstanceStatusPending   StepInstanceStatus = "pending"
	StepInstanceStatusRunning   StepInstanceStatus = "running"
	StepInstanceStatusCompleted StepInstanceStatus = "completed"
	StepInstanceStatusFailed    StepInstanceStatus = "failed"
	StepInstanceStatusCanceled  StepInstanceStatus = "canceled"
	StepInstanceStatusPaused    StepInstanceStatus = "paused"
)

type StepInstance struct {
	Id                 uuid.UUID          `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	WorkflowInstanceId uuid.UUID          `gorm:"type:uuid;not null" json:"workflowInstanceId"`
	StepId             string             `gorm:"type:uuid;not null" json:"stepId"`
	Status             StepInstanceStatus `gorm:"type:varchar(50);not null" json:"status"`
	CreatedAt          time.Time          `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time          `gorm:"autoUpdateTime" json:"updatedAt"`
	StartedAt          *time.Time         `json:"startedAt,omitempty"`
	CompletedAt        *time.Time         `json:"completedAt,omitempty"`
	Input              datatypes.JSONMap  `gorm:"type:text" json:"input"`
	Output             datatypes.JSONMap  `gorm:"type:text" json:"output"`
	Metadata           datatypes.JSONMap  `gorm:"type:jsonb" json:"metadata"`
} // @name StepInstance
