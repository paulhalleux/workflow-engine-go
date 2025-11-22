package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paulhalleux/workflow-engine-go/engine-new/internal/ws"
)

type Handler interface {
	Register(router gin.IRoutes)
}

type HttpServer struct {
	address string
	port    string
	gin     *gin.Engine
	server  *http.Server
	api     *gin.RouterGroup
}

func NewHttpServer(address, port string) HttpServer {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	server := &http.Server{
		Addr:    joinHostPort(address, port),
		Handler: r,
	}

	api := r.Group("/api")

	return HttpServer{
		address: address,
		port:    port,
		gin:     r,
		server:  server,
		api:     api,
	}
}

func (h *HttpServer) Start(
	wsHandler ws.WebsocketServer,
) {
	log.Printf("[engine] HTTP server listening on %s", h.server.Addr)

	h.gin.GET("/ws", func(c *gin.Context) {
		wsHandler.HandleWebSocket(c.Writer, c.Request)
	})

	if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprintf("failed to start HTTP server: %v", err))
	}
}

func (h *HttpServer) Stop() error {
	return h.server.Shutdown(context.Background())
}

func (h *HttpServer) RegisterHandler(handler Handler) {
	handler.Register(h.gin)
}

func (h *HttpServer) RegisterApiHandler(handler Handler) {
	handler.Register(h.api)
}

func joinHostPort(host, port string) string {
	if host == "" {
		return fmt.Sprintf(":%s", port)
	}
	return fmt.Sprintf("%s:%s", host, port)
}
