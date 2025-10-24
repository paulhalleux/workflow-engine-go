package main

import (
	"log"
	"net"

	"github.com/gin-gonic/gin"
	_ "github.com/paulhalleux/workflow-engine-go/docs"
	"github.com/paulhalleux/workflow-engine-go/internal/api"
	"github.com/paulhalleux/workflow-engine-go/internal/config"
	"github.com/paulhalleux/workflow-engine-go/internal/container"
	"github.com/paulhalleux/workflow-engine-go/internal/grpcapi"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/proto"
	"github.com/paulhalleux/workflow-engine-go/internal/worker"
	_ "github.com/swaggo/files"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctn := container.NewContainer(config.Default())
	defer ctn.CancelFunc()

	registerTaskExecutors()
	ctn.StepExecutor.Start(ctn.Context)
	ctn.WorkflowExecutor.Start(ctn.Context)

	go startHttpServer(ctn)
	go startGrpcServer(ctn)

	<-ctn.Context.Done()
}

func registerTaskExecutors() {
	worker.RegisterStepExecutor(models.WorkflowStepTypeWait, worker.NewWaitStepExecutor())
}

func startHttpServer(ctn *container.Container) {
	r := gin.Default()
	group := r.Group("/api/v1")

	// Register REST API handlers
	api.NewWorkflowDefinitionsHandler(ctn.WorkflowDefinitionService, ctn.WorkflowService).RegisterRoutes(group)
	api.NewWorkflowInstancesHandler(ctn.WorkflowInstanceService).RegisterRoutes(group)
	api.NewTypeSchemasHandler(ctn.TypeSchemaService).RegisterRoutes(group)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Start the HTTP server
	log.Println("🚀 Engine HTTP server running on ", ctn.Config.HTTPPort)
	serverErr := r.Run(ctn.Config.HTTPPort)
	if serverErr != nil {
		log.Fatalf("failed to run server: %v", serverErr)
	}
}

func startGrpcServer(ctn *container.Container) {
	lis, _ := net.Listen("tcp", ctn.Config.GRPCPort)
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// Register gRPC services
	proto.RegisterWorkflowEngineServer(grpcServer, grpcapi.NewWorkflowEngineServer(ctn.WorkflowService))

	// Start the gRPC server
	log.Println("🚀 Engine gRPC server running on ", ctn.Config.GRPCPort)
	if serveErr := grpcServer.Serve(lis); serveErr != nil {
		log.Fatalf("failed to serve gRPC server: %v", serveErr)
	}
}
