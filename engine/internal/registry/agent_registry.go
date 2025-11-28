package registry

import (
	"log"

	"github.com/paulhalleux/workflow-engine-go/engine-new/internal/connector"
	"github.com/paulhalleux/workflow-engine-go/proto"
)

type RegisteredAgent struct {
	Name           string                  `json:"name"`
	Version        string                  `json:"version"`
	Address        *string                 `json:"address,omitempty"`
	Port           string                  `json:"port"`
	Protocol       proto.AgentProtocol     `json:"protocol"`
	SupportedTasks []*proto.TaskDefinition `json:"supportedTasks"`
}

type RegisteredAgentsList []*RegisteredAgent

type AgentRegistry struct {
	agents           map[string]RegisteredAgent
	tasks            map[string]*proto.TaskDefinition
	agentByTask      map[string]string
	agentsConnectors map[string]*connector.AgentConnector
}

func NewAgentRegistry() *AgentRegistry {
	return &AgentRegistry{
		agents:           make(map[string]RegisteredAgent),
		tasks:            make(map[string]*proto.TaskDefinition),
		agentByTask:      make(map[string]string),
		agentsConnectors: make(map[string]*connector.AgentConnector),
	}
}

func (ar *AgentRegistry) RegisterAgent(name string, agent RegisteredAgent) error {
	log.Printf("[registry] registering agent %s at %v:%s using protocol %s", name, agent.Address, agent.Port, agent.Protocol.String())
	agentConnector, err := connector.NewAgentConnector(agent.Protocol, agent.Address, agent.Port)
	if err == nil {
		err := agentConnector.Ping()
		if err != nil {
			log.Printf("[registry] failed to ping agent %s at %v:%s: %v", name, agent.Address, agent.Port, err)
			return err
		}

		ar.agents[name] = agent
		ar.agentsConnectors[name] = &agentConnector

		for _, taskDef := range agent.SupportedTasks {
			ar.tasks[taskDef.Id] = taskDef
			ar.agentByTask[taskDef.Id] = name
		}

		log.Printf("[registry] registered agent %s at %v:%s using protocol %s", name, agent.Address, agent.Port, agent.Protocol.String())
	} else {
		log.Printf("[registry] railed to register agent %s: %v", name, err)
	}
	return err
}

func (ar *AgentRegistry) GetAgent(name string) (*RegisteredAgent, bool) {
	agent, exists := ar.agents[name]
	return &agent, exists
}

func (ar *AgentRegistry) GetAgentConnector(name string) (*connector.AgentConnector, bool) {
	conn, exists := ar.agentsConnectors[name]
	return conn, exists
}

func (ar *AgentRegistry) UnregisterAgent(name string) {
	delete(ar.agents, name)
	delete(ar.agentsConnectors, name)
	for taskId, agentName := range ar.agentByTask {
		if agentName == name {
			delete(ar.tasks, taskId)
			delete(ar.agentByTask, taskId)
		}
	}
}

func (ar *AgentRegistry) ListAgents() RegisteredAgentsList {
	agents := make([]*RegisteredAgent, 0, len(ar.agents))
	for _, agent := range ar.agents {
		agents = append(agents, &agent)
	}
	return agents
}
