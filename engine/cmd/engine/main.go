package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/paulhalleux/workflow-engine-go/engine"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	eng := engine.NewEngine(&engine.WorkflowEngineConfig{
		GrpcPort: os.Getenv("GRPC_PORT"),
		HttpPort: os.Getenv("HTTP_PORT"),

		DbName:     os.Getenv("DB_NAME"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbSSLMode:  os.Getenv("DB_SSLMODE"),
	})

	eng.Start()
}
