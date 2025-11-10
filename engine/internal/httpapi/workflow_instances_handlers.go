package httpapi

import (
	"github.com/gin-gonic/gin"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/service"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
)

type WorkflowInstancesHandlers struct {
	service *service.WorkflowInstanceService
}

func NewWorkflowInstancesHandlers(
	service *service.WorkflowInstanceService,
) *WorkflowInstancesHandlers {
	return &WorkflowInstancesHandlers{
		service: service,
	}
}

func (h *WorkflowInstancesHandlers) Register(rg *gin.RouterGroup) {
	rg.GET("/workflow-instances", h.getAll)
	rg.GET("/workflow-instances/:id", h.getByID)
}

// getAll
// @Id 			getAllWorkflowInstances
// @Summary 	Get all workflow instances
// @Description Get all workflow instances
// @Tags 		Workflow Instances
// @Produce 	json
// @Success 	200 {array} models.WorkflowInstance
// @Router 		/api/workflow-instances [get]
func (h *WorkflowInstancesHandlers) getAll(ctx *gin.Context) {
	paginate := utils.Paginate(ctx)
	defs, err := h.service.GetAll(&paginate)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, defs)
}

// getByID
// @Id 			getWorkflowInstanceByID
// @Summary 	Get a workflow instance by ID
// @Description Get a workflow instance by ID
// @Tags 		Workflow Instances
// @Produce 	json
// @Param 		id path string true "Workflow Instance ID"
// @Success 	200 {object} models.WorkflowInstance
// @Router 		/api/workflow-instances/{id} [get]
func (h *WorkflowInstancesHandlers) getByID(ctx *gin.Context) {
	id := ctx.Param("id")
	def, err := h.service.GetByID(id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, def)
}
