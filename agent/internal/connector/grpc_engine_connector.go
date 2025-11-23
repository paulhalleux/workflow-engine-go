package connector

import (
	"context"
	"errors"

	"github.com/paulhalleux/workflow-engine-go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcEngineConnector struct {
	connection *grpc.ClientConn
}

func NewGrpcEngineConnector(address string) (*GrpcEngineConnector, error) {
	connection, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Join(errors.New("failed to create gRPC connection"), err)
	}

	return &GrpcEngineConnector{
		connection: connection,
	}, nil
}

func (g *GrpcEngineConnector) Close() error {
	if g.connection != nil {
		return g.connection.Close()
	}
	return nil
}

func (g *GrpcEngineConnector) Ping(name string) (bool, error) {
	client := proto.NewEngineServiceClient(g.connection)
	result, err := client.Ping(context.Background(), &proto.EnginePingRequest{
		Name: name,
	})

	if err != nil {
		return false, err
	}

	if result == nil {
		return false, errors.New("nil response from engine ping")
	}

	return result.KnowAgent, nil
}

func (g *GrpcEngineConnector) RegisterAgent(
	agent *AgentInfo,
) (bool, error) {
	client := proto.NewEngineServiceClient(g.connection)
	res, err := client.RegisterAgent(context.Background(), &proto.RegisterAgentRequest{
		Name:           agent.Name,
		Version:        agent.Version,
		Address:        &agent.Address,
		Port:           agent.Port,
		Protocol:       proto.AGENT_PROTOCOL_GRPC,
		SupportedTasks: agent.Definitions,
	})

	if err != nil {
		return false, err
	}

	if res == nil {
		return false, errors.New("nil response from engine register agent")
	}

	if res.Success == false {
		return false, errors.New(*res.Message)
	}

	return true, nil
}
