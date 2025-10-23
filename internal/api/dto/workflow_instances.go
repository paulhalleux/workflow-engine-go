package dto

import (
	"github.com/google/uuid"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"gorm.io/datatypes"
)

type CreateWorkflowInstanceRequest struct {
	WorkflowDefinitionId uuid.UUID       `json:"workflowDefinitionId" validate:"required,uuid4"`
	Input                datatypes.JSON  `json:"input" validate:"required"`
	Metadata             *datatypes.JSON `json:"metadata,omitempty"`
} // @name CreateWorkflowInstanceRequest

type UpdateWorkflowInstanceRequest struct {
	Status   *models.WorkflowInstanceStatus `json:"status,omitempty" validate:"omitempty,oneof=pending running completed failed canceled paused"`
	Output   *datatypes.JSON                `json:"output,omitempty"`
	Metadata *datatypes.JSON                `json:"metadata,omitempty"`
} // @name UpdateWorkflowInstanceRequest
