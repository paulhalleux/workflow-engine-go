package models

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/datatypes"
)

type WorkflowStepType string // @name WorkflowStepType
const (
	WorkflowStepTypeTask        WorkflowStepType = "task"
	WorkflowStepTypeFork        WorkflowStepType = "fork"
	WorkflowStepTypeJoin        WorkflowStepType = "join"
	WorkflowStepTypeDecision    WorkflowStepType = "decision"
	WorkflowStepTypeSubWorkflow WorkflowStepType = "subworkflow"
	WorkflowStepTypeWait        WorkflowStepType = "wait"
	WorkflowStepTypeScript      WorkflowStepType = "script"
	WorkflowStepTypeHttp        WorkflowStepType = "http"
	WorkflowStepTypeTerminate   WorkflowStepType = "terminate"
)

// BaseWorkflowStep
// @Description The base structure for a workflow step, containing common fields.
type BaseWorkflowStep struct {
	Id          string               `json:"id,omitempty" validate:"required"`
	DisplayName string               `json:"displayName,omitempty"`
	Description string               `json:"description,omitempty"`
	Input       WorkflowStepInputMap `json:"input,omitempty"`
	Type        WorkflowStepType     `json:"type,omitempty" validate:"required"`
} // @name BaseWorkflowStep

type StepParameterBindingType string // @name StepParameterBindingType
const (
	StepParameterBindingTypeInput StepParameterBindingType = "input"
	StepParameterBindingTypeValue StepParameterBindingType = "value"
)

type StepParameterBinding struct {
	Type      StepParameterBindingType `json:"type,omitempty" validate:"required,oneof=input value"`
	SourceKey string                   `json:"sourceKey,omitempty" validate:"required_if=Type input"`
	Value     interface{}              `json:"value,omitempty" validate:"required_if=Type value"`
} // @name StepParameterBinding

// TaskWorkflowStep
// @Description A simple task workflow step that performs a specific action.
type TaskWorkflowStep struct {
	TaskDefinitionId string `json:"taskDefinitionId,omitempty" validate:"required"`
	NextStepId       string `json:"nextStepId,omitempty"`
} // @name TaskWorkflowStep

// ForkBranch
// @Description A single branch within a fork workflow step, defining the next step to take.
type ForkBranch struct {
	Id         string `json:"id,omitempty"`
	NextStepId string `json:"nextStepId,omitempty"`
} // @name ForkBranch

// ForkWorkflowStep
// @Description A workflow step that forks into multiple branches.
type ForkWorkflowStep struct {
	JoinStepId string       `json:"joinStepId,omitempty"`
	Branches   []ForkBranch `json:"branches,omitempty"`
} // @name ForkWorkflowStep

// DecisionCase
// @Description A single case within a decision workflow step, defining a condition and the next step to take if the condition is met.
type DecisionCase struct {
	Condition  string `json:"condition,omitempty"`
	NextStepId string `json:"nextStepId,omitempty"`
} // @name DecisionCase

// DecisionWorkflowStep
// @Description A workflow step that makes a decision based on conditions to determine the next step.
type DecisionWorkflowStep struct {
	Cases []DecisionCase `json:"cases,omitempty"`
} // @name DecisionWorkflowStep

// JoinWorkflowStep
// @Description A workflow step that joins multiple incoming steps into one.
type JoinWorkflowStep struct {
	OriginId        string   `json:"originId,omitempty"`
	IncomingStepIds []string `json:"IncomingStepIds,omitempty"`
	NextStepId      *string  `json:"nextStepId,omitempty"`
} // @name JoinWorkflowStep

// SubWorkflowStep
// @Description A workflow step that invokes a sub-workflow.
type SubWorkflowStep struct {
	WorkflowDefinitionId string `json:"workflowId,omitempty"`
	NextStepId           string `json:"nextStepId,omitempty"`
} // @name SubWorkflowStep

// WaitWorkflowStep
// @Description A workflow step that waits for a specified duration before proceeding.
type WaitWorkflowStep struct {
	Duration   StepParameterBinding `json:"duration,omitempty"`
	NextStepId *string              `json:"nextStepId,omitempty"`
} // @name WaitWorkflowStep

// TerminateWorkflowStep
// @Description A workflow step that terminates the workflow execution.
type TerminateWorkflowStep struct {
	Reason string `json:"reason,omitempty"`
} // @name TerminateWorkflowStep

// ScriptWorkflowStep
// @Description A workflow step that executes a script.
type ScriptWorkflowStep struct {
	Script     string `json:"script,omitempty"`
	NextStepId string `json:"nextStepId,omitempty"`
} // @name ScriptWorkflowStep

// HttpWorkflowStep
// @Description A workflow step that performs an HTTP request.
type HttpWorkflowStep struct {
	Method     string            `json:"method,omitempty"`
	Url        string            `json:"url,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
	NextStepId string            `json:"nextStepId,omitempty"`
} // @name HttpWorkflowStep

// WorkflowStep
// @Description A step within a workflow, defining its type and specific configurations.
type WorkflowStep struct {
	BaseWorkflowStep
	Task        *TaskWorkflowStep      `json:"task,omitempty" validate:"required_if=Type task"`
	Fork        *ForkWorkflowStep      `json:"fork,omitempty" validate:"required_if=Type fork"`
	Decision    *DecisionWorkflowStep  `json:"decision,omitempty" validate:"required_if=Type decision"`
	Join        *JoinWorkflowStep      `json:"join,omitempty" validate:"required_if=Type join"`
	SubWorkflow *SubWorkflowStep       `json:"subWorkflow,omitempty" validate:"required_if=Type subworkflow"`
	Wait        *WaitWorkflowStep      `json:"wait,omitempty" validate:"required_if=Type wait"`
	Terminate   *TerminateWorkflowStep `json:"terminate,omitempty" validate:"required_if=Type terminate"`
	Script      *ScriptWorkflowStep    `json:"script,omitempty" validate:"required_if=Type script"`
	Http        *HttpWorkflowStep      `json:"http,omitempty" validate:"required_if=Type http"`
} // @name WorkflowStep

type WorkflowStepList []WorkflowStep

func (s *WorkflowStepList) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

func (s *WorkflowStepList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s *StepParameterBinding) GetValue(input datatypes.JSONMap) interface{} {
	if s.Type == StepParameterBindingTypeInput {
		return input[s.SourceKey]
	} else if s.Type == StepParameterBindingTypeValue {
		return s.Value
	}
	return nil
}
