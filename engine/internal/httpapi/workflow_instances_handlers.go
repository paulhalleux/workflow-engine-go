package httpapi

import (
	"github.com/gin-gonic/gin"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/service"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
)

type WorkflowInstancesHandlers struct {
	service *service.WorkflowService
}

func NewWorkflowInstancesHandlers(
	service *service.WorkflowService,
) *WorkflowInstancesHandlers {
	return &WorkflowInstancesHandlers{
		service: service,
	}
}

func (h *WorkflowInstancesHandlers) Register(rg *gin.RouterGroup) {
	rg.GET("/workflow-instances", h.getAll)
	rg.GET("/workflow-instances/:id", h.getByID)
	//rg.POST("/workflow-definitions", h.create)
	//rg.PUT("/workflow-definitions/:id", h.update)
	//rg.DELETE("/workflow-definitions/:id", h.delete)
	//rg.PATCH("/workflow-definitions/:id/enable", h.enable)
	//rg.PATCH("/workflow-definitions/:id/disable", h.disable)
	//rg.PATCH("/workflow-definitions/:id/publish", h.publish)
}

func (h *WorkflowInstancesHandlers) getAll(ctx *gin.Context) {
	paginate := utils.Paginate(ctx)
	defs, err := h.service.GetAll(&paginate)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, defs)
}

func (h *WorkflowInstancesHandlers) getByID(ctx *gin.Context) {
	id := ctx.Param("id")
	def, err := h.service.GetByID(id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, def)
}

//func (h *WorkflowInstancesHandlers) create(ctx *gin.Context) {
//	var def *models.WorkflowInstance
//	if err := ctx.ShouldBindJSON(&def); err != nil {
//		ctx.JSON(400, gin.H{"error": err.Error()})
//		return
//	}
//	err := h.service.Create(def)
//	if err != nil {
//		ctx.JSON(500, gin.H{"error": err.Error()})
//		return
//	}
//	ctx.JSON(201, def)
//}
//
//func (h *WorkflowInstancesHandlers) update(ctx *gin.Context) {
//	id := ctx.Param("id")
//	var updateRequest *models.UpdateWorkflowDefRequest
//	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
//		ctx.JSON(400, gin.H{"error": err.Error()})
//		return
//	}
//	def, err := h.service.GetByID(id)
//	if err != nil {
//		ctx.JSON(500, gin.H{"error": err.Error()})
//		return
//	}
//	updatedDef := updateRequest.ToWorkflowInstance(def)
//	err = h.service.Update(updatedDef)
//	if err != nil {
//		ctx.JSON(500, gin.H{"error": err.Error()})
//		return
//	}
//	ctx.JSON(200, updatedDef)
//}
//
//func (h *WorkflowInstancesHandlers) delete(ctx *gin.Context) {
//	err := h.service.Delete(ctx.Param("id"))
//	if err != nil {
//		ctx.JSON(500, gin.H{"error": err.Error()})
//		return
//	}
//	ctx.Status(204)
//}
