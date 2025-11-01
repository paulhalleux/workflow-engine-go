package models

type UpdateWorkflowDefRequest struct {
	Name        *string                     `json:"name"`
	Description *string                     `json:"description"`
	Steps       *WorkflowStepDefinitionList `json:"steps"`
}

func (r *UpdateWorkflowDefRequest) ToWorkflowDefinition(def *WorkflowDefinition) *WorkflowDefinition {
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
