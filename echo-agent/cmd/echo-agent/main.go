package main

import (
	"context"

	"github.com/paulhalleux/workflow-engine-go/agent"
	"github.com/paulhalleux/workflow-engine-go/echo-agent/internal"
)

func main() {
	ctx := context.Background()

	ag, err := agent.NewAgent(ctx, &agent.Config{
		Name:          "echo-agent",
		Version:       "1.0.0",
		GrpcPort:      "50052",
		EngineGrpcUrl: ":50051",

		MaxQueueSize:     100,
		MaxParallelTasks: 10,
	})

	if err != nil {
		panic(err)
	}

	ag.RegisterTaskDefinition(internal.NewEchoTaskDefinition())

	err = ag.Start()
	if err != nil {
		return
	}
}
