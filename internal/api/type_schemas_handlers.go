package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paulhalleux/workflow-engine-go/internal/dto"
	"github.com/paulhalleux/workflow-engine-go/internal/services"
	"github.com/paulhalleux/workflow-engine-go/internal/utils"
)

type TypeSchemasHandler struct {
	svc services.TypeSchemaService
}

func NewTypeSchemasHandler(svc services.TypeSchemaService) *TypeSchemasHandler {
	return &TypeSchemasHandler{svc: svc}
}

func (h *TypeSchemasHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/type-schemas", h.GetAll)
	r.GET("/type-schemas/:id", h.GetById)
	r.POST("/type-schemas", h.Create)
	r.DELETE("/type-schemas/:id", h.Delete)
}

func (h *TypeSchemasHandler) GetAll(c *gin.Context) {
	wfs, err := h.svc.GetAllTypeSchemas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wfs)
}

func (h *TypeSchemasHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	wf, err := h.svc.GetTypeSchemaById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wf)
}

func (h *TypeSchemasHandler) Create(c *gin.Context) {
	var createTypeSchemaRequest dto.CreateTypeSchemaRequest
	if err := c.ShouldBindJSON(&createTypeSchemaRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErrors, err := utils.Validate(&createTypeSchemaRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"validationErrors": validationErrors})
		return
	}

	wf, err := h.svc.CreateTypeSchema(&createTypeSchemaRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, wf)
}

func (h *TypeSchemasHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteTypeSchema(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
