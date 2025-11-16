package service

import (
	"github.com/paulhalleux/workflow-engine-go/engine/internal"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/persistence"
	"github.com/paulhalleux/workflow-engine-go/proto"
)

type WorkflowExecutionService struct {
	workflowDefinitionsService *WorkflowDefinitionsService
	persistence                *persistence.Persistence
	workflowExecutor           *WorkflowExecutor
	websocketHub               *internal.WebsocketHub
}

func NewWorkflowExecutionService(
	workflowDefinitionsService *WorkflowDefinitionsService,
	persistence *persistence.Persistence,
	workflowExecutor *WorkflowExecutor,
	websocketHub *internal.WebsocketHub,
) *WorkflowExecutionService {
	return &WorkflowExecutionService{
		workflowDefinitionsService: workflowDefinitionsService,
		persistence:                persistence,
		workflowExecutor:           workflowExecutor,
		websocketHub:               websocketHub,
	}
}

func (ws *WorkflowExecutionService) StartWorkflow(definitionID string, parameters map[string]interface{}) (*string, error) {
	def, err := ws.workflowDefinitionsService.GetByID(definitionID)
	if err != nil {
		return nil, err
	}

	instance := def.NewInstance(&parameters)
	err = ws.persistence.WorkflowInstances.Create(instance)
	if err != nil {
		return nil, err
	}

	ws.websocketHub.BroadcastMessage(&proto.WebsocketMessage{
		Type: proto.WebsocketMessageType_WEBSOCKET_MESSAGE_TYPE_WORKFLOW_INSTANCE_EVENT,
		Payload: &proto.WebsocketMessage_WorkflowInstanceEvent{
			WorkflowInstanceEvent: &proto.WorkflowInstanceEvent{
				WorkflowInstanceId: instance.ID,
				EventType:          proto.WorkflowInstanceEventType_WORKFLOW_INSTANCE_EVENT_TYPE_CREATED,
				Details: &proto.WorkflowInstanceEvent_CreatedDetails{
					CreatedDetails: &proto.WorkflowInstanceCreatedDetails{},
				},
			},
		},
		Scope: &proto.WebsocketScope{
			Type: proto.WebsocketScopeType_WEBSOCKET_SCOPE_TYPE_WORKFLOW_INSTANCE,
			Id:   nil,
		},
	})

	exec := &WorkflowExecution{
		WorkflowInstanceID: instance.ID,
		StepCompletionChan: make(chan string, len(*def.Steps)),
		StepInstanceIDs:    make(map[string]string),
	}

	err = ws.workflowExecutor.Enqueue(exec)
	if err != nil {
		return nil, err
	}

	return &instance.ID, nil
}
