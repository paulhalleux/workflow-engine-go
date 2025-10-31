package internal

import (
	"log"

	"github.com/paulhalleux/workflow-engine-go/engine/internal/connector"
	"github.com/paulhalleux/workflow-engine-go/proto"
)

type RegisteredAgent struct {
	Name           string
	Address        *string
	Port           string
	Protocol       proto.AgentProtocol
	SupportedTasks []*proto.TaskDefinition
}

type AgentRegistry struct {
	agents           map[string]RegisteredAgent
	agentsConnectors map[string]*connector.AgentConnector
}

func NewAgentRegistry() *AgentRegistry {
	return &AgentRegistry{
		agents:           make(map[string]RegisteredAgent),
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

func (ar *AgentRegistry) UnregisterAgent(name string) {
	delete(ar.agents, name)
}

func (ar *AgentRegistry) ListAgents() []string {
	names := make([]string, 0, len(ar.agents))
	for name := range ar.agents {
		names = append(names, name)
	}
	return names
}
