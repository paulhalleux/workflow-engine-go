package agent

import (
	"fmt"
	"os"
)

type Config struct {
	Name          string
	Version       string
	EngineGrpcUrl string

	MaxQueueSize     int
	MaxParallelTasks int

	GrpcAddress string
	GrpcPort    string
}

func LoadConfigFromEnv() (*Config, error) {
	cfg := &Config{
		GrpcAddress:   getEnvDefault("GRPC_ADDRESS", ""),
		GrpcPort:      getEnvDefault("GRPC_PORT", "50051"),
		Name:          getEnvDefault("AGENT_NAME", "workflow-agent"),
		Version:       getEnvDefault("AGENT_VERSION", "v1.0.0"),
		EngineGrpcUrl: getEnvDefault("ENGINE_GRPC_URL", "localhost:60051"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	missing := make([]string, 0)
	if c.GrpcPort == "" {
		missing = append(missing, "GRPC_PORT")
	}
	if c.EngineGrpcUrl == "" {
		missing = append(missing, "ENGINE_GRPC_URL")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %v", missing)
	}

	return nil
}

func getEnvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
