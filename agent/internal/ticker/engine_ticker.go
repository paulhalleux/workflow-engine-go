package ticker

import (
	"time"

	"github.com/paulhalleux/workflow-engine-go/agent/internal"
	"github.com/paulhalleux/workflow-engine-go/agent/internal/connector"
)

type EngineTicker struct {
	agentConfig            *internal.WorkflowAgentConfig
	taskDefinitionRegistry *internal.TaskDefinitionRegistry
	engineConnector        connector.EngineConnector
	ticker                 *time.Ticker
}

func NewEngineTicker(
	delaySeconds int,
	agentConfig *internal.WorkflowAgentConfig,
	taskDefinitionRegistry *internal.TaskDefinitionRegistry,
	engineConnector connector.EngineConnector,
) *EngineTicker {
	return &EngineTicker{
		ticker:                 time.NewTicker(time.Duration(delaySeconds) * time.Second),
		agentConfig:            agentConfig,
		taskDefinitionRegistry: taskDefinitionRegistry,
		engineConnector:        engineConnector,
	}
}

func (et *EngineTicker) Start() {
	for range et.ticker.C {
		res, _ := et.engineConnector.Ping(et.agentConfig.Name)
		if !res {
			_, _ = et.engineConnector.RegisterAgent(
				et.agentConfig,
				et.taskDefinitionRegistry,
			)
		}
	}
}

func (et *EngineTicker) Stop() {
	et.ticker.Stop()
}
