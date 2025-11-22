package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/paulhalleux/workflow-engine-go/engine-new/internal/models"
	"github.com/paulhalleux/workflow-engine-go/engine-new/internal/persistance"
	"github.com/paulhalleux/workflow-engine-go/utils/expr"
	"github.com/paulhalleux/workflow-engine-go/utils/pagination"
	"github.com/paulhalleux/workflow-engine-go/utils/semver"
)

type WorkflowDefinitionsHandlers struct {
	repo persistance.WorkflowDefinitionRepository
}

func NewWorkflowDefinitionsHandlers(
	repo persistance.WorkflowDefinitionRepository,
) *WorkflowDefinitionsHandlers {
	return &WorkflowDefinitionsHandlers{
		repo: repo,
	}

}

func (w *WorkflowDefinitionsHandlers) Register(router gin.IRoutes) {
	router.GET("/workflow-definitions", w.GetAllWorkflowDefinitions)
	router.GET("/workflow-definitions/:id", w.GetWorkflowDefinitionByID)
	router.POST("/workflow-definitions/search", w.SearchWorkflowDefinitions)
	router.POST("/workflow-definitions", w.CreateWorkflowDefinition)
	router.PUT("/workflow-definitions/:id", w.UpdateWorkflowDefinition)
	router.DELETE("/workflow-definitions/:id", w.DeleteWorkflowDefinition)
	router.PATCH("/workflow-definitions/:id/publish", w.PublishWorkflowDefinition)
	router.PATCH("/workflow-definitions/:id/enable", w.EnableWorkflowDefinition)
	router.PATCH("/workflow-definitions/:id/disable", w.DisableWorkflowDefinition)
}

func (w *WorkflowDefinitionsHandlers) GetAllWorkflowDefinitions(c *gin.Context) {
	var paginationParams pagination.Pagination
	if err := c.ShouldBindQuery(&paginationParams); err != nil {
		c.JSON(400, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	definitions, err := w.repo.GetAll(paginationParams)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve workflow definitions"})
		return
	}

	c.JSON(200, definitions)
}

func (w *WorkflowDefinitionsHandlers) SearchWorkflowDefinitions(c *gin.Context) {
	var expression expr.Expression
	if err := c.ShouldBindJSON(&expression); err != nil {
		c.JSON(400, gin.H{"error": "Invalid search expression"})
		return
	}

	var paginationParams pagination.Pagination
	if err := c.ShouldBindQuery(&paginationParams); err != nil {
		c.JSON(400, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	definitions, err := w.repo.Search(expression, paginationParams)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to search workflow definitions"})
		return
	}

	c.JSON(200, definitions)
}

func (w *WorkflowDefinitionsHandlers) GetWorkflowDefinitionByID(c *gin.Context) {
	id := c.Param("id")
	definition, err := w.repo.GetByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve workflow definition"})
		return
	}
	if definition == nil {
		c.JSON(404, gin.H{"error": "Workflow definition not found"})
		return
	}
	c.JSON(200, definition)
}

func (w *WorkflowDefinitionsHandlers) CreateWorkflowDefinition(c *gin.Context) {
	var params struct {
		IsDraft bool `form:"draft,default=true"`
	}

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(400, gin.H{"error": "Invalid query parameters"})
		return
	}

	var definition models.WorkflowDefinition
	if err := c.ShouldBindJSON(&definition); err != nil {
		c.JSON(400, gin.H{"error": "Invalid workflow definition data"})
		return
	}

	version, err := semver.Parse(semver.InitialVersion())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse initial version"})
		return
	}

	if params.IsDraft {
		version.PreRelease = "1"
	}

	definition.Version = version.String()
	createdDefinition, err := w.repo.Create(&definition)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create workflow definition"})
		return
	}

	c.JSON(201, createdDefinition)
}

func (w *WorkflowDefinitionsHandlers) UpdateWorkflowDefinition(c *gin.Context) {
	id := c.Param("id")
	var definition models.WorkflowDefinition
	if err := c.ShouldBindJSON(&definition); err != nil {
		c.JSON(400, gin.H{"error": "Invalid workflow definition data"})
		return
	}

	uuidId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid workflow definition ID"})
		return
	}

	definition.ID = uuidId
	updatedDefinition, err := w.repo.Update(&definition)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update workflow definition"})
		return
	}

	c.JSON(200, updatedDefinition)
}

func (w *WorkflowDefinitionsHandlers) DeleteWorkflowDefinition(c *gin.Context) {
	id := c.Param("id")
	err := w.repo.Delete(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete workflow definition"})
		return
	}
	c.Status(204)
}

func (w *WorkflowDefinitionsHandlers) PublishWorkflowDefinition(c *gin.Context) {
	id := c.Param("id")
	definition, err := w.repo.GetByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve workflow definition"})
		return
	}

	version, err := semver.Parse(definition.Version)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse workflow definition version"})
		return
	}

	if !version.IsDraft() {
		c.JSON(400, gin.H{"error": "Workflow definition is already published"})
		return
	}

	version.ReleaseDraft()
	definition.Version = version.String()

	_, err = w.repo.Update(definition)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to publish workflow definition"})
		return
	}

	c.Status(200)
}

func (w *WorkflowDefinitionsHandlers) EnableWorkflowDefinition(c *gin.Context) {
	id := c.Param("id")
	err := w.repo.Enable(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to enable workflow definition"})
		return
	}
	c.Status(200)
}

func (w *WorkflowDefinitionsHandlers) DisableWorkflowDefinition(c *gin.Context) {
	id := c.Param("id")
	err := w.repo.Disable(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to disable workflow definition"})
		return
	}
	c.Status(200)
}
