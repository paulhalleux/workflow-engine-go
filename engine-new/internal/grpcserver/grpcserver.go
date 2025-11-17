package grpcserver

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	address  string
	port     string
	server   *grpc.Server
	listener net.Listener
}

func NewGrpcServer(address, port string) *GrpcServer {
	server := grpc.NewServer()
	reflection.Register(server)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		panic(fmt.Sprintf("failed to listen on %s:%s: %v", address, port, err))
	}

	return &GrpcServer{
		address:  address,
		port:     port,
		server:   server,
		listener: listener,
	}
}

func (g *GrpcServer) Start() {
	log.Printf("[engine] gRPC server listening on %s", g.listener.Addr().String())
	if err := g.server.Serve(g.listener); err != nil {
		panic(fmt.Sprintf("failed to start gRPC server: %v", err))
	}
}

func (g *GrpcServer) Stop() error {
	g.server.GracefulStop()
	return g.listener.Close()
}

func joinHostPort(host, port string) string {
	if host == "" {
		return fmt.Sprintf(":%s", port)
	}
	return fmt.Sprintf("%s:%s", host, port)
}
