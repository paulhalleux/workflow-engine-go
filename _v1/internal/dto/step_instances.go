package dto

import (
	"time"

	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"gorm.io/datatypes"
)

type UpdateStepInstanceRequest struct {
	Status      *models.StepInstanceStatus `json:"status,omitempty" validate:"omitempty,oneof=pending running completed failed canceled paused"`
	Output      *datatypes.JSONMap         `json:"output,omitempty"`
	Metadata    *datatypes.JSONMap         `json:"metadata,omitempty"`
	StartedAt   *time.Time                 `json:"startedAt,omitempty"`
	CompletedAt *time.Time                 `json:"completedAt,omitempty"`
} // @name UpdateStepInstanceRequest
