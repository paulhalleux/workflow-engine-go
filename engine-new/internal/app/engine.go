package app

import (
	"context"
	"fmt"
	"log"

	"github.com/paulhalleux/workflow-engine-go/engine-new/internal/grpcserver"
	"github.com/paulhalleux/workflow-engine-go/engine-new/internal/httpserver"
	"github.com/paulhalleux/workflow-engine-go/engine-new/internal/persistance"
	"github.com/paulhalleux/workflow-engine-go/engine-new/internal/ws"
	"github.com/paulhalleux/workflow-engine-go/proto"
	"gorm.io/gorm"
)

type Engine struct {
	cfg *Config
	ctx context.Context
	db  *gorm.DB
}

func NewEngine(ctx context.Context, cfg *Config) (*Engine, error) {
	db, err := createDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	return &Engine{
		cfg: cfg,
		ctx: ctx,
		db:  db,
	}, nil
}

func (e *Engine) Start() error {
	httpSrv := httpserver.NewHttpServer(
		e.cfg.HttpAddress,
		e.cfg.HttpPort,
	)

	grpcSrv := grpcserver.NewGrpcServer(
		e.cfg.GrpcAddress,
		e.cfg.GrpcPort,
		grpcserver.NewEngineService(),
		grpcserver.NewTaskService(),
	)

	wsSrv := ws.NewServer()
	wsSrv.Registry.RegisterCommand(proto.WEBSOCKET_COMMAND_TYPE_SUBSCRIBE, ws.NewSubscribeCommandHandler())

	wfDefRepo := persistance.NewWorkflowDefinitionRepository(e.db)
	wfDefHandlers := httpserver.NewWorkflowDefinitionsHandlers(wfDefRepo)

	httpSrv.RegisterApiHandler(wfDefHandlers)

	// Lancer les serveurs en goroutines.
	go httpSrv.Start(wsSrv)
	go grpcSrv.Start()

	// Attendre la fin du contexte.
	<-e.ctx.Done()
	log.Println("[engine] shutting down...")

	if err := httpSrv.Stop(); err != nil {
		log.Printf("http shutdown error: %v", err)
	}

	if err := grpcSrv.Stop(); err != nil {
		log.Printf("grpc shutdown error: %v", err)
	}

	return nil
}

func createDatabase(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbSSLMode,
	)
	return persistance.CreateDatabase(dsn)
}
