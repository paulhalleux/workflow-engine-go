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

// GetAllWorkflowDefinitions godoc
// @ID           GetAllWorkflowDefinitions
// @Summary      Get all workflow definitions
// @Description  Retrieve a paginated list of all workflow definitions
// @Tags         Workflow Definitions
// @Accept       json
// @Produce      json
// @Param        page     query    int     false  "Page number"
// @Param        pageSize query    int     false  "Number of items per page"
// @Success      200  {array}   models.WorkflowDefinition
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /api/workflow-definitions [get]
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

// SearchWorkflowDefinitions godoc
// @ID           SearchWorkflowDefinitions
// @Summary      Search workflow definitions
// @Description  Search for workflow definitions based on a given expression
// @Tags         Workflow Definitions
// @Accept       json
// @Produce      json
// @Param        expression  body     expr.Expression  true  "Search expression"
// @Param        page        query    int              false "Page number"
// @Param        pageSize    query    int              false "Number of items per page"
// @Success      200  {array}   models.WorkflowDefinition
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /api/workflow-definitions/search [post]
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

// GetWorkflowDefinitionByID godoc
// @ID           GetWorkflowDefinitionByID
// @Summary      Get workflow definition by ID
// @Description  Retrieve a workflow definition by its ID
// @Tags         Workflow Definitions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Workflow Definition ID"
// @Success      200  {object}  models.WorkflowDefinition
// @Failure      400  {object}  gin.H
// @Failure      404  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /api/workflow-definitions/{id} [get]
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

// CreateWorkflowDefinition godoc
// @ID           CreateWorkflowDefinition
// @Summary      Create a new workflow definition
// @Description  Create a new workflow definition, optionally as a draft
// @Tags         Workflow Definitions
// @Accept       json
// @Produce      json
// @Param        draft  query     bool    false  "Create as draft"  default(true)
// @Param        body   body      models.WorkflowDefinition  true  "Workflow Definition Data"
// @Success      201  {object}  models.WorkflowDefinition
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /api/workflow-definitions [post]
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

// UpdateWorkflowDefinition godoc
// @ID           UpdateWorkflowDefinition
// @Summary      Update an existing workflow definition
// @Description  Update an existing workflow definition by its ID
// @Tags         Workflow Definitions
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "Workflow Definition ID"
// @Param        body  body      models.WorkflowDefinition  true  "Workflow Definition Data"
// @Success      200  {object}  models.WorkflowDefinition
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /api/workflow-definitions/{id} [put]
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

// DeleteWorkflowDefinition godoc
// @ID           DeleteWorkflowDefinition
// @Summary      Delete a workflow definition
// @Description  Delete a workflow definition by its ID
// @Tags         Workflow Definitions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Workflow Definition ID"
// @Success      204
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /api/workflow-definitions/{id} [delete]
func (w *WorkflowDefinitionsHandlers) DeleteWorkflowDefinition(c *gin.Context) {
	id := c.Param("id")
	err := w.repo.Delete(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete workflow definition"})
		return
	}
	c.Status(204)
}

// PublishWorkflowDefinition godoc
// @ID           PublishWorkflowDefinition
// @Summary      Publish a workflow definition
// @Description  Publish a draft workflow definition by its ID
// @Tags         Workflow Definitions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Workflow Definition ID"
// @Success      200
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /api/workflow-definitions/{id}/publish [patch]
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

// EnableWorkflowDefinition godoc
// @ID           EnableWorkflowDefinition
// @Summary      Enable a workflow definition
// @Description  Enable a workflow definition by its ID
// @Tags         Workflow Definitions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Workflow Definition ID"
// @Success      200
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /api/workflow-definitions/{id}/enable [patch]
func (w *WorkflowDefinitionsHandlers) EnableWorkflowDefinition(c *gin.Context) {
	id := c.Param("id")
	err := w.repo.Enable(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to enable workflow definition"})
		return
	}
	c.Status(200)
}

// DisableWorkflowDefinition godoc
// @ID           DisableWorkflowDefinition
// @Summary      Disable a workflow definition
// @Description  Disable a workflow definition by its ID
// @Tags         Workflow Definitions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Workflow Definition ID"
// @Success      200
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /api/workflow-definitions/{id}/disable [patch]
func (w *WorkflowDefinitionsHandlers) DisableWorkflowDefinition(c *gin.Context) {
	id := c.Param("id")
	err := w.repo.Disable(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to disable workflow definition"})
		return
	}
	c.Status(200)
}
