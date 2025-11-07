package connector

import (
	"github.com/paulhalleux/workflow-engine-go/engine/internal/errors"
	"github.com/paulhalleux/workflow-engine-go/proto"
)

type AgentConnector interface {
	Close() error
	Ping() error
	StartTask(req *proto.StartTaskRequest) (*proto.TaskActionResponse, error)
}

func NewAgentConnector(protocol proto.AgentProtocol, address *string, port string) (AgentConnector, error) {
	switch protocol {
	case proto.AgentProtocol_GRPC:
		return NewGrpcAgentConnector(address, port)
	default:
		return nil, errors.ErrUnsupportedProtocol
	}
}
