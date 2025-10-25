package main

import (
	"github.com/paulhalleux/workflow-engine-go/internal/agent"
)

func main() {
	ag := agent.Agent{
		Name: "Echo",
		Grpc: agent.GrpcInfo{
			Port: 50052,
		},
	}

	ag.RegisterTask("echo", NewEchoTask())

	ag.Start()
}
