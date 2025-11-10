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

// getAll
// @Id			 getAllAgents
// @Summary      Get all agents
// @Description  Get a list of all registered agents
// @Tags         agents
// @Produce      json
// @Success      200  {array}   models.AgentResponse
// @Router       /api/agents [get]
func (h *AgentsHandlers) getAll(ctx *gin.Context) {
	agents := h.registry.ListAgents()
	ctx.JSON(200, agents.ToResponseList())
}

// getByName
// @Id  		 getAgentByName
// @Summary      Get agent by name
// @Description  Get details of a registered agent by its name
// @Tags         agents
// @Produce      json
// @Param        name   path      string  true  "Agent Name"
// @Success      200    {object}  models.AgentResponse
// @Failure      404    {object}  map[string]string
// @Router       /api/agents/{name} [get]
func (h *AgentsHandlers) getByName(ctx *gin.Context) {
	name := ctx.Param("name")
	agent, has := h.registry.GetAgent(name)
	if has == false {
		ctx.JSON(404, gin.H{"error": "agent not found"})
		return
	}
	ctx.JSON(200, agent.ToResponse())
}
