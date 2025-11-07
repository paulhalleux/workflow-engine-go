package httpapi

import (
	"github.com/gin-gonic/gin"
	"github.com/paulhalleux/workflow-engine-go/engine/internal"
)

type TaskDefinitionsHandlers struct {
	registry *internal.AgentRegistry
}

func NewTaskDefinitionsHandlers(
	registry *internal.AgentRegistry,
) *TaskDefinitionsHandlers {
	return &TaskDefinitionsHandlers{
		registry: registry,
	}
}

func (h *TaskDefinitionsHandlers) Register(rg *gin.RouterGroup) {
	rg.GET("/task-definitions", h.getAll)
	rg.GET("/task-definitions/:id", h.getById)
}

// getAll
// @Id 			getAllTaskDefinitions
// @Summary 	Get all task definitions
// @Description Get all task definitions
// @Tags 		Task Definitions
// @Produce 	json
// @Success 	200 {array} models.TaskDefinitionResponse
// @Router 		/task-definitions [get]
func (h *TaskDefinitionsHandlers) getAll(ctx *gin.Context) {
	definitions := h.registry.ListTaskDefinitions()
	ctx.JSON(200, definitions)
}

// getById
// @Id 			getTaskDefinitionById
// @Summary 	Get task definition by ID
// @Description Get task definition by ID
// @Tags 		Task Definitions
// @Produce 	json
// @Param 		id path string true "Task Definition ID"
// @Success 	200 {object} models.TaskDefinitionResponse
// @Failure 	404 {object} gin.H
// @Router 		/task-definitions/{id} [get]
func (h *TaskDefinitionsHandlers) getById(ctx *gin.Context) {
	id := ctx.Param("id")
	definition, has := h.registry.GetTaskDefinition(id)
	if has == false {
		ctx.JSON(404, gin.H{"error": "definition not found"})
		return
	}
	ctx.JSON(200, definition)
}
