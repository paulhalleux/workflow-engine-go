package models

import "time"

type WorkflowDefinition struct {
	ID               string                           `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name             string                           `gorm:"type:varchar(255);not null" json:"name"`
	Description      string                           `gorm:"type:text" json:"description"`
	Version          string                           `gorm:"type:varchar(50);not null;uniqueIndex:idx_name_version" json:"version"`
	IsDraft          bool                             `gorm:"not null;default:true" json:"isDraft"`
	IsEnabled        bool                             `gorm:"not null;default:false" json:"isEnabled"`
	InputParameters  *WorkflowParameterDefinitionList `gorm:"type:jsonb" json:"inputParameters"`
	OutputParameters *WorkflowParameterDefinitionList `gorm:"type:jsonb" json:"outputParameters"`
	Steps            *WorkflowStepDefinitionList      `gorm:"type:jsonb;not null" json:"steps"`
	CreatedAt        time.Time                        `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time                        `gorm:"autoUpdateTime" json:"updatedAt"`
	Metadata         map[string]interface{}           `gorm:"type:jsonb" json:"metadata,omitempty"`
}
