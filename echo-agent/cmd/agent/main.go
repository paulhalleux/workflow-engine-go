package main

import "github.com/paulhalleux/workflow-engine-go/agent"

func main() {
	ag := agent.NewWorkflowAgent(&agent.WorkflowAgentConfig{
		Name: "echo-agent",
		Port: "50052",
	})
	ag.Start()
}
