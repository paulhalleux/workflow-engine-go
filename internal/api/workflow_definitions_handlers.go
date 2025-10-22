package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paulhalleux/workflow-engine-go/internal/models"
	"github.com/paulhalleux/workflow-engine-go/internal/persistence"
)

type WorkflowDefinitionsHandler struct {
	repo *persistence.WorkflowDefinitionsRepository
}

func NewWorkflowDefinitionsHandler(repo *persistence.WorkflowDefinitionsRepository) *WorkflowDefinitionsHandler {
	return &WorkflowDefinitionsHandler{repo: repo}
}

func (h *WorkflowDefinitionsHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/workflow-definitions", h.GetAll)
	r.GET("/workflow-definitions/:id", h.GetByID)
	r.POST("/workflow-definitions", h.Create)
	r.PUT("/workflow-definitions/:id", h.Update)
	r.DELETE("/workflow-definitions/:id", h.Delete)
}

func (h *WorkflowDefinitionsHandler) GetAll(c *gin.Context) {
	wfs, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wfs)
}

func (h *WorkflowDefinitionsHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	wf, err := h.repo.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workflow not found"})
		return
	}
	c.JSON(http.StatusOK, wf)
}

func (h *WorkflowDefinitionsHandler) Create(c *gin.Context) {
	var wf models.WorkflowDefinition
	if err := c.ShouldBindJSON(&wf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Create(&wf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, wf)
}

func (h *WorkflowDefinitionsHandler) Update(c *gin.Context) {
	id := c.Param("id")
	wf, err := h.repo.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workflow not found"})
		return
	}
	if err := c.ShouldBindJSON(wf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.Update(wf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wf)
}

func (h *WorkflowDefinitionsHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
