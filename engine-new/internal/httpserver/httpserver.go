package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	address string
	port    string
	gin     *gin.Engine
	server  *http.Server
}

func NewHttpServer(address, port string) HttpServer {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	server := &http.Server{
		Addr:    joinHostPort(address, port),
		Handler: r,
	}

	return HttpServer{
		address: address,
		port:    port,
		gin:     r,
		server:  server,
	}
}

func (h *HttpServer) Start() {
	log.Printf("[engine] HTTP server listening on %s", h.server.Addr)
	if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprintf("failed to start HTTP server: %v", err))
	}
}

func (h *HttpServer) Stop() error {
	return h.server.Shutdown(context.Background())
}

func joinHostPort(host, port string) string {
	if host == "" {
		return fmt.Sprintf(":%s", port)
	}
	return fmt.Sprintf("%s:%s", host, port)
}
