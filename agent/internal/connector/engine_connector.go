package connector

import (
	"errors"

	"github.com/paulhalleux/workflow-engine-go/proto"
)

type AgentInfo struct {
	Name        string
	Version     string
	Address     string
	Port        string
	Definitions []*proto.TaskDefinition
}

type EngineConnector interface {
	Close() error
	Ping(name string) (bool, error)
	RegisterAgent(agent *AgentInfo) (bool, error)
}

func NewEngineConnector(protocol proto.AgentProtocol, address string) (EngineConnector, error) {
	switch protocol {
	case proto.AGENT_PROTOCOL_GRPC:
		return NewGrpcEngineConnector(address)
	default:
		return nil, errors.New("unsupported protocol")
	}
}
