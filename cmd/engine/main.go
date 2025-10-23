package main

import (
	"context"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/paulhalleux/workflow-engine-go/internal/api"
	"github.com/paulhalleux/workflow-engine-go/internal/grpcapi"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/persistence"
	"github.com/paulhalleux/workflow-engine-go/internal/proto"
	"github.com/paulhalleux/workflow-engine-go/internal/queue"
	"github.com/paulhalleux/workflow-engine-go/internal/services"
	"github.com/paulhalleux/workflow-engine-go/internal/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/paulhalleux/workflow-engine-go/docs"
	_ "github.com/swaggo/files"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	dsn := "host=wf-engine-postgres user=db_user password=db_user_password dbname=workflow_engine port=5432 sslmode=disable"
	db, dbErr := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		log.Fatalf("failed to connect database: %v", dbErr)
	}

	migrateErr := db.AutoMigrate(&models.WorkflowDefinition{})
	if migrateErr != nil {
		log.Fatalf("failed to migrate database: %v", migrateErr)
	}

	wfdRepo := persistence.NewWorkflowDefinitionsRepository(db)
	wfiRepo := persistence.NewWorkflowInstancesRepository(db)

	lis, _ := net.Listen("tcp", ":50051")
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	wfQueue := queue.NewMemoryQueue(100)

	workflowSvc := services.NewWorkflowService(wfdRepo, wfiRepo, wfQueue)
	executor := worker.NewWorkflowExecutor(wfdRepo, wfiRepo, wfQueue, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	executor.Start(ctx)

	// Register gRPC services
	proto.RegisterWorkflowEngineServer(grpcServer, grpcapi.NewWorkflowEngineServer(workflowSvc))

	// Start the gRPC server
	go func() {
		log.Println("🚀 Engine gRPC server running on :50051")
		if serveErr := grpcServer.Serve(lis); serveErr != nil {
			log.Fatalf("failed to serve gRPC server: %v", serveErr)
		}
	}()

	r := gin.Default()
	group := r.Group("/api/v1")

	// Register REST API handlers
	api.NewWorkflowDefinitionsHandler(wfdRepo).RegisterRoutes(group)
	api.NewWorkflowInstancesHandler(wfiRepo).RegisterRoutes(group)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Start the HTTP server
	log.Println("🚀 Engine HTTP server running on :8080")
	serverErr := r.Run(":8080")
	if serverErr != nil {
		log.Fatalf("failed to run server: %v", serverErr)
	}
}
