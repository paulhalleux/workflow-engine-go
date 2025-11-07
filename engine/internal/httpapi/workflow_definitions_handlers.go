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

// getAll
// @Id 			getAllWorkflowDefinitions
// @Summary 	Get all workflow definitions
// @Description Get all workflow definitions
// @Tags 		Workflow Definitions
// @Produce 	json
// @Success 	200 {array} models.WorkflowDefinition
// @Router 		/workflow-definitions [get]
func (h *WorkflowDefinitionsHandlers) getAll(ctx *gin.Context) {
	paginate := utils.Paginate(ctx)
	defs, err := h.service.GetAll(&paginate)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, defs)
}

// getByID
// @Id 			getWorkflowDefinitionByID
// @Summary 	Get a workflow definition by ID
// @Description Get a workflow definition by ID
// @Tags 		Workflow Definitions
// @Produce 	json
// @Param 		id path string true "Workflow Definition ID"
// @Success 	200 {object} models.WorkflowDefinition
// @Router 		/workflow-definitions/{id} [get]
func (h *WorkflowDefinitionsHandlers) getByID(ctx *gin.Context) {
	id := ctx.Param("id")
	def, err := h.service.GetByID(id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, def)
}

// create
// @Id 			createWorkflowDefinition
// @Summary 	Create a new workflow definition
// @Description Create a new workflow definition
// @Tags 		Workflow Definitions
// @Accept 		json
// @Produce 	json
// @Param 		definition body models.WorkflowDefinition true "Workflow Definition"
// @Success 	201 {object} models.WorkflowDefinition
// @Router 		/workflow-definitions [post]
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

// update
// @Id 			updateWorkflowDefinition
// @Summary 	Update an existing workflow definition
// @Description Update an existing workflow definition
// @Tags 		Workflow Definitions
// @Accept 		json
// @Produce 	json
// @Param 		id path string true "Workflow Definition ID"
// @Param 		definition body models.UpdateWorkflowDefRequest true "Updated Workflow Definition"
// @Success 	200 {object} models.WorkflowDefinition
// @Router 		/workflow-definitions/{id} [put]
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

// delete
// @Id 			deleteWorkflowDefinition
// @Summary 	Delete a workflow definition
// @Description Delete a workflow definition
// @Tags 		Workflow Definitions
// @Param 		id path string true "Workflow Definition ID"
// @Success 	204
// @Router 		/workflow-definitions/{id} [delete]
func (h *WorkflowDefinitionsHandlers) delete(ctx *gin.Context) {
	err := h.service.Delete(ctx.Param("id"))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(204)
}

// delete
// @Id 			enableWorkflowDefinition
// @Summary 	Delete a workflow definition
// @Description Delete a workflow definition
// @Tags 		Workflow Definitions
// @Param 		id path string true "Workflow Definition ID"
// @Success 	204
// @Router 		/workflow-definitions/{id}/enable [patch]
func (h *WorkflowDefinitionsHandlers) enable(ctx *gin.Context) {
	err := h.service.Enable(ctx.Param("id"))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(200)
}

// delete
// @Id 			disableWorkflowDefinition
// @Summary 	Delete a workflow definition
// @Description Delete a workflow definition
// @Tags 		Workflow Definitions
// @Param 		id path string true "Workflow Definition ID"
// @Success 	204
// @Router 		/workflow-definitions/{id}/disable [patch]
func (h *WorkflowDefinitionsHandlers) disable(ctx *gin.Context) {
	err := h.service.Disable(ctx.Param("id"))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(200)
}

// publish
// @Id 			publishWorkflowDefinition
// @Summary 	Publish a workflow definition
// @Description Publish a workflow definition
// @Tags 		Workflow Definitions
// @Param 		id path string true "Workflow Definition ID"
// @Success 	204
// @Router 		/workflow-definitions/{id}/publish [patch]
func (h *WorkflowDefinitionsHandlers) publish(ctx *gin.Context) {
	err := h.service.Publish(ctx.Param("id"))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(200)
}
