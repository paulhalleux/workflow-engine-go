package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paulhalleux/workflow-engine-go/internal/dto"
	"github.com/paulhalleux/workflow-engine-go/internal/services"
	"github.com/paulhalleux/workflow-engine-go/internal/utils"
)

type WorkflowInstancesHandler struct {
	svc services.WorkflowInstanceService
}

func NewWorkflowInstancesHandler(svc services.WorkflowInstanceService) *WorkflowInstancesHandler {
	return &WorkflowInstancesHandler{svc: svc}
}

func (h *WorkflowInstancesHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/workflow-instances", h.GetAll)
	r.GET("/workflow-instances/:id", h.GetByID)
	r.POST("/workflow-instances", h.Create)
	r.PUT("/workflow-instances/:id", h.Update)
	r.DELETE("/workflow-instances/:id", h.Delete)
}

// GetAll
// @Summary Get all workflow instances
// @Description Retrieve a list of all workflow instances
// @Tags workflow-instances
// @Produce json
// @Success 200 {array} models.WorkflowInstance
// @Router /workflow-instances [get]
func (h *WorkflowInstancesHandler) GetAll(c *gin.Context) {
	wfs, err := h.svc.GetAllWorkflowInstances()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wfs)
}

// GetByID
// @Summary Get a workflow instance by ID
// @Description Retrieve a single workflow instance by its ID
// @Tags workflow-instances
// @Produce json
// @Param id path string true "Workflow Instance ID"
// @Success 200 {object} models.WorkflowInstance
// @Router /workflow-instances/{id} [get]
func (h *WorkflowInstancesHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	wf, err := h.svc.GetWorkflowInstanceById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workflow not found"})
		return
	}
	c.JSON(http.StatusOK, wf)
}

// Create
// @Summary Create a new workflow instance
// @Description Create a new workflow instance
// @Tags workflow-instances
// @Accept json
// @Produce json
// @Param workflow body dto.CreateWorkflowInstanceRequest true "Workflow Instance Data"
// @Success 201 {object} models.WorkflowInstance
// @Router /workflow-instances [post]
func (h *WorkflowInstancesHandler) Create(c *gin.Context) {
	var wfr dto.CreateWorkflowInstanceRequest
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

	wf, err := h.svc.CreateWorkflowInstance(&wfr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, wf)
}

// Update
// @Summary Update an existing workflow instance
// @Description Update an existing workflow instance by its ID
// @Tags workflow-instances
// @Accept json
// @Produce json
// @Param id path string true "Workflow Instance ID"
// @Param workflow body dto.UpdateWorkflowInstanceRequest true "Workflow Instance Data"
// @Success 200 {object} models.WorkflowInstance
// @Router /workflow-instances/{id} [put]
func (h *WorkflowInstancesHandler) Update(c *gin.Context) {
	id := c.Param("id")

	wfr := &dto.UpdateWorkflowInstanceRequest{}
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

	wf, err := h.svc.UpdateWorkflowInstance(id, wfr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wf)
}

// Delete
// @Summary Delete a workflow instance
// @Description Delete a workflow instance by its ID
// @Tags workflow-instances
// @Param id path string true "Workflow Instance ID"
// @Success 204 "No Content"
// @Router /workflow-instances/{id} [delete]
func (h *WorkflowInstancesHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteWorkflowInstance(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
