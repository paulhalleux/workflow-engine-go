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

func (h *TaskDefinitionsHandlers) getAll(ctx *gin.Context) {
	definitions := h.registry.ListTaskDefinitions()
	ctx.JSON(200, definitions)
}

func (h *TaskDefinitionsHandlers) getById(ctx *gin.Context) {
	id := ctx.Param("id")
	definition, has := h.registry.GetTaskDefinition(id)
	if has == false {
		ctx.JSON(404, gin.H{"error": "definition not found"})
		return
	}
	ctx.JSON(200, definition)
}
