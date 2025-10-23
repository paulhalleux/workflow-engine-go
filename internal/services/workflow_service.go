package services

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/persistence"
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
}

func NewWorkflowService(
	wfdRepo *persistence.WorkflowDefinitionsRepository,
	wfiRepo *persistence.WorkflowInstancesRepository,
) WorkflowService {
	return &workflowService{wfdRepo, wfiRepo}
}

func (s *workflowService) StartWorkflow(
	_ context.Context,
	definitionID uuid.UUID,
	input *structpb.Struct,
	metadata *structpb.Struct,
) (uuid.UUID, error) {
	log.Println(
		"WorkflowService: Starting workflow",
		"definitionID=", definitionID,
		"input=", input,
		"metadata=", metadata,
	)

	def, err := s.wfdRepo.GetById(definitionID.String())
	if err != nil {
		return uuid.Nil, err
	}

	// Marshal input as JSON
	contextJSON, _ := utils.StructToJSON(input)
	metadataJSON, _ := utils.StructToJSON(metadata)

	now := time.Now()
	// Create a new instance
	instance := models.WorkflowInstance{
		Id:                   uuid.New(),
		WorkflowDefinitionId: def.Id,
		Status:               models.WorkflowInstanceStatusPending,
		CreatedAt:            now,
		UpdatedAt:            now,
		Input:                contextJSON,
		Metadata:             metadataJSON,
	}

	if err := s.wfiRepo.Create(&instance); err != nil {
		return uuid.Nil, err
	}

	return instance.Id, nil
}
