package connector

import (
	"context"
	"errors"
	"net"

	"github.com/paulhalleux/workflow-engine-go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcAgentConnector struct {
	address *string
	port    string

	connection *grpc.ClientConn
}

func NewGrpcAgentConnector(address *string, port string) (*GrpcAgentConnector, error) {
	connection, err := grpc.NewClient(joinHostPort(address, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Join(errors.New("failed to create gRPC connection"), err)
	}

	return &GrpcAgentConnector{
		address:    address,
		port:       port,
		connection: connection,
	}, nil
}

func (g *GrpcAgentConnector) Close() error {
	if g.connection != nil {
		return g.connection.Close()
	}
	return nil
}

func (g *GrpcAgentConnector) Ping() error {
	client := proto.NewAgentServiceClient(g.connection)
	_, err := client.Ping(context.Background(), &emptypb.Empty{})
	return err
}

func joinHostPort(host *string, port string) string {
	if host != nil {
		return net.JoinHostPort(*host, port)
	}
	return net.JoinHostPort("", port)
}
