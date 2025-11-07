package httpapi

import (
	"github.com/gin-gonic/gin"
	"github.com/paulhalleux/workflow-engine-go/engine/internal"
)

type AgentsHandlers struct {
	registry *internal.AgentRegistry
}

func NewAgentsHandlers(
	registry *internal.AgentRegistry,
) *AgentsHandlers {
	return &AgentsHandlers{
		registry: registry,
	}
}

func (h *AgentsHandlers) Register(rg *gin.RouterGroup) {
	rg.GET("/agents", h.getAll)
	rg.GET("/agents/:name", h.getByName)
}

func (h *AgentsHandlers) getAll(ctx *gin.Context) {
	agents := h.registry.ListAgents()
	ctx.JSON(200, agents.ToResponseList())
}

func (h *AgentsHandlers) getByName(ctx *gin.Context) {
	name := ctx.Param("name")
	agent, has := h.registry.GetAgent(name)
	if has == false {
		ctx.JSON(404, gin.H{"error": "agent not found"})
		return
	}
	ctx.JSON(200, agent.ToResponse())
}
