package models

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
)

type StepType string // @name StepType

const (
	StepTypeTask     StepType = "task"
	StepTypeWorkflow StepType = "workflow"
	StepTypeWait     StepType = "wait"
	StepTypeDecision StepType = "decision"
	StepTypeFork     StepType = "fork"
	StepTypeJoin     StepType = "join"
)

type WorkflowStepDefinition struct {
	StepDefinitionID string                   `json:"stepDefinitionId"`
	Name             string                   `json:"name"`
	Description      string                   `json:"description"`
	Parameters       StepDefinitionParameters `json:"parameters"`
	Metadata         utils.UnknownJson        `json:"metadata,omitempty"`
	Type             StepType                 `json:"type"`
	TaskConfig       *TaskConfig              `json:"taskConfig,omitempty"`
	WorkflowConfig   *WorkflowConfig          `json:"workflowConfig,omitempty"`
	WaitConfig       *WaitConfig              `json:"waitConfig,omitempty"`
	DecisionConfig   *DecisionConfig          `json:"decisionConfig,omitempty"`
	ForkConfig       *ForkConfig              `json:"forkConfig,omitempty"`
	JoinConfig       *JoinConfig              `json:"joinConfig,omitempty"`
} // @name WorkflowStepDefinition

type TaskConfig struct {
	TaskDefinitionID string  `json:"taskDefinitionId"`
	NextStepID       *string `json:"nextStepId"`
} // @name TaskConfig

type WorkflowConfig struct {
	WorkflowDefinitionID string  `json:"workflowDefinitionId"`
	NextStepID           *string `json:"nextStepId"`
} // @name WorkflowConfig

type WaitConfig struct {
	DurationSeconds StepDefinitionParameter `json:"durationSeconds"`
	NextStepID      *string                 `json:"nextStepId"`
} // @name WaitConfig

type JoinConfig struct {
	IncomingStepIDs []string `json:"incomingStepIds"`
	NextStepID      *string  `json:"nextStepId"`
} // @name JoinConfig

type ForkBranch struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	NextStepID  string  `json:"nextStepId"`
} // @name ForkBranch

type ForkConfig struct {
	Branches []ForkBranch `json:"branches"`
} // @name ForkConfig

type DecisionCase struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Condition   string  `json:"condition"`
	NextStepID  string  `json:"nextStepId"`
} // @name DecisionCase

type DecisionConfig struct {
	Cases []DecisionCase `json:"cases"`
} // @name DecisionConfig

type WorkflowStepDefinitionList []WorkflowStepDefinition // @name WorkflowStepDefinitionList

func (list *WorkflowStepDefinitionList) Value() (driver.Value, error) {
	return json.Marshal(list)
}

func (list *WorkflowStepDefinitionList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, list)
}

func (def WorkflowStepDefinition) NewInstance(
	workflowDefinitionId string,
	workflowInstanceId string,
	input *map[string]interface{},
) *StepInstance {
	return &StepInstance{
		ID:                   uuid.New().String(),
		StepID:               def.Name,
		WorkflowDefinitionID: workflowDefinitionId,
		WorkflowInstanceID:   workflowInstanceId,
		Status:               StepStatusPending,
		Input:                utils.UnknownJsonFromMap(input),
		Progress:             0,
	}
}
