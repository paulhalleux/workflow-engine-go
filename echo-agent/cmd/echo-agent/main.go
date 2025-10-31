package main

import (
	"github.com/paulhalleux/workflow-engine-go/agent"
	"github.com/paulhalleux/workflow-engine-go/echo-agent/internal"
)

func main() {
	ag := agent.NewWorkflowAgent(&agent.WorkflowAgentConfig{
		Name:          "echo-agent",
		Version:       "1.0.0",
		GrpcPort:      "50052",
		EngineGrpcUrl: ":50051",
	})

	ag.RegisterTaskDefinition(internal.NewEchoTaskDefinition())
	ag.Start()
}
