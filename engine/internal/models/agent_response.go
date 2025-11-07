package models

import "github.com/paulhalleux/workflow-engine-go/proto"

type AgentResponse struct {
	Name           string              `json:"name"`
	Version        string              `json:"version"`
	Address        *string             `json:"address,omitempty"`
	Port           string              `json:"port"`
	Protocol       proto.AgentProtocol `json:"protocol"`
	SupportedTasks []string            `json:"supportedTasks"`
}
