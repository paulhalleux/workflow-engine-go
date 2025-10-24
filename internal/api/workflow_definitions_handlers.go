package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/paulhalleux/workflow-engine-go/internal/dto"
	"github.com/paulhalleux/workflow-engine-go/internal/services"
	"github.com/paulhalleux/workflow-engine-go/internal/utils"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/datatypes"
)

type WorkflowDefinitionsHandler struct {
	workflowDefinitionService services.WorkflowDefinitionService
	workflowService           services.WorkflowService
}

func NewWorkflowDefinitionsHandler(
	workflowDefinitionService services.WorkflowDefinitionService,
	workflowService services.WorkflowService,
) *WorkflowDefinitionsHandler {
	return &WorkflowDefinitionsHandler{
		workflowDefinitionService,
		workflowService,
	}
}

func (h *WorkflowDefinitionsHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/workflow-definitions", h.GetAll)
	r.POST("/workflow-definitions/search", h.Search)
	r.GET("/workflow-definitions/:id", h.GetByID)
	r.POST("/workflow-definitions", h.Create)
	r.PUT("/workflow-definitions/:id", h.Update)
	r.DELETE("/workflow-definitions/:id", h.Delete)
	r.PATCH("/workflow-definitions/:id/enable", h.Enable)
	r.PATCH("/workflow-definitions/:id/disable", h.Disable)
	r.POST("/workflow-definitions/:id/_start", h.StartWorkflow)
}

// GetAll
// @Summary Get all workflow definitions
// @Description Retrieve a list of all workflow definitions
// @Tags workflow-definitions
// @Produce json
// @Success 200 {array} models.WorkflowDefinition
// @Router /workflow-definitions [get]
func (h *WorkflowDefinitionsHandler) GetAll(c *gin.Context) {
	wfs, err := h.workflowDefinitionService.GetAllWorkflowDefinitions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wfs)
}

// Search
// @Summary Search workflow definitions
// @Description Search for workflow definitions based on criteria
// @Tags workflow-definitions
// @Accept json
// @Produce json
// @Param search body dto.SearchWorkflowDefinitionsRequest true "Search Criteria"
// @Success 200 {array} models.WorkflowDefinition
// @Router /workflow-definitions/search [post]
func (h *WorkflowDefinitionsHandler) Search(c *gin.Context) {
	var searchReq dto.SearchWorkflowDefinitionsRequest
	if err := c.ShouldBindJSON(&searchReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErrors, err := utils.Validate(&searchReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if validationErrors != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validationErrors": validationErrors})
		return
	}

	wfs, err := h.workflowDefinitionService.SearchWorkflowDefinitions(&searchReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wfs)
}

// GetByID
// @Summary Get a workflow definition by ID
// @Description Retrieve a single workflow definition by its ID
// @Tags workflow-definitions
// @Produce json
// @Param id path string true "Workflow Definition ID"
// @Success 200 {object} models.WorkflowDefinition
// @Router /workflow-definitions/{id} [get]
func (h *WorkflowDefinitionsHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	wf, err := h.workflowDefinitionService.GetWorkflowDefinitionById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workflow not found"})
		return
	}
	c.JSON(http.StatusOK, wf)
}

// Create
// @Summary Create a new workflow definition
// @Description Create a new workflow definition
// @Tags workflow-definitions
// @Accept json
// @Produce json
// @Param workflow body dto.CreateWorkflowDefinitionRequest true "Workflow Definition Data"
// @Success 201 {object} models.WorkflowDefinition
// @Router /workflow-definitions [post]
func (h *WorkflowDefinitionsHandler) Create(c *gin.Context) {
	var wfr dto.CreateWorkflowDefinitionRequest
	if err := c.ShouldBindJSON(&wfr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErrors, err := utils.Validate(&wfr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"validationErrors": validationErrors})
		return
	}

	wf, err := h.workflowDefinitionService.CreateWorkflowDefinition(&wfr)
	if err != nil {
		if utils.IsDuplicateKeyError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "workflow definition with the same name and version already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, wf)
}

// Update
// @Summary Update an existing workflow definition
// @Description Update an existing workflow definition by its ID
// @Tags workflow-definitions
// @Accept json
// @Produce json
// @Param id path string true "Workflow Definition ID"
// @Param workflow body dto.UpdateWorkflowDefinitionRequest true "Workflow Definition Data"
// @Success 200 {object} models.WorkflowDefinition
// @Router /workflow-definitions/{id} [put]
func (h *WorkflowDefinitionsHandler) Update(c *gin.Context) {
	id := c.Param("id")
	wfr := &dto.UpdateWorkflowDefinitionRequest{}
	if err := c.ShouldBindJSON(&wfr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErrors, err := utils.Validate(&wfr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"validationErrors": validationErrors})
		return
	}

	wf, err := h.workflowDefinitionService.UpdateWorkflowDefinition(id, wfr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wf)
}

// Delete
// @Summary Delete a workflow definition
// @Description Delete a workflow definition by its ID
// @Tags workflow-definitions
// @Param id path string true "Workflow Definition ID"
// @Success 204 "No Content"
// @Router /workflow-definitions/{id} [delete]
func (h *WorkflowDefinitionsHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.workflowDefinitionService.DeleteWorkflowDefinition(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// Enable
// @Summary Enable a workflow definition
// @Description Enable a workflow definition by its ID
// @Tags workflow-definitions
// @Param id path string true "Workflow Definition ID"
// @Success 200 {object} models.WorkflowDefinition
// @Router /workflow-definitions/{id}/enable [patch]
func (h *WorkflowDefinitionsHandler) Enable(c *gin.Context) {
	id := c.Param("id")
	wf, err := h.workflowDefinitionService.EnableWorkflowDefinition(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wf)
}

// Disable
// @Summary Disable a workflow definition
// @Description Disable a workflow definition by its ID
// @Tags workflow-definitions
// @Param id path string true "Workflow Definition ID"
// @Success 200 {object} models.WorkflowDefinition
// @Router /workflow-definitions/{id}/disable [patch]
func (h *WorkflowDefinitionsHandler) Disable(c *gin.Context) {
	id := c.Param("id")
	wf, err := h.workflowDefinitionService.DisableWorkflowDefinition(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wf)
}

// StartWorkflow
// @Summary Start a workflow instance from a workflow definition
// @Description Start a new workflow instance based on the specified workflow definition ID
// @Tags workflow-definitions
// @Param id path string true "Workflow Definition ID"
// @Success 201 {object} models.WorkflowInstance
// @Router /workflow-definitions/{id}/_start [post]
func (h *WorkflowDefinitionsHandler) StartWorkflow(c *gin.Context) {
	id := c.Param("id")
	startRequest := &dto.StartWorkflowDefinitionRequest{}
	if err := c.ShouldBindJSON(startRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inputJSON, err := jsonToStructPB(startRequest.Input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input JSON"})
		return
	}

	metadataJSON, err := jsonToStructPB(startRequest.Metadata)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid metadata JSON"})
		return
	}

	if _, err := h.workflowService.StartWorkflow(c.Request.Context(), uuid.MustParse(id), inputJSON, metadataJSON); err != nil {
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func jsonToStructPB(data datatypes.JSON) (*structpb.Struct, error) {
	if len(data) == 0 {
		return structpb.NewStruct(nil)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return structpb.NewStruct(m)
}
