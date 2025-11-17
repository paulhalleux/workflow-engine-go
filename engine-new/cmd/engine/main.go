package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/paulhalleux/workflow-engine-go/engine-new/internal/app"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("warning: unable to load .env file: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := app.LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	engine, err := app.NewEngine(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialize engine: %v", err)
	}

	if err := engine.Start(); err != nil {
		log.Printf("engine stopped with error: %v", err)
	}

	log.Println("engine terminated")
}
