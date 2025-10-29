package services

import (
	"github.com/paulhalleux/workflow-engine-go/internal/dto"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/persistence"
	"gorm.io/datatypes"
)

type WorkflowInstanceService interface {
	GetAllWorkflowInstances() ([]*models.WorkflowInstance, error)
	GetWorkflowInstanceById(id string) (*models.WorkflowInstance, error)
	CreateWorkflowInstance(req *dto.CreateWorkflowInstanceRequest) (*models.WorkflowInstance, error)
	UpdateWorkflowInstance(id string, req *dto.UpdateWorkflowInstanceRequest) (*models.WorkflowInstance, error)
	DeleteWorkflowInstance(id string) error
}

type workflowInstanceService struct {
	repo *persistence.WorkflowInstancesRepository
}

func NewWorkflowInstanceService(repo *persistence.WorkflowInstancesRepository) WorkflowInstanceService {
	return &workflowInstanceService{repo}
}

func (s *workflowInstanceService) GetAllWorkflowInstances() ([]*models.WorkflowInstance, error) {
	wfs, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	res := make([]*models.WorkflowInstance, 0, len(wfs))
	for i := range wfs {
		// create a copy to avoid returning the address of a loop variable or slice element
		copyItem := wfs[i]
		res = append(res, &copyItem)
	}
	return res, nil
}

func (s *workflowInstanceService) GetWorkflowInstanceById(id string) (*models.WorkflowInstance, error) {
	return s.repo.GetById(id)
}

func (s *workflowInstanceService) CreateWorkflowInstance(req *dto.CreateWorkflowInstanceRequest) (*models.WorkflowInstance, error) {
	wf := &models.WorkflowInstance{
		WorkflowDefinitionId: req.WorkflowDefinitionId,
		Input:                req.Input,
		Metadata:             datatypes.JSON([]byte("{}")),
		Status:               models.WorkflowInstanceStatusPending,
	}
	if req.Metadata != nil {
		wf.Metadata = *req.Metadata
	}
	return wf, s.repo.Create(wf)
}

func (s *workflowInstanceService) UpdateWorkflowInstance(id string, req *dto.UpdateWorkflowInstanceRequest) (*models.WorkflowInstance, error) {
	wf, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	updatePartialWorkflowInstance(wf, req)
	return wf, s.repo.Update(wf)
}

func (s *workflowInstanceService) DeleteWorkflowInstance(id string) error {
	return s.repo.Delete(id)
}

func updatePartialWorkflowInstance(wf *models.WorkflowInstance, req *dto.UpdateWorkflowInstanceRequest) {
	if req.Output != nil {
		wf.Output = *req.Output
	}
	if req.Metadata != nil {
		wf.Metadata = *req.Metadata
	}
	if req.Status != nil {
		wf.Status = *req.Status
	}
}
