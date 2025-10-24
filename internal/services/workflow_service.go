package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/persistence"
	"github.com/paulhalleux/workflow-engine-go/internal/queue"
	"github.com/paulhalleux/workflow-engine-go/internal/utils"
	"google.golang.org/protobuf/types/known/structpb"
)

type WorkflowService interface {
	StartWorkflow(
		ctx context.Context,
		definitionID uuid.UUID,
		input *structpb.Struct,
		metadata *structpb.Struct,
	) (uuid.UUID, error)
}

type workflowService struct {
	wfdRepo *persistence.WorkflowDefinitionsRepository
	wfiRepo *persistence.WorkflowInstancesRepository
	wfQueue queue.WorkflowQueue
}

func NewWorkflowService(
	wfdRepo *persistence.WorkflowDefinitionsRepository,
	wfiRepo *persistence.WorkflowInstancesRepository,
	wfQueue queue.WorkflowQueue,
) WorkflowService {
	return &workflowService{wfdRepo, wfiRepo, wfQueue}
}

func (s *workflowService) StartWorkflow(
	_ context.Context,
	definitionID uuid.UUID,
	input *structpb.Struct,
	metadata *structpb.Struct,
) (uuid.UUID, error) {
	def, err := s.wfdRepo.GetById(definitionID.String())
	if err != nil {
		return uuid.Nil, err
	}

	// Marshal input as JSON
	inputJSON, _ := utils.StructToJSON(input)
	metadataJSON, _ := utils.StructToJSON(metadata)

	// Create a new instance
	now := time.Now()
	instance := models.WorkflowInstance{
		Id:                   uuid.New(),
		WorkflowDefinitionId: def.Id,
		Status:               models.WorkflowInstanceStatusPending,
		CreatedAt:            now,
		UpdatedAt:            now,
		Input:                inputJSON,
		Metadata:             metadataJSON,
	}

	// Save the instance
	if err := s.wfiRepo.Create(&instance); err != nil {
		return uuid.Nil, err
	}

	// Enqueue the workflow job
	if err := s.wfQueue.Enqueue(queue.WorkflowJob{
		Instance: &instance,
	}); err != nil {
		return uuid.Nil, err
	}

	return instance.Id, nil
}
