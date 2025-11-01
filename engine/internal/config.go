package internal

type WorkflowEngineConfig struct {
	GrpcPort             string
	GrpcAddress          *string
	HttpPort             string
	HttpAddress          *string
	DbHost               string
	DbPort               string
	DbUser               string
	DbPassword           string
	DbName               string
	DbSSLMode            string
	MaxWorkflowQueueSize int
	MaxParallelWorkflows int
	MaxStepQueueSize     int
	MaxParallelSteps     int
}
