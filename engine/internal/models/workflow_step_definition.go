package models

import (
	"database/sql/driver"
	"encoding/json"
)

type StepType string

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
	Metadata         map[string]interface{}   `json:"metadata,omitempty"`
	Type             StepType                 `json:"type"`
	TaskConfig       *TaskConfig              `json:"taskConfig,omitempty"`
	WorkflowConfig   *WorkflowConfig          `json:"workflowConfig,omitempty"`
	WaitConfig       *WaitConfig              `json:"waitConfig,omitempty"`
	DecisionConfig   *DecisionConfig          `json:"decisionConfig,omitempty"`
	ForkConfig       *ForkConfig              `json:"forkConfig,omitempty"`
	JoinConfig       *JoinConfig              `json:"joinConfig,omitempty"`
}

type TaskConfig struct {
	TaskDefinitionID string  `json:"taskDefinitionId"`
	NextStepID       *string `json:"nextStepId"`
}

type WorkflowConfig struct {
	WorkflowDefinitionID string  `json:"workflowDefinitionId"`
	NextStepID           *string `json:"nextStepId"`
}

type WaitConfig struct {
	DurationSeconds StepDefinitionParameter `json:"durationSeconds"`
	NextStepID      *string                 `json:"nextStepId"`
}

type JoinConfig struct {
	IncomingStepIDs []string `json:"incomingStepIds"`
	NextStepID      *string  `json:"nextStepId"`
}

type ForkBranch struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	NextStepID  string  `json:"nextStepId"`
}

type ForkConfig struct {
	Branches []ForkBranch `json:"branches"`
}

type DecisionCase struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Condition   string  `json:"condition"`
	NextStepID  string  `json:"nextStepId"`
}

type DecisionConfig struct {
	Cases []DecisionCase `json:"cases"`
}

type WorkflowStepDefinitionList []WorkflowStepDefinition

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
