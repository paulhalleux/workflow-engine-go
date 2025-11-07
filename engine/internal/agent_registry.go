package internal

import (
	"log"

	"github.com/paulhalleux/workflow-engine-go/engine/internal/connector"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/models"
	"github.com/paulhalleux/workflow-engine-go/proto"
)

type RegisteredAgent struct {
	Name           string
	Version        string
	Address        *string
	Port           string
	Protocol       proto.AgentProtocol
	SupportedTasks []*proto.TaskDefinition
}

type RegisteredAgentsList []RegisteredAgent

type AgentRegistry struct {
	agents           map[string]RegisteredAgent
	tasks            map[string]*proto.TaskDefinition
	agentsConnectors map[string]*connector.AgentConnector
}

func NewAgentRegistry() *AgentRegistry {
	return &AgentRegistry{
		agents:           make(map[string]RegisteredAgent),
		tasks:            make(map[string]*proto.TaskDefinition),
		agentsConnectors: make(map[string]*connector.AgentConnector),
	}
}

func (ar *AgentRegistry) RegisterAgent(name string, agent RegisteredAgent) error {
	agentConnector, err := connector.NewAgentConnector(agent.Protocol, agent.Address, agent.Port)
	if err == nil {
		err := agentConnector.Ping()
		if err != nil {
			return err
		}

		ar.agents[name] = agent
		ar.agentsConnectors[name] = &agentConnector

		for _, taskDef := range agent.SupportedTasks {
			ar.tasks[taskDef.Id] = taskDef
		}

		log.Printf("Registered agent %s at %v:%s using protocol %s", name, agent.Address, agent.Port, agent.Protocol.String())
	} else {
		log.Printf("Failed to register agent %s: %v", name, err)
	}
	return err
}

func (ar *AgentRegistry) GetAgent(name string) (RegisteredAgent, bool) {
	agent, exists := ar.agents[name]
	return agent, exists
}

func (ar *AgentRegistry) GetAgentConnector(name string) (*connector.AgentConnector, bool) {
	conn, exists := ar.agentsConnectors[name]
	return conn, exists
}

func (ar *AgentRegistry) UnregisterAgent(name string) {
	delete(ar.agents, name)
}

func (ar *AgentRegistry) ListAgents() RegisteredAgentsList {
	agents := make([]RegisteredAgent, 0, len(ar.agents))
	for _, agent := range ar.agents {
		agents = append(agents, agent)
	}
	return agents
}

func (ar *AgentRegistry) ListAgentNames() []string {
	names := make([]string, 0, len(ar.agents))
	for name := range ar.agents {
		names = append(names, name)
	}
	return names
}

func (ar *AgentRegistry) ListTaskDefinitions() []*models.TaskDefinitionResponse {
	taskDefs := make([]*proto.TaskDefinition, 0, len(ar.tasks))
	for _, taskDef := range ar.tasks {
		taskDefs = append(taskDefs, taskDef)
	}
	return taskDefResponseListFromProto(taskDefs)
}

func (ar *AgentRegistry) GetTaskDefinition(taskId string) (*models.TaskDefinitionResponse, bool) {
	taskDef, exists := ar.tasks[taskId]
	return taskDefResponseFromProto(taskDef), exists
}

func (ar *AgentRegistry) GetAgentByTaskName(taskName string) (*RegisteredAgent, bool) {
	for _, agent := range ar.agents {
		for _, taskDef := range agent.SupportedTasks {
			if taskDef.Id == taskName {
				return &agent, true
			}
		}
	}
	return nil, false
}

func (ag RegisteredAgent) ToResponse() models.AgentResponse {
	supportedTasks := make([]string, 0, len(ag.SupportedTasks))
	for _, taskDef := range ag.SupportedTasks {
		supportedTasks = append(supportedTasks, taskDef.Id)
	}

	return models.AgentResponse{
		Name:           ag.Name,
		Version:        ag.Version,
		Address:        ag.Address,
		Port:           ag.Port,
		Protocol:       ag.Protocol,
		SupportedTasks: supportedTasks,
	}
}

func (agl RegisteredAgentsList) ToResponseList() []models.AgentResponse {
	responses := make([]models.AgentResponse, 0, len(agl))
	for _, ag := range agl {
		responses = append(responses, ag.ToResponse())
	}
	return responses
}

func taskDefResponseFromProto(td *proto.TaskDefinition) *models.TaskDefinitionResponse {
	if td == nil {
		return nil
	}

	return &models.TaskDefinitionResponse{
		Id:               td.Id,
		Name:             td.Name,
		Description:      td.Description,
		InputParameters:  td.InputParameters.AsMap(),
		OutputParameters: td.OutputParameters.AsMap(),
	}
}

func taskDefResponseListFromProto(tds []*proto.TaskDefinition) []*models.TaskDefinitionResponse {
	responses := make([]*models.TaskDefinitionResponse, 0, len(tds))
	for _, td := range tds {
		responses = append(responses, taskDefResponseFromProto(td))
	}
	return responses
}
