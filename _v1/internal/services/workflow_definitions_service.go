package services

import (
	"github.com/paulhalleux/workflow-engine-go/internal/dto"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/persistence"
)

type WorkflowDefinitionService interface {
	GetAllWorkflowDefinitions() ([]*models.WorkflowDefinition, error)
	GetWorkflowDefinitionById(id string) (*models.WorkflowDefinition, error)
	SearchWorkflowDefinitions(query *dto.SearchWorkflowDefinitionsRequest) ([]*models.WorkflowDefinition, error)
	CreateWorkflowDefinition(def *dto.CreateWorkflowDefinitionRequest) (*models.WorkflowDefinition, error)
	UpdateWorkflowDefinition(id string, def *dto.UpdateWorkflowDefinitionRequest) (*models.WorkflowDefinition, error)
	DeleteWorkflowDefinition(id string) error
	EnableWorkflowDefinition(id string) (*models.WorkflowDefinition, error)
	DisableWorkflowDefinition(id string) (*models.WorkflowDefinition, error)
}

type workflowDefinitionService struct {
	repo *persistence.WorkflowDefinitionsRepository
}

func NewWorkflowDefinitionService(
	repo *persistence.WorkflowDefinitionsRepository,
) WorkflowDefinitionService {
	return &workflowDefinitionService{repo}
}

func (s *workflowDefinitionService) GetAllWorkflowDefinitions() ([]*models.WorkflowDefinition, error) {
	return s.repo.GetAll()
}

func (s *workflowDefinitionService) GetWorkflowDefinitionById(id string) (*models.WorkflowDefinition, error) {
	return s.repo.GetById(id)
}

func (s *workflowDefinitionService) SearchWorkflowDefinitions(query *dto.SearchWorkflowDefinitionsRequest) ([]*models.WorkflowDefinition, error) {
	return s.repo.Search(query)
}

func (s *workflowDefinitionService) CreateWorkflowDefinition(def *dto.CreateWorkflowDefinitionRequest) (*models.WorkflowDefinition, error) {
	wf := &models.WorkflowDefinition{
		Name:             def.Name,
		Description:      def.Description,
		Version:          def.Version,
		Metadata:         def.Metadata,
		Steps:            &def.Steps,
		InputParameters:  &def.InputParameters,
		OutputParameters: &def.OutputParameters,
		IsEnabled:        def.IsEnabled,
	}
	return wf, s.repo.Create(wf)
}

func (s *workflowDefinitionService) UpdateWorkflowDefinition(id string, def *dto.UpdateWorkflowDefinitionRequest) (*models.WorkflowDefinition, error) {
	wf, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	updatePartialWorkflowDefinition(wf, def)
	return wf, s.repo.Update(wf)
}

func (s *workflowDefinitionService) DeleteWorkflowDefinition(id string) error {
	return s.repo.Delete(id)
}

func (s *workflowDefinitionService) EnableWorkflowDefinition(id string) (*models.WorkflowDefinition, error) {
	def, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	def.IsEnabled = true
	return def, s.repo.Update(def)
}

func (s *workflowDefinitionService) DisableWorkflowDefinition(id string) (*models.WorkflowDefinition, error) {
	def, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	def.IsEnabled = false
	return def, s.repo.Update(def)
}

func updatePartialWorkflowDefinition(wf *models.WorkflowDefinition, wfr *dto.UpdateWorkflowDefinitionRequest) {
	if wfr.Description != nil {
		wf.Description = *wfr.Description
	}
	if wfr.Version != nil {
		wf.Version = *wfr.Version
	}
	if wfr.Metadata != nil {
		wf.Metadata = *wfr.Metadata
	}
	if wfr.Steps != nil {
		wf.Steps = wfr.Steps
	}
	if wfr.InputParameters != nil {
		wf.InputParameters = wfr.InputParameters
	}
	if wfr.OutputParameters != nil {
		wf.OutputParameters = wfr.OutputParameters
	}
	if wfr.IsEnabled != nil {
		wf.IsEnabled = *wfr.IsEnabled
	}
}
