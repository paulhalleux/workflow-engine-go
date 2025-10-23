package dto

import (
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"gorm.io/datatypes"
)

type CreateWorkflowDefinitionRequest struct {
	Name        string                  `json:"name" validate:"required,min=3,max=100"`
	Description string                  `json:"description,omitempty" validate:"max=500"`
	Version     string                  `json:"version,omitempty" validate:"required,wf_version"`
	Metadata    datatypes.JSON          `json:"metadata,omitempty"`
	Steps       models.WorkflowStepList `json:"steps,omitempty" validate:"required,gt=0,dive,required"`
	IsEnabled   bool                    `json:"isEnabled,omitempty"`
} // @name CreateWorkflowDefinitionRequest

type UpdateWorkflowDefinitionRequest struct {
	Description *string                  `json:"description,omitempty" validate:"max=500"`
	Version     *string                  `json:"version,omitempty" validate:"wf_version"`
	Metadata    *datatypes.JSON          `json:"metadata,omitempty"`
	Steps       *models.WorkflowStepList `json:"steps,omitempty" validate:"gt=0,dive,required"`
	IsEnabled   *bool                    `json:"isEnabled,omitempty"`
} // @name UpdateWorkflowDefinitionRequest

type SearchWorkflowDefinitionsRequest struct {
	Name         *string `form:"name,omitempty"`
	Version      *string `form:"version,omitempty" validate:"omitnil,wf_version"`
	MajorVersion *int    `form:"majorVersion,omitempty" validate:"omitnil,gte=0"`
	MinorVersion *int    `form:"minorVersion,omitempty" validate:"omitnil,gte=0"`
	PatchVersion *int    `form:"patchVersion,omitempty" validate:"omitnil,gte=0"`
	IsRelease    *bool   `form:"isRelease,omitempty"`
	IsEnabled    *bool   `form:"isEnabled,omitempty"`
} // @name SearchWorkflowDefinitionsRequest
