package main

import (
	"github.com/paulhalleux/workflow-engine-go/internal/agent"
)

func main() {
	ag := agent.NewAgent("Echo", nil, 50052)
	ag.RegisterTask("echo", NewEchoTask())
	ag.Start()
}
