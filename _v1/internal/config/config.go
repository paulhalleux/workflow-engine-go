package config

type Config struct {
	DatabaseURL          string
	QueueBuffer          int
	MaxParallelWorkflows int
	MaxParallelSteps     int
	GRPCPort             string
	HTTPPort             string
}

func Default() *Config {
	return &Config{
		DatabaseURL:          "host=wf-engine-postgres user=db_user password=db_user_password dbname=workflow_engine port=5432 sslmode=disable",
		QueueBuffer:          100,
		MaxParallelWorkflows: 2,
		MaxParallelSteps:     10,
		GRPCPort:             ":50051",
		HTTPPort:             ":8080",
	}
}
