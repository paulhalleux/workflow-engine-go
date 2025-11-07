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

// getAll
// @Id 			getAllStepInstances
// @Summary 	Get all step instances
// @Description Get all step instances with pagination
// @Tags		StepInstances
// @Produce 	json
// @Param 		limit query int false "Page limit"
// @Param 		offset query int false "Page offset"
// @Success 	200 {array} models.StepInstance
// @Failure 	500 {object} gin.H
// @Router 		/step-instances [get]
func (h *StepInstancesHandlers) getAll(ctx *gin.Context) {
	paginate := utils.Paginate(ctx)
	defs, err := h.service.GetAll(&paginate)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, defs)
}

// getByID
// @Id 			getStepInstanceById
// @Summary 	Get step instance by ID
// @Description Get a step instance by its ID
// @Tags 		StepInstances
// @Produce 	json
// @Param 		id path string true "Step Instance ID"
// @Success 	200 {object} models.StepInstance
// @Failure 	500 {object} gin.H
// @Router 		/step-instances/{id} [get]
func (h *StepInstancesHandlers) getByID(ctx *gin.Context) {
	id := ctx.Param("id")
	def, err := h.service.GetByID(id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, def)
}
