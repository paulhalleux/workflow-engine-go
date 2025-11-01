package httpapi

import (
	"github.com/gin-gonic/gin"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/models"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/service"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
)

type WorkflowDefinitionsHandlers struct {
	service *service.WorkflowDefinitionsService
}

func NewWorkflowDefinitionsHandlers(
	service *service.WorkflowDefinitionsService,
) *WorkflowDefinitionsHandlers {
	return &WorkflowDefinitionsHandlers{
		service: service,
	}
}

func (h *WorkflowDefinitionsHandlers) Register(rg *gin.RouterGroup) {
	rg.GET("/workflow-definitions", h.getAll)
	rg.GET("/workflow-definitions/:id", h.getByID)
	rg.POST("/workflow-definitions", h.create)
	rg.PUT("/workflow-definitions/:id", h.update)
	rg.DELETE("/workflow-definitions/:id", h.delete)
	rg.PATCH("/workflow-definitions/:id/enable", h.enable)
	rg.PATCH("/workflow-definitions/:id/disable", h.disable)
	rg.PATCH("/workflow-definitions/:id/publish", h.publish)
}

func (h *WorkflowDefinitionsHandlers) getAll(ctx *gin.Context) {
	paginate := utils.Paginate(ctx)
	defs, err := h.service.GetAll(&paginate)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, defs)
}

func (h *WorkflowDefinitionsHandlers) getByID(ctx *gin.Context) {
	id := ctx.Param("id")
	def, err := h.service.GetByID(id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, def)
}

func (h *WorkflowDefinitionsHandlers) create(ctx *gin.Context) {
	var def *models.WorkflowDefinition
	if err := ctx.ShouldBindJSON(&def); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := h.service.Create(def)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, def)
}

func (h *WorkflowDefinitionsHandlers) update(ctx *gin.Context) {
	id := ctx.Param("id")
	var updateRequest *models.UpdateWorkflowDefRequest
	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	def, err := h.service.GetByID(id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	updatedDef := updateRequest.ToWorkflowDefinition(def)
	err = h.service.Update(updatedDef)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, updatedDef)
}

func (h *WorkflowDefinitionsHandlers) delete(ctx *gin.Context) {
	err := h.service.Delete(ctx.Param("id"))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(204)
}

func (h *WorkflowDefinitionsHandlers) enable(ctx *gin.Context) {
	err := h.service.Enable(ctx.Param("id"))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(200)
}

func (h *WorkflowDefinitionsHandlers) disable(ctx *gin.Context) {
	err := h.service.Disable(ctx.Param("id"))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(200)
}

func (h *WorkflowDefinitionsHandlers) publish(ctx *gin.Context) {
	err := h.service.Publish(ctx.Param("id"))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(200)
}
