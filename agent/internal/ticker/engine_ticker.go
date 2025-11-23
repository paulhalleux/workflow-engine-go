package ticker

import (
	"log"
	"time"

	"github.com/paulhalleux/workflow-engine-go/agent/internal/connector"
	"github.com/paulhalleux/workflow-engine-go/agent/internal/registry"
)

type EngineTicker struct {
	agentInfo              *connector.AgentInfo
	taskDefinitionRegistry *registry.TaskDefinitionRegistry
	ticker                 *time.Ticker
	engineConnector        connector.EngineConnector
}

func NewEngineTicker(
	delaySeconds int,
	agentInfo *connector.AgentInfo,
	taskDefinitionRegistry *registry.TaskDefinitionRegistry,
	engineConnector connector.EngineConnector,
) *EngineTicker {
	return &EngineTicker{
		ticker:                 time.NewTicker(time.Duration(delaySeconds) * time.Second),
		agentInfo:              agentInfo,
		taskDefinitionRegistry: taskDefinitionRegistry,
		engineConnector:        engineConnector,
	}
}

func (et *EngineTicker) Start() {
	for ; true; <-et.ticker.C {
		res, err := et.engineConnector.Ping(et.agentInfo.Name)
		if err != nil {
			log.Printf("[ticker] ping error: %v", err)
			continue
		}

		log.Printf("[ticker] ping result: %v", res)
		if !res {
			log.Printf("[ticker] engine does not know agent, registering agent: %s", et.agentInfo.Name)
			_, _ = et.engineConnector.RegisterAgent(
				et.agentInfo,
			)
		}
	}
}

func (et *EngineTicker) Stop() {
	et.ticker.Stop()
}
