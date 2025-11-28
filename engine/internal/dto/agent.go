package dto

import (
	"github.com/paulhalleux/workflow-engine-go/engine-new/internal/registry"
	"github.com/paulhalleux/workflow-engine-go/proto"
	"github.com/paulhalleux/workflow-engine-go/utils/array"
)

type AgentOverviewResponse struct {
	Name     string                      `json:"name"`
	Version  string                      `json:"version"`
	Protocol string                      `json:"protocol"`
	Address  *string                     `json:"address,omitempty"`
	Port     string                      `json:"port"`
	Tasks    []AgentTaskOverviewResponse `json:"tasks"`
} // @name AgentOverviewResponse

type AgentTaskOverviewResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
} // @name AgentTaskOverviewResponse

type AgentResponse struct {
	Name           string              `json:"name"`
	Version        string              `json:"version"`
	Protocol       string              `json:"protocol"`
	Address        *string             `json:"address,omitempty"`
	Port           string              `json:"port"`
	SupportedTasks []AgentTaskResponse `json:"supportedTasks"`
} // @name AgentResponse

type AgentTaskResponse struct {
	Id               string      `json:"id"`
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	InputParameters  interface{} `json:"inputParameters,omitempty"`
	OutputParameters interface{} `json:"outputParameters,omitempty"`
} // @name AgentTaskResponse

func NewAgentResponse(
	agent *registry.RegisteredAgent,
) AgentResponse {
	return AgentResponse{
		Name:     agent.Name,
		Version:  agent.Version,
		Protocol: agent.Protocol.String(),
		Address:  agent.Address,
		Port:     agent.Port,
		SupportedTasks: array.ToMapped(
			agent.SupportedTasks,
			func(taskDef *proto.TaskDefinition) AgentTaskResponse {
				return AgentTaskResponse{
					Id:               taskDef.Id,
					Name:             taskDef.Name,
					Description:      taskDef.Description,
					InputParameters:  taskDef.InputParameters,
					OutputParameters: taskDef.OutputParameters,
				}
			},
		),
	}
}

func NewAgentsOverviewResponse(
	agents []*registry.RegisteredAgent,
) []AgentOverviewResponse {
	return array.ToMapped(agents, NewAgentOverviewResponse)
}

func NewAgentOverviewResponse(
	agent *registry.RegisteredAgent,
) AgentOverviewResponse {
	return AgentOverviewResponse{
		Name:     agent.Name,
		Version:  agent.Version,
		Protocol: agent.Protocol.String(),
		Address:  agent.Address,
		Port:     agent.Port,
		Tasks: array.ToMapped(
			agent.SupportedTasks,
			func(taskDef *proto.TaskDefinition) AgentTaskOverviewResponse {
				return AgentTaskOverviewResponse{
					Id:          taskDef.Id,
					Name:        taskDef.Name,
					Description: taskDef.Description,
				}
			},
		),
	}
}
