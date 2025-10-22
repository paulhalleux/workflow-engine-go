package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paulhalleux/workflow-engine-go/internal/api/dto"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/persistence"
	"github.com/paulhalleux/workflow-engine-go/internal/utils"
)

var validate, _ = utils.CreateValidator()

type WorkflowDefinitionsHandler struct {
	repo *persistence.WorkflowDefinitionsRepository
}

func NewWorkflowDefinitionsHandler(repo *persistence.WorkflowDefinitionsRepository) *WorkflowDefinitionsHandler {
	return &WorkflowDefinitionsHandler{repo: repo}
}

func (h *WorkflowDefinitionsHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/workflow-definitions", h.GetAll)
	r.GET("/workflow-definitions/:id", h.GetByID)
	r.POST("/workflow-definitions", h.Create)
	r.PUT("/workflow-definitions/:id", h.Update)
	r.DELETE("/workflow-definitions/:id", h.Delete)
	r.PATCH("/workflow-definitions/:id/enable", h.Enable)
	r.PATCH("/workflow-definitions/:id/disable", h.Disable)
}

// GetAll
// @Summary Get all workflow definitions
// @Description Retrieve a list of all workflow definitions
// @Tags workflow-definitions
// @Produce json
// @Success 200 {array} models.WorkflowDefinition
// @Router /workflow-definitions [get]
func (h *WorkflowDefinitionsHandler) GetAll(c *gin.Context) {
	wfs, err := h.repo.GetAll()
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
	wf, err := h.repo.GetById(id)
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

	wf := models.WorkflowDefinition{
		Name:        wfr.Name,
		Description: wfr.Description,
		Version:     wfr.Version,
		Metadata:    wfr.Metadata,
		Steps:       &wfr.Steps,
		IsEnabled:   wfr.IsEnabled,
	}

	if err := h.repo.Create(&wf); err != nil {
		if h.repo.IsDuplicateKeyError(err) {
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
	wf, err := h.repo.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workflow not found"})
		return
	}

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

	updatePartialWorkflowDefinition(wf, wfr)
	if err := h.repo.Update(wf); err != nil {
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
	if err := h.repo.Delete(id); err != nil {
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
	wf, err := h.repo.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workflow not found"})
		return
	}

	wf.IsEnabled = true

	if err := h.repo.Update(wf); err != nil {
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
	wf, err := h.repo.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workflow not found"})
		return
	}

	wf.IsEnabled = false

	if err := h.repo.Update(wf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wf)
}

func updatePartialWorkflowDefinition(wf *models.WorkflowDefinition, wfr *dto.UpdateWorkflowDefinitionRequest) {
	if wfr.Name != nil {
		wf.Name = *wfr.Name
	}
	if wfr.Description != nil {
		wf.Description = *wfr.Description
	}
	if wfr.Version != nil {
		wf.Version = *wfr.Version
	}
	if wfr.Metadata != nil {
		wf.Metadata = *wfr.Metadata
	}
	if wfr.Steps != nil {
		wf.Steps = wfr.Steps
	}
	if wfr.IsEnabled != nil {
		wf.IsEnabled = *wfr.IsEnabled
	}
}
