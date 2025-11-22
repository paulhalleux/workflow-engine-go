package models

import (
	"database/sql/driver"
	"encoding/json"
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
	StepDefinitionID string                    `json:"stepDefinitionId" validate:"required"`
	Name             string                    `json:"name" validate:"required"`
	Description      string                    `json:"description"`
	Parameters       *StepDefinitionParameters `json:"parameters,omitempty"`
	Metadata         *map[string]interface{}   `json:"metadata,omitempty"`
	Type             StepType                  `json:"type" validate:"required"`
	TaskConfig       *TaskConfig               `json:"taskConfig,omitempty"`
	WorkflowConfig   *WorkflowConfig           `json:"workflowConfig,omitempty"`
	WaitConfig       *WaitConfig               `json:"waitConfig,omitempty"`
	DecisionConfig   *DecisionConfig           `json:"decisionConfig,omitempty"`
	ForkConfig       *ForkConfig               `json:"forkConfig,omitempty"`
	JoinConfig       *JoinConfig               `json:"joinConfig,omitempty"`
} // @name WorkflowStepDefinition

type TaskConfig struct {
	TaskDefinitionID string  `json:"taskDefinitionId" validate:"required"`
	NextStepID       *string `json:"nextStepId,omitempty"`
} // @name TaskConfig

type WorkflowConfig struct {
	WorkflowDefinitionID string  `json:"workflowDefinitionId" validate:"required"`
	NextStepID           *string `json:"nextStepId,omitempty"`
} // @name WorkflowConfig

type WaitConfig struct {
	DurationSeconds StepDefinitionParameter `json:"durationSeconds" validate:"required"`
	NextStepID      *string                 `json:"nextStepId,omitempty"`
} // @name WaitConfig

type JoinConfig struct {
	IncomingStepIDs []string `json:"incomingStepIds" validate:"required"`
	NextStepID      *string  `json:"nextStepId,omitempty"`
} // @name JoinConfig

type ForkBranch struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	NextStepID  string  `json:"nextStepId" validate:"required"`
} // @name ForkBranch

type ForkConfig struct {
	JoinStepID string       `json:"joinStepId" validate:"required"`
	Branches   []ForkBranch `json:"branches" validate:"required"`
} // @name ForkConfig

type DecisionCase struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Condition   string  `json:"condition" validate:"required"`
	NextStepID  string  `json:"nextStepId" validate:"required"`
} // @name DecisionCase

type DecisionConfig struct {
	JoinStepID string         `json:"joinStepId" validate:"required"`
	Cases      []DecisionCase `json:"cases" validate:"required"`
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
