package httpapi

import (
	"github.com/gin-gonic/gin"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/service"
	"github.com/paulhalleux/workflow-engine-go/engine/internal/utils"
)

type StepInstancesHandlers struct {
	service *service.StepInstanceService
}

func NewStepInstancesHandlers(
	service *service.StepInstanceService,
) *StepInstancesHandlers {
	return &StepInstancesHandlers{
		service: service,
	}
}

func (h *StepInstancesHandlers) Register(rg *gin.RouterGroup) {
	rg.GET("/step-instances", h.getAll)
	rg.GET("/step-instances/:id", h.getByID)
}

func (h *StepInstancesHandlers) getAll(ctx *gin.Context) {
	paginate := utils.Paginate(ctx)
	defs, err := h.service.GetAll(&paginate)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, defs)
}

func (h *StepInstancesHandlers) getByID(ctx *gin.Context) {
	id := ctx.Param("id")
	def, err := h.service.GetByID(id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, def)
}
