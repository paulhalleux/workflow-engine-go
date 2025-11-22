package app

import (
	"fmt"
	"os"
)

type Config struct {
	GrpcAddress string
	GrpcPort    string

	HttpAddress string
	HttpPort    string

	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	DbSSLMode  string
}

func LoadConfigFromEnv() (*Config, error) {
	cfg := &Config{
		GrpcAddress: getEnvDefault("GRPC_ADDRESS", ""),
		GrpcPort:    getEnvDefault("GRPC_PORT", "50051"),

		HttpAddress: getEnvDefault("HTTP_ADDRESS", ""),
		HttpPort:    getEnvDefault("HTTP_PORT", "8080"),

		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
		DbSSLMode:  getEnvDefault("DB_SSLMODE", "disable"),
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
	if c.HttpPort == "" {
		missing = append(missing, "HTTP_PORT")
	}
	if c.DbHost == "" {
		missing = append(missing, "DB_HOST")
	}
	if c.DbPort == "" {
		missing = append(missing, "DB_PORT")
	}
	if c.DbUser == "" {
		missing = append(missing, "DB_USER")
	}
	if c.DbName == "" {
		missing = append(missing, "DB_NAME")
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
