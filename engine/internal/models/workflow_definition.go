package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulhalleux/workflow-engine-go/engine-new/internal/errors"
)

type WorkflowDefinition struct {
	ID               uuid.UUID                        `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id" validate:"required"`
	Name             string                           `gorm:"type:varchar(255);not null" json:"name" validate:"required"`
	Description      string                           `gorm:"type:text" json:"description"`
	Version          string                           `gorm:"type:varchar(50);not null;uniqueIndex:idx_name_version" json:"version" validate:"required"`
	IsEnabled        bool                             `gorm:"not null;default:false" json:"isEnabled" validate:"required"`
	InputParameters  *WorkflowParameterDefinitionList `gorm:"type:jsonb" json:"inputParameters,omitempty"`
	OutputParameters *interface{}                     `gorm:"type:jsonb" json:"outputParameters,omitempty"`
	Steps            *WorkflowStepDefinitionList      `gorm:"type:jsonb;not null" json:"steps,omitempty" validate:"required"`
	CreatedAt        time.Time                        `gorm:"autoCreateTime" json:"createdAt" validate:"required"`
	UpdatedAt        time.Time                        `gorm:"autoUpdateTime" json:"updatedAt" validate:"required"`
	Metadata         *map[string]interface{}          `gorm:"type:jsonb" json:"metadata,omitempty"`
} // @name WorkflowDefinition

func (def WorkflowDefinition) GetFirstStep() (*WorkflowStepDefinition, error) {
	if def.Steps == nil || len(*def.Steps) == 0 {
		return nil, errors.ErrWorkflowDefinitionNoSteps
	}
	return &(*def.Steps)[0], nil
}

func (def WorkflowDefinition) GetStepByID(id string) (*WorkflowStepDefinition, bool) {
	if def.Steps == nil {
		return nil, false
	}
	for _, step := range *def.Steps {
		if step.StepDefinitionID == id {
			return &step, true
		}
	}
	return nil, false
}
