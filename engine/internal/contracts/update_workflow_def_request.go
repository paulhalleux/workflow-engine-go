package contracts

import "github.com/paulhalleux/workflow-engine-go/engine/internal/models"

type UpdateWorkflowDefRequest struct {
	Name        *string                            `json:"name"`
	Description *string                            `json:"description"`
	Steps       *models.WorkflowStepDefinitionList `json:"steps"`
}

func (r *UpdateWorkflowDefRequest) ToWorkflowDefinition(def *models.WorkflowDefinition) *models.WorkflowDefinition {
	if r.Name != nil {
		def.Name = *r.Name
	}
	if r.Description != nil {
		def.Description = *r.Description
	}
	if r.Steps != nil {
		def.Steps = r.Steps
	}
	return def
}
