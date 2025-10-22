package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/paulhalleux/workflow-engine-go/internal/api"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/persistence"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	repo := persistence.NewWorkflowRepository(db)
	handler := api.NewWorkflowDefinitionsHandler(repo)

	r := gin.Default()
	handler.RegisterRoutes(r)

	log.Println("🚀 Engine HTTP server running on :8080")
	serverErr := r.Run(":8080")
	if serverErr != nil {
		log.Fatalf("failed to run server: %v", serverErr)
	}
}
