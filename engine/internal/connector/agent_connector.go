package connector

import (
	"errors"

	"github.com/paulhalleux/workflow-engine-go/proto"
)

type AgentConnector interface {
	Close() error
	Ping() error
	StartTask(req *proto.StartTaskRequest) (*proto.TaskActionResponse, error)
}

func NewAgentConnector(protocol proto.AgentProtocol, address *string, port string) (AgentConnector, error) {
	switch protocol {
	case proto.AGENT_PROTOCOL_GRPC:
		return NewGrpcAgentConnector(address, port)
	default:
		return nil, errors.New("unsupported protocol")
	}
}
