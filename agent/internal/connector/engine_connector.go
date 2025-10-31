package connector

import (
	"github.com/paulhalleux/workflow-engine-go/agent/internal"
	"github.com/paulhalleux/workflow-engine-go/agent/internal/errors"
	"github.com/paulhalleux/workflow-engine-go/proto"
)

type EngineConnector interface {
	Close() error
	Ping(name string) (bool, error)
	RegisterAgent(config *internal.WorkflowAgentConfig, registry *internal.TaskDefinitionRegistry) (bool, error)
}

func NewEngineConnector(protocol proto.AgentProtocol, address string) (EngineConnector, error) {
	switch protocol {
	case proto.AgentProtocol_GRPC:
		return NewGrpcEngineConnector(address)
	default:
		return nil, errors.ErrUnsupportedProtocol
	}
}
