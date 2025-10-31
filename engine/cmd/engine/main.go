package main

import "github.com/paulhalleux/workflow-engine-go/engine"

func main() {
	eng := engine.NewEngine(&engine.WorkflowEngineConfig{
		GrpcAddress: nil,
		GrpcPort:    "50051",
	})
	eng.Start()
}
