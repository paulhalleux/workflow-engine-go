package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/errors"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
)

type WorkflowDefinition struct {
	ID               string                           `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name             string                           `gorm:"type:varchar(255);not null" json:"name"`
	Description      string                           `gorm:"type:text" json:"description"`
	Version          string                           `gorm:"type:varchar(50);not null;uniqueIndex:idx_name_version" json:"version"`
	IsDraft          bool                             `gorm:"not null;default:true" json:"isDraft"`
	IsEnabled        bool                             `gorm:"not null;default:false" json:"isEnabled"`
	InputParameters  *WorkflowParameterDefinitionList `gorm:"type:jsonb" json:"inputParameters"`
	OutputParameters *interface{}                     `gorm:"type:jsonb" json:"outputParameters"`
	Steps            *WorkflowStepDefinitionList      `gorm:"type:jsonb;not null" json:"steps"`
	CreatedAt        time.Time                        `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time                        `gorm:"autoUpdateTime" json:"updatedAt"`
	Metadata         *utils.UnknownJson               `gorm:"type:jsonb" json:"metadata,omitempty"`
} // @name WorkflowDefinition

func (def WorkflowDefinition) NewInstance(input *map[string]interface{}) *WorkflowInstance {
	return &WorkflowInstance{
		ID:                   uuid.New().String(),
		WorkflowDefinitionID: def.ID,
		Status:               WorkflowStatusPending,
		Input:                utils.UnknownJsonFromMap(input),
	}
}

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
