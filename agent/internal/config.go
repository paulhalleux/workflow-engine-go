package internal

type WorkflowAgentConfig struct {
	Name          string
	Version       string
	GrpcPort      string
	GrpcAddress   *string
	EngineGrpcUrl string

	MaxQueueSize     int
	MaxParallelTasks int
}
