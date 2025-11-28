package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/paulhalleux/workflow-engine-go/engine-new/internal/dto"
	"github.com/paulhalleux/workflow-engine-go/engine-new/internal/registry"
)

type AgentsHandlers struct {
	registry *registry.AgentRegistry
}

func NewAgentsHandlers(
	registry *registry.AgentRegistry,
) *AgentsHandlers {
	return &AgentsHandlers{
		registry: registry,
	}

}

func (w *AgentsHandlers) Register(router gin.IRoutes) {
	router.GET("/agents", w.GetAllAgents)
	router.GET("/agents/:name", w.GetAgentByName)
}

// GetAllAgents godoc
// @ID GetAllAgents
// @Summary Get all registered agents
// @Description Retrieve a list of all registered agents in the system.
// @Tags agents
// @Produce json
// @Success 200 {array} dto.AgentOverviewResponse
// @Router /api/agents [get]
func (w *AgentsHandlers) GetAllAgents(c *gin.Context) {
	c.JSON(200, dto.NewAgentsOverviewResponse(w.registry.ListAgents()))
}

// GetAgentByName godoc
// @ID GetAgentByName
// @Summary Get agent by name
// @Description Retrieve detailed information about a specific agent by its name.
// @Tags agents
// @Produce json
// @Param name path string true "Agent Name"
// @Success 200 {object} dto.AgentResponse
// @Failure 404 {object} gin.H
// @Router /api/agents/{name} [get]
func (w *AgentsHandlers) GetAgentByName(c *gin.Context) {
	agentName := c.Param("name")
	agent, found := w.registry.GetAgent(agentName)
	if !found {
		c.JSON(404, gin.H{"error": "agent not found"})
		return
	}
	c.JSON(200, dto.NewAgentResponse(agent))
}
